package chat

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// UnreadDivider is a "new messages" line that fades.
func UnreadDivider(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("chat: UnreadDivider needs a label")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
	since := ui.Local(root, "since", func() time.Time { return c.Now() })
	fade := float32(1)
	if left := 3*time.Second - c.Now().Sub(*since); left <= 0 {
		fade = 0
	} else {
		c.After(left)
	}
	root.Opacity(root.AnimateTo("fade", fade, ui.Smooth))
	return root.Children(func() {
		ui.Box(c).Grow(1).Height(1).Background(t.Accent.Alpha(0.4))
		ui.Text(c, label).FontSize(t.Rem(0.8)).Bold().TextColor(t.Accent).SingleLine()
		ui.Box(c).Grow(1).Height(1).Background(t.Accent.Alpha(0.4))
	})
}
