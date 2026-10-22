package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VerticalTabs are tabs in a column; the indicator slides.
func VerticalTabs(c *ui.Context, label string, names []string, active *int) ui.Element {
	tabs("VerticalTabs", label, names, active)
	t := c.Theme()
	shown := *active
	parts := ui.TabsBase(c, active, len(names))
	root := parts.List.Label(label).Column().AlignSelf(ui.Start).Width(t.Space(48)).Gap(t.Space(0.5))
	id := root.ID()
	root.Children(func() {
		for i, name := range names {
			c.Key(i)
			tab := parts.Tab(i).FillWidth().Radius(t.Radius)
			bg, fg := t.SurfaceHover.Alpha(0), t.TextMuted
			switch {
			case i == shown:
				fg = t.Text
			case tab.Hovered():
				bg, fg = t.SurfaceHover, t.Text
			}
			tab.Background(tab.AnimateColor("bg", bg, ui.Snappy)).TextColor(tab.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
				ui.Row(c).FillWidth().AlignItems(ui.Stretch).Children(func() {
					bar := ui.Box(c).Width(2).Radius(999)
					if i == shown {
						bar.Background(t.Accent).Transition(glide(id))
					}
					ui.Text(c, name).SingleLine().Padding(t.Space(2), t.Space(3))
				})
			})
		}
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("VerticalTabs", "picked", "index", *active)
	}
	return root
}
