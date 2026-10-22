package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Attention pulses a ring around its content while on.
func Attention(c *ui.Context, on bool, content func()) ui.Element {
	if content == nil {
		panic("feedback: Attention needs content")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	return root.Children(func() {
		content()
		if on {
			ring := ui.Box(c).Absolute().Top(0).Left(0).Right(0).Bottom(0).Radius(t.Radius + t.Space(1)).PassThrough()
			k := ring.Loop("pulse", 1400*time.Millisecond, ui.EaseOut)
			ring.Border(t.Space(0.5), t.Accent.Alpha(0.7*(1-k))).Scale(1+0.08*k, 1+0.08*k)
		}
	})
}
