package chart

import (
	"github.com/egoist/mygo/ui"
)

// Point is a named place; Size scales a bubble.
type Point struct {
	Name       string
	X, Y, Size float64
}

// ScatterOptions tune a Scatter.
type ScatterOptions = LineOptions

// Scatter places points on two axes.
func Scatter(c *ui.Context, label string, points []Point, at *int, o ScatterOptions) ui.Element {
	return dots(c, label, points, at, o, false, "Scatter")
}
