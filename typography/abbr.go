package typography

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Abbr is dotted text with a full-term tooltip.
func Abbr(c *ui.Context, abbr, full string) ui.Element {
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Tooltip(full)
	p := root.Keyframes("dots", []ui.Keyframe{{Value: 0}, {At: 160 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	root.DrawOver(func(painter *ui.Painter, r ui.Rect) {
		for x := r.X; x < r.X+r.W; x += 3 {
			painter.Fill(ui.Rect{X: x, Y: r.Y + r.H - 1, W: 1, H: 1}, t.TextMuted.Alpha(p), 0.5)
		}
	})
	return root.Children(func() { ui.Text(c, abbr) })
}
