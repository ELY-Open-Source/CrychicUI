package layout

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Board scrolls fixed-width columns sideways; columns slide in.
func Board(c *ui.Context, columns []func()) ui.Element {
	t := c.Theme()
	return ui.ScrollHorizontal(c).FillWidth().Children(func() {
		ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Start).Padding(t.Space(1)).Children(func() {
			for i, col := range columns {
				b := ui.Column(c).Key(i).Width(t.Rem(18)).Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).
					Background(t.Surface)
				kit.Arrive(b, ui.Pose{X: 16, Fade: 1}, time.Duration(i)*50*time.Millisecond).Children(col)
			}
		})
	})
}
