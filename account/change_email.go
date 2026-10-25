package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChangeEmailOptions word the card; Current shows the present address.
type ChangeEmailOptions struct{ Title, Current, Email, Verify, Sent string }

// ChangeEmail verifies a new email; cue Send.
func ChangeEmail(c *ui.Context, email *string, o ChangeEmailOptions) (root, verify ui.Element) {
	if email == nil {
		panic("account: ChangeEmail needs email")
	}
	def(&o.Title, "Change email")
	def(&o.Email, "New email")
	def(&o.Verify, "Verify")
	def(&o.Sent, "We sent a confirmation link to %s.")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	sent := ui.Local(root, "sent", func() string { return "" })
	root.Children(func() {
		head(c, "mail", o.Title, o.Current)
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				in := field(c, o.Email, "at-sign", email, false)
				if in.Changed() {
					root.SetChanged()
				}
			})
			verify = primary(c, o.Verify, "send").Disabled(!validEmail(*email))
		})
		if *sent != "" && *sent == *email {
			c.Key("sent")
			note(c, "mail-check", fmt.Sprintf(o.Sent, *sent), t.Success)
		}
	})
	if verify.Clicked() && validEmail(*email) {
		*sent = *email
		verify.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("ChangeEmail", "verify")
	}
	return root, verify
}
