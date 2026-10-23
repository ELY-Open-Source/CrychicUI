package chat

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SendState is where an outgoing message stands.
type SendState uint8

// Send states.
const (
	Sending SendState = iota
	Failed
	sendStateCount
)

// MessageStatusOptions name the states and retry.
type MessageStatusOptions struct {
	Sending, Failed, Retry string
}

// MessageStatus shows a sending clock or a failed retry.
func MessageStatus(c *ui.Context, state SendState, retry func(), o MessageStatusOptions) ui.Element {
	if state >= sendStateCount || retry == nil {
		panic(fmt.Sprintf("chat: MessageStatus state %d needs a retry func", state))
	}
	def(&o.Sending, "Sending")
	def(&o.Failed, "Failed")
	def(&o.Retry, "Retry")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleStatus)
	if state == Sending {
		root.Label(o.Sending).Tooltip(o.Sending).Children(func() {
			ic := ui.Icon(c, icon.Must("clock-3")).Size(16, 16).TextColor(t.TextMuted)
			ic.Rotate(360 * float32(ic.Loop("spin", 2*time.Second, ui.Linear)))
		})
		return root
	}
	root.Label(o.Failed).Children(func() {
		ui.Icon(c, icon.Must("triangle-alert")).Size(16, 16).TextColor(t.Danger)
		btn := kit.Press(ui.ButtonBase(c).Label(o.Retry).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(3)).Radius(999).TextColor(t.Danger))
		btn.Background(btn.AnimateColor("bg", t.Danger.Alpha(0.12*unit(btn.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("rotate-cw")).Size(16, 16)
			ui.Text(c, o.Retry).FontSize(t.Rem(0.9)).Bold().SingleLine()
		})
		if btn.Clicked() {
			retry()
			root.SetChanged()
			sound.Play(sound.Send)
			kit.Log("MessageStatus", "retried")
		}
	})
	return root
}
