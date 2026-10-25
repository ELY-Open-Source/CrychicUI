package account

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PasswordChange is one password change form.
type PasswordChange struct {
	Current, New string
}

// ChangePasswordOptions word the form and its two rules.
type ChangePasswordOptions struct{ Title, Current, New, Submit, Long, Different, Changed string }

// ChangePassword changes the current password; cue Success.
func ChangePassword(c *ui.Context, value *PasswordChange, o ChangePasswordOptions) (root, submit ui.Element) {
	if value == nil {
		panic("account: ChangePassword needs a value")
	}
	def(&o.Title, "Change password")
	def(&o.Current, "Current password")
	def(&o.New, "New password")
	def(&o.Submit, "Update password")
	def(&o.Long, "At least 8 characters")
	def(&o.Different, "Different from the current one")
	def(&o.Changed, "Password updated")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	done := ui.Local(root, "done", func() bool { return false })
	root.Children(func() {
		head(c, "lock-keyhole", o.Title, "")
		ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
			current := field(c, o.Current, "lock", &value.Current, true)
			newPassword := field(c, o.New, "lock-keyhole", &value.New, true)
			if current.Changed() || newPassword.Changed() {
				*done = false
				root.SetChanged()
			}
		})
		long := utf8.RuneCountInString(value.New) >= 8
		different := value.New != "" && value.New != value.Current
		ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
			for i, rule := range []struct {
				text string
				ok   bool
			}{{o.Long, long}, {o.Different, different}} {
				c.Key(i)
				row := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center)
				row.TextColor(row.AnimateColor("ink", t.TextMuted.Mix(t.Success, unit(rule.ok)), ui.Snappy)).Children(func() {
					c.Key(rule.ok)
					mark := "circle"
					if rule.ok {
						mark = "circle-check"
					}
					kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.5, Fade: 1}, 0)
					ui.Text(c, rule.text).FontSize(t.Rem(0.78)).Wrap()
				})
			}
		})
		valid := value.Current != "" && long && different
		submit = primary(c, o.Submit, "").FillWidth().Disabled(!valid)
		if *done {
			c.Key("done")
			note(c, "circle-check", o.Changed, t.Success)
		}
		if submit.Clicked() && valid {
			*done = true
			submit.SetChanged()
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("ChangePassword", "changed")
		}
	})
	return root, submit
}
