package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MetricTileOptions word the tile; Note sits under the value.
type MetricTileOptions struct{ Icon, Note string }

// MetricTile shows a value, its delta and a sparkline.
func MetricTile(c *ui.Context, label, value, delta string, points []float64, o MetricTileOptions) ui.Element {
	if label == "" || value == "" || len(points) < 2 {
		panic(fmt.Sprintf("dashboard: MetricTile %q needs a label, a value and two points", label))
	}
	def(&o.Icon, "activity")
	t := c.Theme()
	ink := t.Accent
	if falling(delta) {
		ink = t.Danger
	}
	root := card(c).Gap(t.Space(2)).Padding(t.Space(4))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, o.Icon, ink)
			kit.Muted(c, label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
			if delta != "" {
				deltaChip(c, delta)
			}
		})
		ui.Column(c).Gap(t.Space(0.5)).Children(func() {
			kit.Numeric(c, value).Bold().FontSize(t.Rem(1.8)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			if o.Note != "" {
				kit.Muted(c, o.Note).FontSize(t.Rem(0.75)).Shrink(0)
			}
		})
		values := make([]Point, len(points))
		for i, point := range points {
			values[i] = Point{Value: point}
		}
		series(c, fmt.Sprintf("MetricTile %q", label), values, ink).Height(t.Space(14))
	})
	return root
}
