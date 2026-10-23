package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MetricTile shows value, delta, and sparkline.
func MetricTile(c *ui.Context, label, value, delta string, points []float64) ui.Element {
	if label == "" || value == "" || delta == "" || len(points) < 2 {
		panic("dashboard: MetricTile has invalid data")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		kit.Muted(c, label).FontSize(t.Rem(0.75))
		ui.Row(c).FillWidth().Children(func() {
			kit.Numeric(c, value).Bold().Grow(1)
			ui.Text(c, delta).TextColor(t.Accent)
		})
		values := make([]Point, len(points))
		for i, point := range points {
			values[i] = Point{Value: point}
		}
		series(c, values, t.Accent)
	})
	return root
}
