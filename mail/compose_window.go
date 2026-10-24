package mail

import (
	"net/mail"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Compose is one outgoing message.
type Compose struct {
	To, Subject, Body string
}

// ComposeWindowOptions are the compose labels.
type ComposeWindowOptions struct{ To, Subject, Body, Send string }

// ComposeWindow edits and sends mail; cue Send.
func ComposeWindow(c *ui.Context, value *Compose, o ComposeWindowOptions) (root, send ui.Element) {
	if value == nil {
		panic("mail: ComposeWindow needs a value")
	}
	def(&o.To, "To")
	def(&o.Subject, "Subject")
	def(&o.Body, "Message")
	def(&o.Send, "Send")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		to := ui.TextInput(c, &value.To).Label(o.To).Placeholder(o.To).FillWidth()
		subject := ui.TextInput(c, &value.Subject).Label(o.Subject).Placeholder(o.Subject).FillWidth()
		body := ui.TextArea(c, &value.Body).Label(o.Body).Placeholder(o.Body).FillWidth().Lines(4, 10)
		if to.Changed() || subject.Changed() || body.Changed() {
			root.SetChanged()
		}
		_, err := mail.ParseAddress(value.To)
		valid := err == nil && strings.TrimSpace(value.Subject) != "" && strings.TrimSpace(value.Body) != ""
		send = button(c, o.Send).Background(t.Accent).TextColor(t.AccentText).Disabled(!valid)
		send.Children(func() { ui.Text(c, o.Send) })
	})
	if send.Clicked() {
		send.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("ComposeWindow", "sent", "bodyRunes", len([]rune(value.Body)))
	}
	return root, send
}
