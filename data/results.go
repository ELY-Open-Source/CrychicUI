package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Result is a search result's title and snippet.
type Result struct{ Title, Snippet string }

// ResultsOptions tune Results.
type ResultsOptions struct {
	// Empty shows when there are no results.
	Empty string
}

// Results lists results with query words marked.
func Results(c *ui.Context, label, query string, results []Result, chosen *int, o ResultsOptions) ui.Element {
	if label == "" || chosen == nil || *chosen < -1 || *chosen >= len(results) {
		panic(fmt.Sprintf("data: Results %q needs a choice or -1", label))
	}
	for _, r := range results {
		if r.Title == "" {
			panic(fmt.Sprintf("data: Results %q results need titles", label))
		}
	}
	def(&o.Empty, "No results")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList).Label(label)
	root.Children(func() {
		if len(results) == 0 {
			kit.Muted(c, o.Empty).Padding(t.Space(3)).Key("empty").Transition(kit.Enter)
			return
		}
		for i, r := range results {
			c.Key(struct {
				q string
				i int
			}{query, i})
			b := ui.ButtonBase(c).Label(r.Title).FillWidth().Gap(t.Space(1)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius)
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if i == *chosen {
				bg = t.Accent.Alpha(0.1)
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy))
			kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
					ui.RichText(c, marked(t, r.Title, query)...).Bold().TextColor(t.Accent)
					if r.Snippet != "" {
						ui.RichText(c, marked(t, r.Snippet, query)...).TextColor(t.TextMuted).FontSize(t.Rem(0.9))
					}
				})
			})
			if b.Clicked() {
				*chosen = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Results", "chose", "index", i)
			}
		}
	})
	return root
}
