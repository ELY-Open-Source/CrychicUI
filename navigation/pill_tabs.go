package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PillTabs are tabs whose pill slides to the chosen one.
func PillTabs(c *ui.Context, label string, names []string, active *int) ui.Element {
	tabs("PillTabs", label, names, active)
	t := c.Theme()
	shown := *active
	parts := ui.TabsBase(c, active, len(names))
	root := parts.List.Label(label).AlignSelf(ui.Start).Gap(t.Space(0.5)).Padding(t.Space(0.5)).Radius(999).Background(t.Surface)
	id := root.ID()
	root.Children(func() {
		for i, name := range names {
			c.Key(i)
			tab := parts.Tab(i).Radius(999)
			fg := t.TextMuted
			if i == shown || tab.Hovered() {
				fg = t.Text
			}
			tab.TextColor(tab.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
				stack(c, i == shown, func() ui.Element {
					return ui.Box(c).Radius(999).Background(t.Background).Shadow(0, 1, 2, 0, kit.Shade(0.12)).Transition(glide(id))
				}, func() ui.Element {
					return ui.Box(c).Padding(t.Space(1.5), t.Space(4)).Children(func() { ui.Text(c, name).SingleLine() })
				})
			})
		}
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("PillTabs", "picked", "index", *active)
	}
	return root
}
