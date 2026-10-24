package account

import (
	"unicode"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ResetValue is one password reset form.
type ResetValue struct {
	Password, Confirm string
}

// ResetPasswordOptions are the reset labels.
type ResetPasswordOptions struct{ Password, Confirm, Submit, Weak, Strong string }

// ResetPassword resets with strength and confirmation; cue Success.
func ResetPassword(c *ui.Context, value *ResetValue, o ResetPasswordOptions) (root, submit ui.Element) {
	if value == nil {
		panic("account: ResetPassword needs a value")
	}
	def(&o.Password, "New password")
	def(&o.Confirm, "Confirm password")
	def(&o.Submit, "Reset password")
	def(&o.Weak, "Weak")
	def(&o.Strong, "Strong")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		password := ui.TextInput(c, &value.Password).Label(o.Password).Placeholder(o.Password).Password().FillWidth()
		confirm := ui.TextInput(c, &value.Confirm).Label(o.Confirm).Placeholder(o.Confirm).Password().FillWidth()
		if password.Changed() || confirm.Changed() {
			root.SetChanged()
		}
		score := passwordStrength(value.Password)
		longEnough := len([]rune(value.Password)) >= 8
		label, color := o.Weak, t.Warning
		if longEnough && score >= 3 {
			label, color = o.Strong, t.Success
		}
		ui.Text(c, label).TextColor(color)
		bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
		w := bar.AnimateTo("strength", float32(score)/4, ui.Smooth)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, color, 999)
		}))
		valid := longEnough && score >= 2 && value.Password == value.Confirm
		submit = primary(c, o.Submit).Disabled(!valid)
		submit.Children(func() { ui.Text(c, o.Submit) })
		if submit.Clicked() && valid {
			submit.SetChanged()
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("ResetPassword", "reset", "strength", score)
		}
	})
	return root, submit
}

func passwordStrength(s string) int {
	score := 0
	if len([]rune(s)) >= 8 {
		score++
	}
	classes := [4]bool{}
	for _, r := range s {
		switch {
		case unicode.IsLower(r):
			classes[0] = true
		case unicode.IsUpper(r):
			classes[1] = true
		case unicode.IsDigit(r):
			classes[2] = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			classes[3] = true
		}
	}
	for _, present := range classes {
		if present {
			score++
		}
	}
	return min(4, score)
}
