package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ConicBorder rings content with a turning conic gradient.
func ConicBorder(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("fx: ConicBorder needs content")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Padding(t.Space(4)).Radius(t.Radius).Border(t.Space(0.5), ui.Color{})
	cols := palette(t)
	stops := make([]ui.GradientStop, 0, len(cols)+1)
	for i, col := range append(cols, cols[0]) {
		stops = append(stops, at(col, float32(i)/float32(len(cols))))
	}
	turn := 360 * phase(c, root, 6*time.Second)
	return root.BorderConicGradient(ui.ConicGradient{Angle: turn, Stops: stops, Oklab: true}).Children(content)
}
