package marketing

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NewsletterOptions word the signup band.
type NewsletterOptions struct{ Title, Text, Label, Send, Sent, Note string }

// Newsletter signs up a valid email; cue Success.
func Newsletter(c *ui.Context, email *string, sent *bool, o NewsletterOptions) (root, send ui.Element) {
	def(&o.Title, "Get the monthly letter")
	if email == nil || sent == nil {
		panic(fmt.Sprintf("marketing: Newsletter %q needs an email and a sent state", o.Title))
	}
	def(&o.Label, "Email address")
	def(&o.Send, "Subscribe")
	def(&o.Sent, "You're on the list")
	t := c.Theme()
	root = card(c).Gap(t.Space(3)).Padding(t.Space(6), t.Space(5)).Background(t.Background.Mix(t.Accent, 0.04)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "mail", t.Accent, t.Space(11))
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, o.Title).FontSize(t.Rem(1.05)).Bold().FillWidth().Wrap()
				if o.Text != "" {
					kit.Muted(c, o.Text).FontSize(t.Rem(0.8)).FillWidth().Wrap()
				}
			})
		})
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			in := emailField(c, o.Label, email)
			if in.Changed() {
				*sent = false
				root.SetChanged()
				kit.Log("Newsletter", "edited", "runes", utf8.RuneCountInString(*email))
			}
			ok := validEmail(*email) && !*sent
			c.Key("send")
			send = action(c, o.Send, "arrow-right", true, t.Space(60)).Disabled(!ok)
			if (send.Clicked() || in.Submitted()) && ok {
				*sent = true
				root.SetChanged()
				sound.Play(sound.Success)
				kit.Log("Newsletter", "subscribed")
			}
		})
		ui.Column(c).FillWidth().Children(func() {
			if *sent {
				c.Key("sent")
				note := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).TextColor(t.Success)
				kit.Arrive(note, ui.Pose{Y: t.Space(1.5), Fade: 1}, 0).Children(func() {
					glyph(c, "circle-check", 16)
					ui.Text(c, o.Sent).FontSize(t.Rem(0.82)).Bold().Grow(1).MinWidth(0).Wrap()
				})
			} else if o.Note != "" {
				kit.Muted(c, o.Note).FontSize(t.Rem(0.74)).FillWidth().Wrap()
			}
		})
	})
	return root, send
}
