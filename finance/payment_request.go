package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PaymentRequestOptions are the payment labels.
type PaymentRequestOptions struct{ Pay string }

// PaymentRequest shows amount, payer, and pay; cue Success.
func PaymentRequest(c *ui.Context, amount float64, payer, currency string, paid *bool, o PaymentRequestOptions) (root, pay ui.Element) {
	if amount <= 0 || !finite(amount) || payer == "" || currency == "" || paid == nil {
		panic("finance: PaymentRequest has invalid data")
	}
	def(&o.Pay, "Pay")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Column(c).Grow(1).Children(func() {
			money(c, amount, currency).Bold()
			kit.Muted(c, payer).FontSize(t.Rem(0.8))
		})
		pay = kit.Press(ui.ButtonBase(c).Label(o.Pay).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
			Background(t.Accent).TextColor(t.AccentText).Disabled(*paid))
		pay.Children(func() { ui.Text(c, o.Pay).SingleLine() })
	})
	if pay.Clicked() && !*paid {
		*paid = true
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("PaymentRequest", "paid")
	}
	return root, pay
}
