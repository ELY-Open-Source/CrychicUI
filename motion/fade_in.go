package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var refade = replay(200*time.Millisecond, ui.EaseOut)

// FadeIn fades content in when built and when run changes.
func FadeIn(c *ui.Context, run uint64, content func()) ui.Element {
	if content == nil {
		panic("motion: FadeIn needs content")
	}
	root := ui.Box(c)
	return root.Opacity(root.Keyframes("crychic.in", refade, ui.KeyframeOptions{Run: run})).Children(content)
}
