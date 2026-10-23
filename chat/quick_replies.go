package chat

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuickReplies are reply chips; the tapped one sends.
func QuickReplies(c *ui.Context, label string, replies []string, chosen *int) ui.Element {
	if label == "" || len(replies) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(replies) {
		panic(fmt.Sprintf("chat: QuickReplies %q needs a label, replies and a choice or -1", label))
	}
	for _, r := range replies {
		if r == "" {
			panic(fmt.Sprintf("chat: QuickReplies %q replies must not be empty", label))
		}
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1.5)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	return root.Children(func() {
		for i, r := range replies {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(r).Padding(t.Space(1.5), t.Space(3)).Radius(999).Border(1, t.Accent).TextColor(t.Accent))
			b.Background(b.AnimateColor("bg", t.Accent.Alpha(0.12*unit(b.Hovered())), ui.Snappy))
			kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond)
			b.Children(func() { ui.Text(c, r).SingleLine() })
			if b.Clicked() {
				*chosen = i
				root.SetChanged()
				sound.Play(sound.Send)
				kit.Log("QuickReplies", "sent", "index", i)
			}
		}
	})
}
