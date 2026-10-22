package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var bob = []ui.Keyframe{{Value: 0}, {At: 1200 * time.Millisecond, Value: 1, Ease: ui.EaseInOut},
	{At: 2400 * time.Millisecond, Value: 0, Ease: ui.EaseInOut}}

// Bob floats content up and down in a loop.
func Bob(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: Bob needs content")
	}
	root := ui.Box(c)
	return root.Translate(0, -c.Theme().Space(1.5)*root.Keyframes("crychic.bob", bob, forever)).Children(content)
}
