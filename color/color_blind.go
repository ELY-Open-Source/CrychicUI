package color

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ColorBlindOptions tune a ColorBlind.
type ColorBlindOptions struct {
	// Names title the normal, protan, deutan and tritan rows.
	Names [4]string
}

// ColorBlind shows colors as color-blind eyes see them.
func ColorBlind(c *ui.Context, colors []ui.Color, o ColorBlindOptions) ui.Element {
	if len(colors) == 0 {
		panic("color: ColorBlind needs colors")
	}
	if o.Names == ([4]string{}) {
		o.Names = [4]string{"Normal", "Protan", "Deutan", "Tritan"}
	}
	t := c.Theme()
	return ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).Children(func() {
		for k, name := range o.Names {
			c.Key(k)
			row := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(name)
			kit.Arrive(row, ui.Pose{X: -t.Space(2), Fade: 1}, time.Duration(k)*60*time.Millisecond).Children(func() {
				kit.Muted(c, name).Width(t.Space(18))
				for i, col := range colors {
					if k > 0 {
						col = seen(col, k-1)
					}
					c.Key(i)
					swatch(c, col, t.Space(8), t.Space(8))
				}
			})
		}
	})
}
