package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ForgotPasswordOptions word the recovery form.
type ForgotPasswordOptions struct{ Title, Subtitle, Email, Send, Sent string }

// ForgotPassword sends a recovery link; cue Send.
func ForgotPassword(c *ui.Context, email *string, o ForgotPasswordOptions) (root, send ui.Element) {
	if email == nil {
		panic("account: ForgotPassword needs email")
	}
	def(&o.Title, "Forgot your password?")
	def(&o.Subtitle, "We will email you a link to reset it.")
	def(&o.Email, "Email")
	def(&o.Send, "Send link")
	def(&o.Sent, "Link sent. Check your inbox.")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	sent := ui.Local(root, "sent", func() bool { return false })
	root.Children(func() {
		head(c, "key-round", o.Title, o.Subtitle)
		in := field(c, o.Email, "mail", email, false)
		if in.Changed() {
			*sent = false
			root.SetChanged()
		}
		send = primary(c, o.Send, "send").FillWidth().Disabled(!validEmail(*email))
		if *sent {
			c.Key("sent")
			note(c, "mail-check", o.Sent, t.Success)
		}
	})
	if send.Clicked() && validEmail(*email) {
		*sent = true
		send.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("ForgotPassword", "sent")
	}
	return root, send
}
