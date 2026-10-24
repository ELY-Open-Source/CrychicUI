package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PasswordChange is one password change form.
type PasswordChange struct {
	Current, New string
}

// ChangePasswordOptions are the password change labels.
type ChangePasswordOptions struct{ Current, New, Submit string }

// ChangePassword changes the current password; cue Success.
func ChangePassword(c *ui.Context, value *PasswordChange, o ChangePasswordOptions) (root, submit ui.Element) {
	if value == nil {
		panic("account: ChangePassword needs a value")
	}
	def(&o.Current, "Current password")
	def(&o.New, "New password")
	def(&o.Submit, "Change password")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		current := ui.TextInput(c, &value.Current).Label(o.Current).Placeholder(o.Current).Password().FillWidth()
		newPassword := ui.TextInput(c, &value.New).Label(o.New).Placeholder(o.New).Password().FillWidth()
		if current.Changed() || newPassword.Changed() {
			root.SetChanged()
		}
		valid := value.Current != "" && len([]rune(value.New)) >= 8 && value.New != value.Current
		submit = primary(c, o.Submit).Disabled(!valid)
		submit.Children(func() { ui.Text(c, o.Submit) })
		if submit.Clicked() && valid {
			submit.SetChanged()
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("ChangePassword", "changed")
		}
	})
	return root, submit
}
