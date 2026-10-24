package shop

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// shakeTime is how long a refusal pushes the field.
const shakeTime = 60 * time.Millisecond

// CouponOptions word the coupon field; Applied takes the code.
type CouponOptions struct{ Placeholder, Apply, Invalid, Applied string }

// Coupon applies a valid code or shakes; cue Success, Error.
func Coupon(c *ui.Context, code *string, codes []string, applied *string, o CouponOptions) (root, apply ui.Element) {
	if code == nil || applied == nil || len(codes) == 0 {
		panic("shop: Coupon needs a code, the valid codes and bound state")
	}
	def(&o.Placeholder, "Coupon code")
	def(&o.Apply, "Apply")
	def(&o.Invalid, "Invalid coupon")
	def(&o.Applied, "%s applied")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	bad := ui.Local(root, "bad", func() bool { return false })
	shook := ui.Local(root, "shook", func() time.Time { return time.Time{} })
	kick := c.Now().Sub(*shook) < shakeTime
	if kick {
		c.After(shakeTime - c.Now().Sub(*shook))
	}
	var in ui.Element
	root.Children(func() {
		field := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).Radius(t.Radius)
		field.Border(1, field.AnimateColor("line", t.Border.Mix(t.Danger, unit(*bad)), ui.Snappy))
		// A brief push sprung back by Bouncy makes the shake.
		field.Translate(field.AnimateTo("shake", t.Space(3)*unit(kick), ui.Bouncy), 0).Children(func() {
			ui.Icon(c, icon.Must("ticket-percent")).Size(16, 16).Shrink(0).TextColor(t.TextMuted)
			in = ui.TextInput(c, code).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1)
			apply = kit.Press(ui.ButtonBase(c).Label(o.Apply).Padding(t.Space(1.25), t.Space(3)).Radius(t.Radius)).Disabled(strings.TrimSpace(*code) == "")
			apply.Background(apply.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(apply.Hovered())), ui.Snappy)).TextColor(t.AccentText).
				Children(func() { ui.Text(c, o.Apply).Bold().Shrink(0) })
		})
		switch {
		case *bad:
			ui.Row(c).Key("bad").Gap(t.Space(1.5)).AlignItems(ui.Center).Transition(kit.Enter).Children(func() {
				ui.Icon(c, icon.Must("circle-x")).Size(12, 12).TextColor(t.Danger)
				ui.Text(c, o.Invalid).FontSize(t.Rem(0.8)).TextColor(t.Danger)
			})
		case *applied != "":
			ui.Row(c).Key("applied").Gap(t.Space(1.5)).AlignItems(ui.Center).Transition(kit.Enter).Children(func() {
				ui.Icon(c, icon.Must("circle-check")).Size(12, 12).TextColor(t.Success)
				ui.Text(c, fmt.Sprintf(o.Applied, *applied)).FontSize(t.Rem(0.8)).Bold().TextColor(t.Success)
			})
		}
	})
	if in.Changed() {
		*bad = false
		root.SetChanged()
	}
	if apply.Clicked() || in.Submitted() {
		entered := strings.TrimSpace(*code)
		*bad = true
		for _, valid := range codes {
			if entered != "" && strings.EqualFold(entered, strings.TrimSpace(valid)) {
				*applied, *bad = entered, false
			}
		}
		apply.SetChanged()
		root.SetChanged()
		if *bad {
			*shook = c.Now()
			c.After(shakeTime)
			sound.Play(sound.Error)
			kit.Log("Coupon", "rejected")
		} else {
			sound.Play(sound.Success)
			kit.Log("Coupon", "applied")
		}
	}
	return root, apply
}
