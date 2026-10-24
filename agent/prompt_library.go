package agent

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SavedPrompt is one prompt of a PromptLibrary.
type SavedPrompt struct {
	Title, Body, Tag string
}

// PromptLibraryOptions tune a PromptLibrary.
type PromptLibraryOptions struct {
	Placeholder, Empty string
}

// PromptLibrary searches saved prompts; a click selects one.
func PromptLibrary(c *ui.Context, label string, prompts []SavedPrompt, query *string, chosen *int, o PromptLibraryOptions) (root, field ui.Element) {
	if label == "" || len(prompts) == 0 || query == nil || chosen == nil || *chosen < -1 || *chosen >= len(prompts) {
		panic(fmt.Sprintf("agent: PromptLibrary %q needs a label, prompts, a query and a choice or -1", label))
	}
	seen := map[string]bool{}
	for _, p := range prompts {
		if p.Title == "" || p.Body == "" || seen[p.Title] {
			panic(fmt.Sprintf("agent: PromptLibrary %q prompts need unique titles and a body", label))
		}
		seen[p.Title] = true
	}
	def(&o.Placeholder, "Search prompts")
	def(&o.Empty, "No prompts match")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		field = input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder})
		q := strings.ToLower(strings.TrimSpace(*query))
		var hits []int
		for i, p := range prompts {
			hay := strings.ToLower(p.Title + "\n" + p.Body + "\n" + p.Tag)
			if q == "" || strings.Contains(hay, q) {
				hits = append(hits, i)
			}
		}
		ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Children(func() {
			if len(hits) == 0 {
				empty := ui.Column(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(6))
				kit.Arrive(empty, ui.Pose{Fade: 1}, 0).Children(func() {
					ui.Icon(c, icon.Must("search-x")).Size(16, 16).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty)
				})
				return
			}
			for n, i := range hits {
				p := prompts[i]
				c.Key(p.Title)
				on := i == *chosen
				b := kit.Press(ui.ButtonBase(c).Label(p.Title).Column().FillWidth().Gap(t.Space(0.5)).AlignItems(ui.Start).
					Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).Checked(on))
				kit.Arrive(b, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(n)*25*time.Millisecond)
				bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
				if on {
					bg = t.Accent.Alpha(0.1)
				}
				b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icon.Must("message-square-text")).Size(16, 16).TextColor(t.Accent)
						ui.Text(c, p.Title).Bold().Grow(1).SingleLine()
						if p.Tag != "" {
							tag(c, p.Tag)
						}
					})
					kit.Muted(c, p.Body).FontSize(t.Rem(0.8)).SingleLine()
				})
				if b.Clicked() && !on {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("PromptLibrary", "selected", "index", i)
				}
			}
		})
	})
	return root, field
}
