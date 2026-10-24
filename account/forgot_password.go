package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ForgotPasswordOptions are the recovery labels.
type ForgotPasswordOptions struct{ Email, Send string }

// ForgotPassword sends a recovery link; cue Send.
func ForgotPassword(c *ui.Context, email *string, o ForgotPasswordOptions) (root, send ui.Element) {
	if email == nil {
		panic("account: ForgotPassword needs email")
	}
	def(&o.Email, "Email")
	def(&o.Send, "Send link")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		in := ui.TextInput(c, email).Label(o.Email).Placeholder(o.Email).FillWidth()
		if in.Changed() {
			root.SetChanged()
		}
		send = primary(c, o.Send).Disabled(!validEmail(*email))
		send.Children(func() { ui.Text(c, o.Send) })
	})
	if send.Clicked() && validEmail(*email) {
		send.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("ForgotPassword", "sent")
	}
	return root, send
}
