package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var wiggle = []ui.Keyframe{{Value: 0}, {At: 80 * time.Millisecond, Value: -6}, {At: 160 * time.Millisecond, Value: 6},
	{At: 240 * time.Millisecond, Value: -4}, {At: 320 * time.Millisecond, Value: 4}, {At: 400 * time.Millisecond, Value: 0}}

// Wiggle rocks content a few degrees each time run changes.
func Wiggle(c *ui.Context, run uint64, content func()) ui.Element {
	if content == nil {
		panic("motion: Wiggle needs content")
	}
	root := ui.Box(c)
	return root.Rotate(onChange(root, run, wiggle)).Children(content)
}
