package chart

import "github.com/egoist/mygo/ui"

// SplineOptions tune a Spline.
type SplineOptions = LineOptions

// Spline draws series as smooth curves through their points.
func Spline(c *ui.Context, label string, labels []string, series []Series, at *int, o SplineOptions) ui.Element {
	return lines(c, label, labels, series, at, look{smooth: true, h: o.Height, format: o.Format, comp: "Spline"})
}
