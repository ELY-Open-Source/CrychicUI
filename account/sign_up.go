package account

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SignUpValue is one registration form.
type SignUpValue struct {
	Name, Email, Password string
	Terms                 bool
}

// SignUpOptions are the registration labels.
type SignUpOptions struct{ Name, Email, Password, Terms, Submit string }

// SignUp registers with terms; cue Success.
func SignUp(c *ui.Context, value *SignUpValue, o SignUpOptions) (root, submit ui.Element) {
	if value == nil {
		panic("account: SignUp needs a value")
	}
	def(&o.Name, "Name")
	def(&o.Email, "Email")
	def(&o.Password, "Password")
	def(&o.Terms, "Accept terms")
	def(&o.Submit, "Create account")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		name := ui.TextInput(c, &value.Name).Label(o.Name).Placeholder(o.Name).FillWidth()
		email := ui.TextInput(c, &value.Email).Label(o.Email).Placeholder(o.Email).FillWidth()
		password := ui.TextInput(c, &value.Password).Label(o.Password).Placeholder(o.Password).Password().FillWidth()
		if name.Changed() || email.Changed() || password.Changed() {
			root.SetChanged()
		}
		terms := button(c, o.Terms).Checked(value.Terms)
		terms.Children(func() { ui.Text(c, o.Terms) })
		if terms.Clicked() {
			value.Terms = !value.Terms
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("SignUp", "terms", value.Terms)
		}
		valid := strings.TrimSpace(value.Name) != "" && validEmail(value.Email) && len([]rune(value.Password)) >= 8 && value.Terms
		submit = primary(c, o.Submit).Disabled(!valid)
		submit.Children(func() { ui.Text(c, o.Submit) })
		if submit.Clicked() && valid {
			submit.SetChanged()
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("SignUp", "created")
		}
	})
	return root, submit
}
