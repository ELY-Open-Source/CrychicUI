package shop

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CouponOptions are the coupon labels.
type CouponOptions struct{ Placeholder, Apply, Invalid string }

// Coupon applies a valid code or shows an error.
func Coupon(c *ui.Context, code *string, codes []string, applied *string, o CouponOptions) (root, apply ui.Element) {
	if code == nil || applied == nil || len(codes) == 0 {
		panic("shop: Coupon needs a code and valid codes")
	}
	def(&o.Placeholder, "Coupon code")
	def(&o.Apply, "Apply")
	def(&o.Invalid, "Invalid coupon")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	bad := ui.Local(root, "bad", func() bool { return false })
	var in ui.Element
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			in = ui.TextInput(c, code).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1)
			apply = addButton(c, o.Apply).Disabled(strings.TrimSpace(*code) == "")
			apply.Children(func() { ui.Text(c, o.Apply).SingleLine() })
		})
		if *bad {
			ui.Text(c, o.Invalid).FontSize(t.Rem(0.8)).TextColor(t.Danger)
		}
	})
	if in.Changed() {
		*bad = false
		root.SetChanged()
	}
	if apply.Clicked() || in.Submitted() {
		*bad = true
		for _, valid := range codes {
			if strings.EqualFold(strings.TrimSpace(*code), strings.TrimSpace(valid)) {
				*applied = strings.TrimSpace(*code)
				*bad = false
			}
		}
		apply.SetChanged()
		root.SetChanged()
		if *bad {
			sound.Play(sound.Error)
			kit.Log("Coupon", "rejected")
		} else {
			sound.Play(sound.Success)
			kit.Log("Coupon", "applied")
		}
	}
	return root, apply
}
