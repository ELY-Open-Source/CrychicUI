package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LegendItem pairs a swatch color with a label.
type LegendItem struct {
	Label string
	Color ui.Color
}

// Legend fades in swatches with labels.
func Legend(c *ui.Context, items []LegendItem) ui.Element {
	t := c.Theme()
	return ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).Children(func() {
		for i, item := range items {
			if item.Label == "" || item.Color.A == 0 {
				panic("display: Legend items need a label and color")
			}
			row := ui.Row(c).Key(i).Gap(t.Space(1)).AlignItems(ui.Center)
			kit.Arrive(row, ui.Pose{Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(t.Space(0.5)).Background(item.Color)
				ui.Text(c, item.Label).FontSize(t.Rem(0.85))
			})
		}
	})
}
