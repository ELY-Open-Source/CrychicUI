package mail

import (
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Compose is one outgoing message.
type Compose struct {
	To, Subject, Body string
}

// ComposeWindowOptions word the window and its fields.
type ComposeWindowOptions struct{ Title, To, Subject, Body, Send string }

// ComposeWindow edits and sends mail; cue Send.
func ComposeWindow(c *ui.Context, value *Compose, o ComposeWindowOptions) (root, send ui.Element) {
	if value == nil {
		panic("mail: ComposeWindow needs a value")
	}
	def(&o.Title, "New message")
	def(&o.To, "To")
	def(&o.Subject, "Subject")
	def(&o.Body, "Message")
	def(&o.Send, "Send")
	t := c.Theme()
	root = card(c)
	root.Children(func() {
		kit.CardHead(c, "pen-line", o.Title, nil)
		field := func(label string, text *string) ui.Element {
			var in ui.Element
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
				BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
				kit.Muted(c, label).FontSize(t.Rem(0.8)).Width(t.Space(14)).Shrink(0).SingleLine()
				in = ui.TextInputBase(c, text).Label(label).Placeholder(label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0)
			})
			return in
		}
		to := field(o.To, &value.To)
		subject := field(o.Subject, &value.Subject)
		body := ui.TextAreaBase(c, &value.Body).Label(o.Body).Placeholder(o.Body).FontSize(t.Rem(0.85)).FillWidth().
			Lines(5, 10).Padding(t.Space(2.5), t.Space(3))
		if to.Changed() || subject.Changed() || body.Changed() {
			root.SetChanged()
		}
		_, err := mail.ParseAddress(value.To)
		valid := err == nil && strings.TrimSpace(value.Subject) != "" && strings.TrimSpace(value.Body) != ""
		ui.Row(c).FillWidth().Justify(ui.End).Padding(t.Space(2), t.Space(3)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			send = kit.Press(ui.ButtonBase(c).Label(o.Send).Padding(t.Space(2), t.Space(4)).Radius(999).TextColor(t.AccentText).Disabled(!valid))
			send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy))
			send.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, o.Send).FontSize(t.Rem(0.85)).Bold().Shrink(0)
					glyph(c, "send", 14)
				})
			})
		})
	})
	if send.Clicked() {
		send.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("ComposeWindow", "sent", "bodyRunes", utf8.RuneCountInString(value.Body))
	}
	return root, send
}
