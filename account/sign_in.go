package account

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Credentials is one sign-in form.
type Credentials struct {
	Email, Password string
}

// SignInOptions are the sign-in labels and error.
type SignInOptions struct {
	Email, Password, Forgot, Submit string
	Error                           string
}

// SignIn signs in with providers and forgot; cue Tap, Error.
func SignIn(c *ui.Context, value *Credentials, providers []string, action *string, o SignInOptions) (root, submit ui.Element) {
	if value == nil || action == nil {
		panic("account: SignIn needs value and action")
	}
	def(&o.Email, "Email")
	def(&o.Password, "Password")
	def(&o.Forgot, "Forgot password")
	def(&o.Submit, "Sign in")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	shown := ui.Local(root, "error", func() string { return o.Error })
	if o.Error != *shown {
		if *shown = o.Error; o.Error != "" {
			sound.Play(sound.Error)
			kit.Log("SignIn", "error")
		}
	}
	valid := false
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
			email := ui.TextInput(c, &value.Email).Label(o.Email).Placeholder(o.Email).FillWidth()
			password := ui.TextInput(c, &value.Password).Label(o.Password).Placeholder(o.Password).Password().FillWidth()
			if email.Changed() || password.Changed() {
				*action = ""
				root.SetChanged()
			}
			if o.Error != "" {
				ui.Text(c, o.Error).TextColor(t.Danger).Wrap()
			}
		})
		valid = validEmail(value.Email) && value.Password != ""
		submit = primary(c, o.Submit).FillWidth().Disabled(!valid)
		submit.Children(func() { ui.Text(c, o.Submit) })
		ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for _, provider := range providers {
				if strings.TrimSpace(provider) == "" {
					panic("account: provider must not be empty")
				}
				c.Key(provider)
				b := button(c, provider).FillWidth()
				b.Children(func() { ui.Text(c, provider) })
				if b.Clicked() {
					*action = provider
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("SignIn", "provider")
				}
			}
		})
		forgot := button(c, o.Forgot)
		forgot.Children(func() { ui.Text(c, o.Forgot) })
		if forgot.Clicked() {
			*action = "forgot"
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("SignIn", "forgot")
		}
	})
	if submit.Clicked() && valid {
		*action = "signin"
		submit.SetChanged()
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("SignIn", "submit")
	}
	return root, submit
}
