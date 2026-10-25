package account

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TwoFactorOptions word the verification card.
type TwoFactorOptions struct{ Title, Subtitle, Code, Resend, Verify, Verified string }

// TwoFactor verifies six digits with resend; cue Success.
func TwoFactor(c *ui.Context, code *string, resent *int, verified *bool, o TwoFactorOptions) (root, verify ui.Element) {
	if code == nil || resent == nil || *resent < 0 || verified == nil {
		panic("account: TwoFactor needs code and state")
	}
	def(&o.Title, "Two-step verification")
	def(&o.Subtitle, "Enter the code from your authenticator app.")
	def(&o.Code, "Six-digit code")
	def(&o.Resend, "Resend")
	def(&o.Verify, "Verify")
	def(&o.Verified, "Verified")
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		head(c, "smartphone", o.Title, o.Subtitle)
		ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
			in := field(c, o.Code, "key-round", code, false)
			if in.Changed() {
				root.SetChanged()
			}
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				for i := range 6 {
					filled := i < len(*code) && (*code)[i] >= '0' && (*code)[i] <= '9'
					pip := ui.Box(c).Grow(1).Height(t.Space(1)).Radius(999)
					pip.Background(pip.AnimateColor("ink", t.SurfaceHover.Mix(t.Accent, unit(filled)), ui.Snappy))
				}
			})
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			resend := outline(c, o.Resend, o.Resend, "rotate-cw", t.Text).Disabled(*resent == math.MaxInt)
			if resend.Clicked() && *resent < math.MaxInt {
				*resent++
				root.SetChanged()
				sound.Play(sound.Send)
				kit.Log("TwoFactor", "resent", "count", *resent)
			}
			valid := sixDigits(*code)
			verify = primary(c, o.Verify, "shield-check").Grow(1).Disabled(!valid || *verified)
			if verify.Clicked() && valid && !*verified {
				*verified = true
				verify.SetChanged()
				root.SetChanged()
				sound.Play(sound.Success)
				kit.Log("TwoFactor", "verified")
			}
		})
		if *verified {
			c.Key("verified")
			note(c, "circle-check", o.Verified, t.Success)
		}
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
