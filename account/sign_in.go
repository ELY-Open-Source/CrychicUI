package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Credentials is one sign-in form.
type Credentials struct {
	Email, Password string
}

// SignInOptions word the form; Error shows a failed attempt.
type SignInOptions struct {
	Title, Subtitle, Email, Password, Forgot, Submit, Or, Continue string
	Error                                                          string
}

// SignIn signs in with providers and forgot; cue Tap, Error.
func SignIn(c *ui.Context, value *Credentials, providers []string, action *string, o SignInOptions) (root, submit ui.Element) {
	if value == nil || action == nil {
		panic("account: SignIn needs value and action")
	}
	if len(providers) > 0 {
		unique("SignIn", providers)
	}
	def(&o.Title, "Welcome back")
	def(&o.Subtitle, "Sign in to your account")
	def(&o.Email, "Email")
	def(&o.Password, "Password")
	def(&o.Forgot, "Forgot password")
	def(&o.Submit, "Sign in")
	def(&o.Or, "or")
	def(&o.Continue, "Continue with %s")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5))
	shown := ui.Local(root, "error", func() string { return o.Error })
	if o.Error != *shown {
		if *shown = o.Error; o.Error != "" {
			sound.Play(sound.Error)
			kit.Log("SignIn", "error")
		}
	}
	valid := false
	root.Children(func() {
		head(c, "log-in", o.Title, o.Subtitle)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
			email := field(c, o.Email, "mail", &value.Email, false)
			password := field(c, o.Password, "lock", &value.Password, true)
			if email.Changed() || password.Changed() {
				*action = ""
				root.SetChanged()
			}
			forgot := kit.Press(ui.ButtonBase(c).Label(o.Forgot).AlignSelf(ui.End).Padding(t.Space(0.5), t.Space(1)).Radius(t.Radius))
			forgot.TextColor(forgot.AnimateColor("ink", t.Accent.Mix(t.AccentHover, unit(forgot.Hovered())), ui.Snappy))
			forgot.Children(func() { ui.Text(c, o.Forgot).FontSize(t.Rem(0.78)).Bold().Shrink(0) })
			if forgot.Clicked() {
				*action = "forgot"
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("SignIn", "forgot")
			}
			if o.Error != "" {
				c.Key(o.Error)
				note(c, "circle-alert", o.Error, t.Danger)
			}
		})
		valid = validEmail(value.Email) && value.Password != ""
		submit = primary(c, o.Submit, "").FillWidth().Disabled(!valid)
		if len(providers) == 0 {
			return
		}
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Grow(1).Height(1).Background(t.Border)
			kit.Muted(c, o.Or).FontSize(t.Rem(0.75)).Shrink(0)
			ui.Box(c).Grow(1).Height(1).Background(t.Border)
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(2))
		list.Children(func() {
			for i, provider := range providers {
				c.Key(provider)
				b := providerButton(c, provider, o.Continue, *action == provider)
				kit.Stagger(c, list, b, i)
				if b.Clicked() {
					*action = provider
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("SignIn", "provider", "index", i)
				}
			}
		})
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
