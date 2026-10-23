package chat

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MessageTime shows a message's time on hover.
func MessageTime(c *ui.Context, at time.Time, content func()) ui.Element {
	if at.IsZero() || content == nil {
		panic("chat: MessageTime needs a time and content")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End)
	show := root.AnimateTo("show", unit(root.Hovered()), ui.Snappy)
	return root.Children(func() {
		content()
		kit.Muted(c, when(c, at)).FontSize(t.Rem(0.75)).Shrink(0).Opacity(show)
	})
}
