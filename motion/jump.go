package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Jump bounces content once each time run changes.
func Jump(c *ui.Context, run uint64, content func()) ui.Element {
	if content == nil {
		panic("motion: Jump needs content")
	}
	t := c.Theme()
	ms := time.Millisecond
	root := ui.Box(c)
	y := onChange(root, run, []ui.Keyframe{{Value: 0}, {At: 160 * ms, Value: -t.Space(3), Ease: ui.EaseOut},
		{At: 320 * ms, Value: 0, Ease: ui.EaseIn}, {At: 400 * ms, Value: -t.Space(1), Ease: ui.EaseOut},
		{At: 480 * ms, Value: 0, Ease: ui.EaseIn}})
	return root.Translate(0, y).Children(content)
}
