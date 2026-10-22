package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// refade clears first; replays start from the shown value.
var refade = []ui.Keyframe{{Value: 0}, {At: time.Millisecond, Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}}

// FadeIn fades content in when built and when run changes.
func FadeIn(c *ui.Context, run uint64, content func()) ui.Element {
	if content == nil {
		panic("motion: FadeIn needs content")
	}
	root := ui.Box(c)
	return root.Opacity(root.Keyframes("crychic.in", refade, ui.KeyframeOptions{Run: run})).Children(content)
}
