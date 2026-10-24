package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Notice is one notification row.
type Notice struct {
	ID, Text, Time, Icon string
}

// Notification opens a notification row; cue Tap.
func Notification(c *ui.Context, notice Notice, opened *bool) ui.Element {
	if notice.ID == "" || notice.Text == "" || notice.Time == "" || opened == nil {
		panic("mail: Notification needs notice and state")
	}
	if notice.Icon == "" {
		notice.Icon = "bell"
	}
	t := c.Theme()
	root := button(c, notice.Text).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
		Border(1, t.Border).Padding(t.Space(2))
	root.Children(func() {
		ui.Icon(c, icon.Must(notice.Icon)).Size(16, 16).TextColor(t.Accent)
		ui.Text(c, notice.Text).Grow(1).Wrap()
		kit.Muted(c, notice.Time).SingleLine()
	})
	if root.Clicked() {
		*opened = true
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("Notification", "opened")
	}
	return root
}
