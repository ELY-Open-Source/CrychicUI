package typography

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Highlight sweeps a marker behind its text.
func Highlight(c *ui.Context, text string) ui.Element {
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	sweep := root.Keyframes("sweep", []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		if sweep > 0 {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H*0.35, W: r.W * sweep, H: r.H * 0.65}, t.Warning.Alpha(0.35), 0)
		}
	})
	return root.Children(func() { ui.Text(c, text) })
}
