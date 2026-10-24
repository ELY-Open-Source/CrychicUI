package shop

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CheckoutStep is one checkout stage.
type CheckoutStep uint8

// Checkout stages.
const (
	Cart CheckoutStep = iota
	Shipping
	Payment
	Done
)

// CheckoutStepsOptions name the stages.
type CheckoutStepsOptions struct{ Cart, Shipping, Payment, Done string }

// CheckoutSteps steps through checkout, done stages checked; cue Tick.
func CheckoutSteps(c *ui.Context, current *CheckoutStep, o CheckoutStepsOptions) ui.Element {
	if current == nil || *current > Done {
		panic("shop: CheckoutSteps needs a known step")
	}
	def(&o.Cart, "Cart")
	def(&o.Shipping, "Shipping")
	def(&o.Payment, "Payment")
	def(&o.Done, "Done")
	t := c.Theme()
	names := [...]string{o.Cart, o.Shipping, o.Payment, o.Done}
	root := ui.Row(c).FillWidth().AlignItems(ui.Start).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(names[*current])
	root.Children(func() {
		for i, name := range names {
			step := CheckoutStep(i)
			c.Key(step)
			if i > 0 {
				line := ui.Box(c).Grow(1).Height(t.Space(0.5)).Margin(t.Space(4), t.Space(1), 0, t.Space(1)).Radius(999).Background(t.Border)
				fill := line.AnimateTo("fill", unit(step <= *current), ui.Smooth)
				line.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, t.Accent, 999)
				}))
			}
			done, on := step < *current, step == *current
			b := kit.Press(ui.ButtonBase(c).Label(name).Column().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1)).Radius(t.Radius).Checked(on))
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered() && !on)), ui.Snappy))
			b.Children(func() {
				dot := ui.Box(c).Size(t.Space(7), t.Space(7)).Center().Radius(999)
				dot.Background(dot.AnimateColor("bg", t.Surface.Mix(t.Accent, unit(done || on)), ui.Snappy))
				dot.TextColor(dot.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(done || on)), ui.Snappy))
				ring := dot.AnimateTo("ring", unit(on), ui.Snappy)
				dot.Border(t.Space(0.75)*ring, t.Accent.Alpha(0.25))
				dot.Children(func() {
					if done {
						ui.Icon(c, icon.Must("check")).Size(14, 14)
					} else {
						ui.Text(c, strconv.Itoa(i+1)).FontSize(t.Rem(0.75)).Bold()
					}
				})
				txt := ui.Text(c, name).FontSize(t.Rem(0.75)).Shrink(0)
				if on {
					txt.Bold()
				} else {
					txt.TextColor(t.TextMuted)
				}
			})
			if b.Clicked() && !on {
				*current = step
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("CheckoutSteps", "step", "value", int(step))
			}
		}
	})
	return root
}
