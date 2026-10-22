package chart

import (
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// AreaOptions tune an Area.
type AreaOptions = LineOptions

// Area fills under each line series.
func Area(c *ui.Context, label string, labels []string, series []Series, at *int, o AreaOptions) ui.Element {
	return lines(c, label, labels, series, at, look{curve: mc.CurveMonotone, fill: true, h: o.Height, format: o.Format, comp: "Area"})
}
