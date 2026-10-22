package overlay

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Emoji is a character and its searchable name.
type Emoji struct{ Char, Name string }

// EmojiGroup is a titled group of emoji.
type EmojiGroup struct {
	Title string
	Emoji []Emoji
}

// EmojiPickerOptions tune an EmojiPicker.
type EmojiPickerOptions struct {
	// Empty shows when search finds nothing.
	Empty string
}

// EmojiPicker is a searchable grid of emoji by group.
func EmojiPicker(c *ui.Context, label string, groups []EmojiGroup, chosen *string, o EmojiPickerOptions) ui.Element {
	if label == "" || len(groups) == 0 || chosen == nil {
		panic(fmt.Sprintf("overlay: EmojiPicker %q needs groups and a choice", label))
	}
	names, titles := map[string]bool{}, map[string]bool{}
	for _, g := range groups {
		if g.Title == "" || titles[g.Title] || len(g.Emoji) == 0 {
			panic(fmt.Sprintf("overlay: EmojiPicker %q groups need unique titles and emoji", label))
		}
		titles[g.Title] = true
		for _, e := range g.Emoji {
			if e.Char == "" || e.Name == "" || names[e.Name] {
				panic(fmt.Sprintf("overlay: EmojiPicker %q emoji need a character and a unique name", label))
			}
			names[e.Name] = true
		}
	}
	def(&o.Empty, "No emoji found")
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(80)).Gap(t.Space(2))
	query := ui.Local(root, "query", func() string { return "" })
	shown := *chosen
	root.Children(func() {
		ui.SearchField(c, query).Label(label).FillWidth()
		q := strings.ToLower(strings.TrimSpace(*query))
		ui.Scroll(c).FillWidth().MaxHeight(t.Space(64)).Children(func() {
			ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
				found := false
				for g, grp := range groups {
					var hits []Emoji
					for _, e := range grp.Emoji {
						if strings.Contains(strings.ToLower(e.Name), q) {
							hits = append(hits, e)
						}
					}
					if len(hits) == 0 {
						continue
					}
					found = true
					c.Key(grp.Title)
					ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
						kit.Muted(c, grp.Title).Bold().FontSize(t.Rem(0.75))
						ui.Row(c).Wrap().Gap(t.Space(0.5)).FocusGroup(ui.Horizontal).Children(func() {
							for _, e := range hits {
								c.Key(e.Name)
								b := kit.Press(ui.ButtonBase(c).Label(e.Name).Tooltip(e.Name).Size(t.Space(9), t.Space(9)).Center().Radius(t.Radius))
								bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
								if e.Char == shown {
									bg = t.Accent.Alpha(0.16)
								}
								b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
									glyph := ui.Text(c, e.Char).FontSize(t.Rem(1.3))
									s := 1 + 0.2*glyph.AnimateTo("grow", unit(b.Hovered()), ui.Bouncy)
									glyph.Scale(s, s)
								})
								if b.Clicked() {
									*chosen = e.Char
									root.SetChanged()
									sound.Play(sound.Select)
									kit.Log("EmojiPicker", "chose", "group", g)
								}
							}
						})
					})
				}
				if !found {
					kit.Muted(c, o.Empty).Padding(t.Space(3)).Key("empty").Transition(kit.Enter)
				}
			})
		})
	})
	return root
}
