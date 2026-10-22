package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Shake shakes content sideways each time run changes.
func Shake(c *ui.Context, run uint64, content func()) ui.Element {
	if content == nil {
		panic("motion: Shake needs content")
	}
	d := c.Theme().Space(1.5)
	ms := time.Millisecond
	root := ui.Box(c)
	x := onChange(root, run, []ui.Keyframe{{Value: 0}, {At: 60 * ms, Value: -d}, {At: 120 * ms, Value: d},
		{At: 180 * ms, Value: -0.6 * d}, {At: 240 * ms, Value: 0.6 * d}, {At: 300 * ms, Value: 0}})
	return root.Translate(x, 0).Children(content)
}
