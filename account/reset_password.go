package account

import (
	"unicode"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ResetValue is one password reset form.
type ResetValue struct {
	Password, Confirm string
}

// ResetPasswordOptions word the reset form.
type ResetPasswordOptions struct{ Title, Subtitle, Password, Confirm, Submit, Weak, Strong, Mismatch string }

// ResetPassword resets with strength and confirmation; cue Success.
func ResetPassword(c *ui.Context, value *ResetValue, o ResetPasswordOptions) (root, submit ui.Element) {
	if value == nil {
		panic("account: ResetPassword needs a value")
	}
	def(&o.Title, "Choose a new password")
	def(&o.Subtitle, "Mix letters, numbers and symbols.")
	def(&o.Password, "New password")
	def(&o.Confirm, "Confirm password")
	def(&o.Submit, "Reset password")
	def(&o.Weak, "Weak")
	def(&o.Strong, "Strong")
	def(&o.Mismatch, "Passwords do not match")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		head(c, "shield-check", o.Title, o.Subtitle)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
			password := field(c, o.Password, "lock", &value.Password, true)
			score := passwordStrength(value.Password)
			long := utf8.RuneCountInString(value.Password) >= 8
			label, ink := o.Weak, t.Danger
			if score >= 2 {
				ink = t.Warning
			}
			if long && score >= 3 {
				label, ink = o.Strong, t.Success
			}
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Row(c).Grow(1).Gap(t.Space(1)).Children(func() {
					for i := range 4 {
						seg := ui.Box(c).Grow(1).Height(t.Space(1.25)).Radius(999)
						seg.Background(seg.AnimateColor("ink", t.SurfaceHover.Mix(ink, unit(i < score)), ui.Snappy))
					}
				})
				text := ui.Text(c, label).FontSize(t.Rem(0.75)).Bold().Shrink(0)
				text.TextColor(text.AnimateColor("ink", ink, ui.Snappy))
			})
			confirm := field(c, o.Confirm, "lock-keyhole", &value.Confirm, true)
			if password.Changed() || confirm.Changed() {
				root.SetChanged()
			}
			if value.Confirm != "" && value.Confirm != value.Password {
				c.Key("mismatch")
				note(c, "circle-alert", o.Mismatch, t.Danger)
			}
		})
		score := passwordStrength(value.Password)
		valid := utf8.RuneCountInString(value.Password) >= 8 && score >= 2 && value.Password == value.Confirm
		submit = primary(c, o.Submit, "").FillWidth().Disabled(!valid)
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
	if utf8.RuneCountInString(s) >= 8 {
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
