package account

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TwoFactorOptions are the verification labels.
type TwoFactorOptions struct{ Code, Resend, Verify string }

// TwoFactor verifies six digits with resend; cue Success.
func TwoFactor(c *ui.Context, code *string, resent *int, verified *bool, o TwoFactorOptions) (root, verify ui.Element) {
	if code == nil || resent == nil || *resent < 0 || verified == nil {
		panic("account: TwoFactor needs code and state")
	}
	def(&o.Code, "Six-digit code")
	def(&o.Resend, "Resend")
	def(&o.Verify, "Verify")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		in := ui.TextInput(c, code).Label(o.Code).Placeholder(o.Code).FillWidth()
		if in.Changed() {
			root.SetChanged()
		}
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			resend := button(c, o.Resend).Disabled(*resent == math.MaxInt)
			resend.Children(func() { ui.Text(c, o.Resend) })
			if resend.Clicked() && *resent < math.MaxInt {
				*resent++
				root.SetChanged()
				sound.Play(sound.Send)
				kit.Log("TwoFactor", "resent", "count", *resent)
			}
			valid := sixDigits(*code)
			verify = primary(c, o.Verify).Disabled(!valid || *verified)
			verify.Children(func() { ui.Text(c, o.Verify) })
			if verify.Clicked() && valid && !*verified {
				*verified = true
				verify.SetChanged()
				root.SetChanged()
				sound.Play(sound.Success)
				kit.Log("TwoFactor", "verified")
			}
		})
	})
	return root, verify
}

func sixDigits(s string) bool {
	if len(s) != 6 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
