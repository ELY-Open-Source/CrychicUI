package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SidebarLayout puts a side column by main; side springs open.
func SidebarLayout(c *ui.Context, open *bool, side, main func()) ui.Element {
	if open == nil {
		panic("layout: SidebarLayout needs open")
	}
	t := c.Theme()
	root := ui.Row(c).Fill().AlignItems(ui.Stretch)
	return root.Children(func() {
		col := ui.Column(c).ClipX().Background(t.Surface)
		w := float32(0)
		if *open {
			w = t.Rem(16)
		}
		col.Width(col.AnimateTo("w", w, ui.Snappy)).Children(func() {
			ui.Column(c).Width(t.Rem(16)).Grow(1).Padding(t.Space(3)).Gap(t.Space(1)).Children(side)
		})
		ui.Column(c).Grow(1).Basis(0).Children(func() {
			glyph := "panel-left-open"
			if *open {
				glyph = "panel-left-close"
			}
			b := ui.Button(c, "").Label("Toggle sidebar").Margin(t.Space(2)).Children(func() {
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
			})
			if kit.Clicked(kit.Press(b), sound.Toggle) {
				*open = !*open
				root.SetChanged()
				kit.Log("SidebarLayout", "toggled", "open", *open)
			}
			ui.Column(c).Grow(1).Children(main)
		})
	})
}
