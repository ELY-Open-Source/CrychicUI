package chat

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// EmojiBurst is a big emoji message that bounces.
func EmojiBurst(c *ui.Context, emoji string) ui.Element {
	if emoji == "" {
		panic("chat: EmojiBurst needs an emoji")
	}
	t := c.Theme()
	root := ui.Text(c, emoji).FontSize(t.Rem(3)).AlignSelf(ui.Start).Label(emoji)
	k := root.Keyframes("burst", []ui.Keyframe{
		{Value: 0.2},
		{At: 140 * time.Millisecond, Value: 1.15, Ease: ui.EaseOut},
		{At: 260 * time.Millisecond, Value: 1, Ease: ui.EaseOut},
	})
	return root.Scale(k, k)
}
