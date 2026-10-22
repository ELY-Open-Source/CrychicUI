package chart

import (
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// StreamOptions tune a Stream.
type StreamOptions = LineOptions

// Stream stacks areas around a center line.
func Stream(c *ui.Context, label string, labels []string, series []Series, at *int, o StreamOptions) ui.Element {
	return lines(c, label, labels, series, at, look{curve: mc.CurveMonotone, fill: true, stack: true, center: true,
		h: o.Height, format: o.Format, comp: "Stream"})
}
