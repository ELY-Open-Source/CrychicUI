package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PaymentRequestOptions word the request; Requests takes the payer.
type PaymentRequestOptions struct{ Pay, Paid, Requests, Note string }

// PaymentRequest shows who asks for how much; cue Success.
func PaymentRequest(c *ui.Context, amount float64, payer, currency string, paid *bool, o PaymentRequestOptions) (root, pay ui.Element) {
	if amount <= 0 || !finite(amount) || payer == "" || currency == "" || paid == nil {
		panic(fmt.Sprintf("finance: PaymentRequest from %q needs a positive amount, a currency and bound state", payer))
	}
	def(&o.Pay, "Pay")
	def(&o.Paid, "Paid")
	def(&o.Requests, "%s requests")
	t := c.Theme()
	root = kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			display.Avatar(c, payer, nil).Size(t.Space(10), t.Space(10)).Shrink(0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, fmt.Sprintf(o.Requests, payer)).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
				if o.Note != "" {
					kit.Muted(c, o.Note).FontSize(t.Rem(0.8)).FillWidth().Wrap()
				}
			})
		})
		money(c, amount, currency).Bold().FontSize(t.Rem(2))
		name, glyph := o.Pay, "send-horizontal"
		if *paid {
			name, glyph = o.Paid, "circle-check"
		}
		pay = kit.Press(ui.ButtonBase(c).Label(name).FillWidth().Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).
			TextColor(t.AccentText).Disabled(*paid))
		pay.Background(pay.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(pay.Hovered())).Mix(t.Success, unit(*paid)), ui.Snappy))
		pay.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
				ui.Text(c, name).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
			})
		})
	})
	if pay.Clicked() && !*paid {
		*paid = true
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("PaymentRequest", "paid")
	}
	return root, pay
}
