package chart

import (
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Series is a named row of values, one per label.
type Series struct {
	Name   string
	Values []float64
}

// LineOptions tune a line chart.
type LineOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
}

// Line draws series as lines over axes.
func Line(c *ui.Context, label string, labels []string, series []Series, at *int, o LineOptions) ui.Element {
	return lines(c, label, labels, series, at, look{curve: mc.CurveMonotone, h: o.Height, format: o.Format, comp: "Line"})
}
