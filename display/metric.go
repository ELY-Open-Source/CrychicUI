package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Metric is a rolling value with a caller's unit.
func Metric(c *ui.Context, value float64, unit string) ui.Element {
	if unit == "" {
		panic("display: Metric needs a unit")
	}
	return kit.Numeric(c, formatNumber(c, value)+" "+unit)
}
