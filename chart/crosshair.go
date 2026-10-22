package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// CrosshairOptions tune a Crosshair.
type CrosshairOptions = LineOptions

// Crosshair crosses the hovered point, reading values.
func Crosshair(c *ui.Context, label string, labels []string, series []Series, at *int, o CrosshairOptions) ui.Element {
	def(&o.Format, "%g")
	t := c.Theme()
	return lines(c, label, labels, series, at, look{curve: mc.CurveMonotone, h: o.Height, format: o.Format, comp: "Crosshair",
		over: func(_ ui.Element, _ mc.Linear, _ []float32, aims [][]float32, i int) {
			if i < 0 {
				return
			}
			for s := range series {
				c.Key(fmt.Sprintf("rule.%d", s))
				mc.RuleY(c, aims[s][i], 1).Background(tone(t, s).Alpha(0.5)).Transition(spring)
				c.Key(fmt.Sprintf("read.%d", s))
				mc.Anchor(c, 0, aims[s][i], ui.Start, ui.Center, func() {
					ui.Text(c, fmt.Sprintf(o.Format, series[s].Values[i])).NoWrap().FontSize(t.Rem(0.7)).Bold().
						Padding(t.Space(0.25), t.Space(1)).Radius(t.Radius / 2).Background(tone(t, s)).TextColor(t.AccentText)
				}).Transition(spring)
			}
		}})
}
