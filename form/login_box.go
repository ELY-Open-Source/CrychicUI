package form

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Login is the value a LoginBox edits.
type Login struct {
	Email, Password string
	Remember        bool
}

// LoginBoxOptions tune a LoginBox.
type LoginBoxOptions struct {
	// Email, Password, Remember and Submit are labels.
	Email, Password, Remember, Submit string
	// Error is the caller's sign-in failure; it shakes in.
	Error string
}

// LoginBox takes email, password and remember on a card.
func LoginBox(c *ui.Context, value *Login, o LoginBoxOptions) (root, submit ui.Element) {
	if value == nil {
		panic("form: LoginBox needs a value")
	}
	def(&o.Email, "Email")
	def(&o.Password, "Password")
	def(&o.Remember, "Remember me")
	def(&o.Submit, "Sign in")
	t := c.Theme()
	root = kit.Arrive(kit.Card(c).FillWidth().Gap(t.Space(4)).Padding(t.Space(6)), ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	shown := ui.Local(root, "error", func() string { return o.Error })
	if o.Error != *shown {
		if *shown = o.Error; o.Error != "" {
			sound.Play(sound.Error)
			kit.Log("LoginBox", "error")
		}
	}
	root.Children(func() {
		FieldGroup(c, o.Email, FieldGroupOptions{}, func() {
			edited(root, ui.TextInput(c, &value.Email).FillWidth())
		})
		FieldGroup(c, o.Password, FieldGroupOptions{}, func() {
			edited(root, ui.TextInput(c, &value.Password).Password().FillWidth())
		})
		edited(root, choice.Check(c, o.Remember, &value.Remember))
		feedback.FieldError(c, o.Error)
		ok := emailPattern.MatchString(value.Email) && value.Password != ""
		submit = primary(c, o.Submit).FillWidth().Disabled(!ok)
	})
	if kit.Clicked(submit, sound.Tap) {
		kit.Log("LoginBox", "submit", "remember", value.Remember)
	}
	return root, submit
}
