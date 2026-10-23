package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CheckoutStep is one checkout stage.
type CheckoutStep uint8

const (
	Cart CheckoutStep = iota
	Shipping
	Payment
	Done
)

// CheckoutStepsOptions are the checkout labels.
type CheckoutStepsOptions struct{ Cart, Shipping, Payment, Done string }

// CheckoutSteps selects checkout stages; cue Tick.
func CheckoutSteps(c *ui.Context, current *CheckoutStep, o CheckoutStepsOptions) ui.Element {
	if current == nil || *current > Done {
		panic("shop: CheckoutSteps needs a valid step")
	}
	def(&o.Cart, "Cart")
	def(&o.Shipping, "Shipping")
	def(&o.Payment, "Payment")
	def(&o.Done, "Done")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		for _, step := range []struct {
			value CheckoutStep
			name  string
		}{{Cart, o.Cart}, {Shipping, o.Shipping}, {Payment, o.Payment}, {Done, o.Done}} {
			b := kit.Press(ui.ButtonBase(c).Label(step.name).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(*current == step.value))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*current == step.value)), ui.Snappy))
			b.Children(func() { ui.Text(c, step.name).FontSize(t.Rem(0.8)) })
			if b.Clicked() && *current != step.value {
				*current = step.value
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("CheckoutSteps", "step", "value", *current)
			}
		}
	})
	return root
}
