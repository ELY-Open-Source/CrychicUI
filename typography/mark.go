package typography

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Mark fades in a search-match highlight.
func Mark(c *ui.Context, text string) ui.Element {
	t := c.Theme()
	e := ui.Text(c, text)
	p := e.Keyframes("mark", []ui.Keyframe{{Value: 0}, {At: 160 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	return e.TextBackground(t.Warning.Alpha(0.35 * p))
}
