package chart

import (
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// StackedAreaOptions tune a StackedArea.
type StackedAreaOptions = LineOptions

// StackedArea piles filled series on each other.
func StackedArea(c *ui.Context, label string, labels []string, series []Series, at *int, o StackedAreaOptions) ui.Element {
	return lines(c, label, labels, series, at, look{curve: mc.CurveMonotone, fill: true, stack: true, h: o.Height, format: o.Format, comp: "StackedArea"})
}
