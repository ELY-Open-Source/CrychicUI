package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuickLinks is a grid of lifting icon tiles.
func QuickLinks(c *ui.Context, label string, links []NavItem, picked *int) ui.Element {
	if label == "" || len(links) == 0 || picked == nil || *picked < -1 || *picked >= len(links) {
		panic(fmt.Sprintf("navigation: QuickLinks %q needs links and a pick or -1", label))
	}
	for _, l := range links {
		if l.Label == "" {
			panic(fmt.Sprintf("navigation: QuickLinks %q links need labels", label))
		}
		icon.Must(l.Icon)
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).Label(label)
	return root.Children(func() {
		for i, l := range links {
			c.Key(i)
			tile := kit.Press(ui.ButtonBase(c).Label(l.Label).FillWidth().Padding(t.Space(3)).Radius(t.Radius).Background(t.Surface).Border(1, t.Border))
			k := tile.AnimateTo("lift", unit(tile.Hovered()), ui.Snappy)
			tile.Translate(0, -k*t.Space(1)).Shadow(0, 1+3*k, 2+6*k, 0, kit.Shade(0.06+0.08*k)).Children(func() {
				ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
					ui.Icon(c, icon.Must(l.Icon)).Size(16, 16).TextColor(t.Accent)
					ui.Text(c, l.Label).SingleLine()
				})
			})
			if kit.Clicked(tile, sound.Tap) {
				*picked = i
				root.SetChanged()
				kit.Log("QuickLinks", "picked", "index", i)
			}
		}
	})
}
