package account

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SignUpValue is one registration form.
type SignUpValue struct {
	Name, Email, Password string
	Terms                 bool
}

// SignUpOptions word the form.
type SignUpOptions struct{ Title, Subtitle, Name, Email, Password, Rule, Terms, Submit string }

// SignUp registers with terms; cue Success.
func SignUp(c *ui.Context, value *SignUpValue, o SignUpOptions) (root, submit ui.Element) {
	if value == nil {
		panic("account: SignUp needs a value")
	}
	def(&o.Title, "Create your account")
	def(&o.Subtitle, "Start your 14-day free trial")
	def(&o.Name, "Name")
	def(&o.Email, "Email")
	def(&o.Password, "Password")
	def(&o.Rule, "At least 8 characters")
	def(&o.Terms, "Accept terms")
	def(&o.Submit, "Create account")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		head(c, "user-plus", o.Title, o.Subtitle)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
			name := field(c, o.Name, "user", &value.Name, false)
			email := field(c, o.Email, "mail", &value.Email, false)
			password := field(c, o.Password, "lock", &value.Password, true)
			if name.Changed() || email.Changed() || password.Changed() {
				root.SetChanged()
			}
			long := utf8.RuneCountInString(value.Password) >= 8
			rule := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center)
			rule.TextColor(rule.AnimateColor("ink", t.TextMuted.Mix(t.Success, unit(long)), ui.Snappy)).Children(func() {
				c.Key(long)
				mark := "circle"
				if long {
					mark = "circle-check"
				}
				kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.5, Fade: 1}, 0)
				ui.Text(c, o.Rule).FontSize(t.Rem(0.75)).Shrink(0)
			})
		})
		terms := kit.Press(ui.ButtonBase(c).Label(o.Terms).AlignSelf(ui.Start).Padding(t.Space(1), 0).Radius(t.Radius).Checked(value.Terms))
		terms.Children(func() {
			ui.Row(c).Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
				tick(c, value.Terms, terms.Hovered())
				ui.Text(c, o.Terms).FontSize(t.Rem(0.82)).Shrink(0)
			})
		})
		if terms.Clicked() {
			value.Terms = !value.Terms
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("SignUp", "terms", value.Terms)
		}
		valid := strings.TrimSpace(value.Name) != "" && validEmail(value.Email) && utf8.RuneCountInString(value.Password) >= 8 && value.Terms
		submit = primary(c, o.Submit, "").FillWidth().Disabled(!valid)
		if submit.Clicked() && valid {
			submit.SetChanged()
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("SignUp", "created")
		}
	})
	return root, submit
}
