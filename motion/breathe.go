package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var breathe = []ui.Keyframe{{Value: 1}, {At: 1600 * time.Millisecond, Value: 0.45, Ease: ui.EaseInOut},
	{At: 3200 * time.Millisecond, Value: 1, Ease: ui.EaseInOut}}

// Breathe pulses the opacity of content in a loop.
func Breathe(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: Breathe needs content")
	}
	root := ui.Box(c)
	return root.Opacity(root.Keyframes("crychic.breathe", breathe, forever)).Children(content)
}
