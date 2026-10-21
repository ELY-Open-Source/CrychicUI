package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Watermark tiles faint text over its content.
func Watermark(c *ui.Context, text string, content func()) ui.Element {
	t := c.Theme()
	root := ui.Box(c)
	return root.Children(func() {
		content()
		overlay := ui.Box(c).Absolute().Fill().Clip().PassThrough()
		r := overlay.LayoutBounds()
		kit.Settle(c, root, []float32{r.W, r.H}, func() []float32 {
			b := overlay.LayoutBounds()
			return []float32{b.W, b.H}
		})
		fade := overlay.Keyframes("fade", []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		overlay.Opacity(fade).Children(func() {
			for y := float32(0); y < r.H; y += t.Space(12) {
				for x := float32(0); x < r.W; x += t.Space(18) {
					ui.Text(c, text).Absolute().Left(x).Top(y).Rotate(-20).TextColor(t.Text.Alpha(0.06)).PassThrough()
				}
			}
		})
	})
}
