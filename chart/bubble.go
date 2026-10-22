package chart

import "github.com/egoist/mygo/ui"

// BubbleOptions tune a Bubble.
type BubbleOptions = LineOptions

// Bubble places points sized by Size.
func Bubble(c *ui.Context, label string, points []Point, at *int, o BubbleOptions) ui.Element {
	return dots(c, label, points, at, o, true, "Bubble")
}
