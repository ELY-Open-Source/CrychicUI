package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// KeyPath loops a dot around a rectangle of keyframes.
func KeyPath(c *ui.Context) ui.Element {
	t := c.Theme()
	w, h, d := t.Space(40), t.Space(16), t.Space(3)
	ms := time.Millisecond
	xs := []ui.Keyframe{{Value: 0}, {At: 600 * ms, Value: w - d, Ease: ui.EaseInOut}, {At: 1200 * ms, Value: w - d},
		{At: 1800 * ms, Value: 0, Ease: ui.EaseInOut}, {At: 2400 * ms, Value: 0}}
	ys := []ui.Keyframe{{Value: 0}, {At: 600 * ms, Value: 0}, {At: 1200 * ms, Value: h - d, Ease: ui.EaseInOut},
		{At: 1800 * ms, Value: h - d}, {At: 2400 * ms, Value: 0, Ease: ui.EaseInOut}}
	root := ui.Box(c).Size(w, h).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
		var track ui.Path
		x0, y0, x1, y1 := r.X+d/2, r.Y+d/2, r.X+r.W-d/2, r.Y+r.H-d/2
		p.StrokePath(track.MoveTo(x0, y0).LineTo(x1, y0).LineTo(x1, y1).LineTo(x0, y1).Close(), 1, t.Border)
	})
	return root.Children(func() {
		dot := ui.Box(c).Absolute().Top(0).Left(0).Size(d, d).Radius(999).Background(t.Accent)
		dot.Translate(dot.Keyframes("x", xs, forever), dot.Keyframes("y", ys, forever))
	})
}
