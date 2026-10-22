package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// ThresholdOptions tune a Threshold.
type ThresholdOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
	// Limit words the limit line; it takes the value.
	Limit string
}

// Threshold draws a line; points past limit turn danger.
func Threshold(c *ui.Context, label string, labels []string, vs []float64, limit float64, at *int, o ThresholdOptions) ui.Element {
	finite("Threshold", label, limit)
	def(&o.Format, "%g")
	def(&o.Limit, "Limit %s")
	t := c.Theme()
	return lines(c, label, labels, []Series{{Name: label, Values: vs}}, at, look{curve: mc.CurveMonotone, h: o.Height,
		format: o.Format, comp: "Threshold", also: []float64{limit},
		over: func(_ ui.Element, y mc.Linear, xs []float32, aims [][]float32, _ int) {
			c.Key("limit")
			mc.RuleY(c, y.At(limit), 1).Background(t.Danger.Alpha(0.6)).Transition(spring)
			c.Key("limit.label")
			mc.Anchor(c, 1, y.At(limit), ui.End, ui.End, func() {
				kit.Muted(c, fmt.Sprintf(o.Limit, fmt.Sprintf(o.Format, limit))).NoWrap().FontSize(t.Rem(0.7)).
					TextColor(t.Danger).Margin(0, 0, t.Space(0.5), 0)
			}).Transition(spring)
			for i, v := range vs {
				c.Key(fmt.Sprintf("pt.%d", i))
				ink := tone(t, 0)
				if v > limit {
					ink = t.Danger
				}
				d := mc.Point(c, xs[i], aims[0][i], t.Space(2.5), t.Space(2.5)).Radius(999).Transition(spring)
				d.Background(d.AnimateColor("ink", ink, ui.Snappy))
			}
		}})
}
