package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var beat = []ui.Keyframe{{Value: 1}, {At: 140 * time.Millisecond, Value: 1.15, Ease: ui.EaseOut},
	{At: 280 * time.Millisecond, Value: 1, Ease: ui.EaseIn}, {At: 420 * time.Millisecond, Value: 1.1, Ease: ui.EaseOut},
	{At: 560 * time.Millisecond, Value: 1, Ease: ui.EaseIn}, {At: 1200 * time.Millisecond, Value: 1}}

// Beat scales content like a heartbeat in a loop.
func Beat(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: Beat needs content")
	}
	root := ui.Box(c)
	s := root.Keyframes("crychic.beat", beat, forever)
	return root.Scale(s, s).Children(content)
}
