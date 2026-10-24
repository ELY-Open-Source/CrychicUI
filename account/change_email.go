package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChangeEmailOptions are the email change labels.
type ChangeEmailOptions struct{ Email, Verify string }

// ChangeEmail verifies a new email; cue Send.
func ChangeEmail(c *ui.Context, email *string, o ChangeEmailOptions) (root, verify ui.Element) {
	if email == nil {
		panic("account: ChangeEmail needs email")
	}
	def(&o.Email, "New email")
	def(&o.Verify, "Verify")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		in := ui.TextInput(c, email).Label(o.Email).Placeholder(o.Email).Grow(1)
		if in.Changed() {
			root.SetChanged()
		}
		verify = primary(c, o.Verify).Disabled(!validEmail(*email))
		verify.Children(func() { ui.Text(c, o.Verify) })
	})
	if verify.Clicked() && validEmail(*email) {
		verify.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("ChangeEmail", "verify")
	}
	return root, verify
}
