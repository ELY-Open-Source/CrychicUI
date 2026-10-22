package chart

import (
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// StepLineOptions tune a StepLine.
type StepLineOptions = LineOptions

// StepLine draws series in steps.
func StepLine(c *ui.Context, label string, labels []string, series []Series, at *int, o StepLineOptions) ui.Element {
	return lines(c, label, labels, series, at, look{curve: mc.CurveStepAfter, h: o.Height, format: o.Format, comp: "StepLine"})
}
