package chart

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Note labels point At of series Series.
type Note struct {
	At, Series int
	Text       string
}

// AnnotatedOptions tune an Annotated.
type AnnotatedOptions = LineOptions

// Annotated draws lines with labeled notes on points.
func Annotated(c *ui.Context, label string, labels []string, series []Series, notes []Note, at *int, o AnnotatedOptions) ui.Element {
	for _, n := range notes {
		if n.Text == "" || n.At < 0 || n.At >= len(labels) || n.Series < 0 || n.Series >= len(series) {
			panic(fmt.Sprintf("chart: Annotated %q notes need text on a known point", label))
		}
	}
	t := c.Theme()
	return lines(c, label, labels, series, at, look{curve: mc.CurveMonotone, h: o.Height, format: o.Format, comp: "Annotated",
		over: func(_ ui.Element, _ mc.Linear, xs []float32, aims [][]float32, _ int) {
			for k, n := range notes {
				x, y := xs[n.At], aims[n.Series][n.At]
				c.Key(fmt.Sprintf("dot.%d", k))
				mc.Point(c, x, y, t.Space(2.5), t.Space(2.5)).Radius(999).Background(t.Warning).Transition(spring)
				c.Key(fmt.Sprintf("note.%d", k))
				mc.Anchor(c, x, y, ui.Center, ui.End, func() {
					note := ui.Text(c, n.Text).NoWrap().FontSize(t.Rem(0.7)).Padding(t.Space(0.5), t.Space(1.5)).
						Radius(t.Radius).Background(t.Warning.Mix(t.Background, 0.8)).Border(1, t.Warning).
						Margin(0, 0, t.Space(2.5), 0)
					kit.Arrive(note, ui.Pose{Y: t.Space(2), Fade: 1}, 300*time.Millisecond+time.Duration(k)*80*time.Millisecond)
				}).Transition(spring)
			}
		}})
}
