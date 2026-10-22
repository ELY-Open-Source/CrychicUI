package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ConicSpin turns a soft conic gradient behind content.
func ConicSpin(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	cols := palette(t)
	stops := make([]ui.GradientStop, 0, len(cols)+1)
	for i, col := range append(cols, cols[0]) {
		stops = append(stops, at(col.Mix(t.Surface, 0.6), float32(i)/float32(len(cols))))
	}
	turn := 360 * phase(c, root, 10*time.Second)
	return root.ConicGradient(ui.ConicGradient{Angle: turn, Stops: stops, Oklab: true}).Children(content)
}
