package choice

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var column = ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{X: -16}, Exit: &ui.Motion{X: -16}}

// Cascade is sliding columns that drill down a tree.
func Cascade(c *ui.Context, label string, items []TreeItem, path *[]string) ui.Element {
	if label == "" || len(items) == 0 || path == nil {
		panic("choice: " + who("Cascade", label) + " needs a label, items and a path")
	}
	at := items
	for _, step := range *path {
		i := slices.IndexFunc(at, func(it TreeItem) bool { return it.Label == step })
		if i < 0 {
			panic("choice: " + who("Cascade", label) + " path leaves the tree")
		}
		at = at[i].Children
	}
	t := c.Theme()
	root := ui.Row(c).Label(label).AlignSelf(ui.Start).AlignItems(ui.Stretch).Radius(t.Radius).Border(1, t.Border).
		Background(t.Background).Clip()
	return root.Children(func() {
		level := items
		for depth := 0; len(level) > 0; depth++ {
			labels := make([]string, len(level))
			for i, it := range level {
				labels[i] = it.Label
			}
			distinct(who("Cascade", label), labels)
			chosen := ""
			if depth < len(*path) {
				chosen = (*path)[depth]
			}
			// Columns are keyed by their ancestor path.
			c.Key(fmt.Sprintf("%q", (*path)[:depth]))
			col := ui.Column(c).Width(t.Space(40)).Padding(t.Space(1)).Gap(t.Space(0.5)).Transition(column)
			if depth > 0 {
				col.BorderStartWidth(1).BorderColor(t.Border)
			}
			var next []TreeItem
			col.Children(func() {
				for _, it := range level {
					c.Key(it.Label)
					row := ui.ButtonBase(c).Label(it.Label).Justify(ui.Start).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2.5)).
						Radius(t.Radius)
					bg, fg := ui.Color{}, t.Text
					if it.Label == chosen {
						bg, fg, next = t.Accent.Alpha(0.14), t.Accent, it.Children
					} else if row.Hovered() {
						bg = t.SurfaceHover
					}
					row.Background(row.AnimateColor("bg", bg, quick)).TextColor(fg).Children(func() {
						ui.Text(c, it.Label).SingleLine().Grow(1)
						if len(it.Children) > 0 {
							ui.Icon(c, icon.Must("chevron-right")).Size(16, 16)
						}
					})
					if kit.Clicked(row, sound.Select) && it.Label != chosen {
						*path = append(slices.Clone((*path)[:depth]), it.Label)
						root.SetChanged()
						kit.Log("Cascade", "chose", "depth", depth)
					}
				}
			})
			level = next
		}
	})
}
