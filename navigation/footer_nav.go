package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FooterNav is columns of links under titles.
func FooterNav(c *ui.Context, label string, gs []LinkGroup, picked *[2]int) ui.Element {
	groups("FooterNav", label, gs, picked)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(10)).Padding(t.Space(4)).AlignItems(ui.Start).Label(label)
	return root.Children(func() {
		for g, grp := range gs {
			c.Key(g)
			col := ui.Column(c).Gap(t.Space(2)).AlignItems(ui.Start)
			kit.Arrive(col, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() {
				ui.Text(c, grp.Title).Bold().FontSize(t.Rem(0.85)).SingleLine()
				for i, l := range grp.Links {
					c.Key(i)
					b := kit.Press(ui.ButtonBase(c).Label(l).Radius(t.Radius))
					b.TextColor(b.AnimateColor("fg", t.TextMuted.Mix(t.Text, unit(b.Hovered())), ui.Snappy)).Children(func() {
						ui.Text(c, l).SingleLine()
					})
					if kit.Clicked(b, sound.Tap) {
						*picked = [2]int{g, i}
						root.SetChanged()
						kit.Log("FooterNav", "picked", "group", g, "index", i)
					}
				}
			})
		}
	})
}
