package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Bill is one due bill.
type Bill struct {
	Name, Due string
	Amount    float64
}

// BillsOptions are the bill labels.
type BillsOptions struct{ Pay string }

// Bills shows due bills with pay; cue Success.
func Bills(c *ui.Context, bills []Bill, paid *int, currency string, o BillsOptions) ui.Element {
	if len(bills) == 0 || paid == nil || *paid < -1 || *paid >= len(bills) || currency == "" {
		panic("finance: Bills needs bills and payment state")
	}
	def(&o.Pay, "Pay")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, bill := range bills {
			if bill.Name == "" || bill.Due == "" || bill.Amount <= 0 || !finite(bill.Amount) {
				panic("finance: Bills has an invalid bill")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, bill.Name).SingleLine()
					kit.Muted(c, bill.Due).FontSize(t.Rem(0.75))
				})
				money(c, bill.Amount, currency)
				b := kit.Press(ui.ButtonBase(c).Label(o.Pay).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
				b.Children(func() { ui.Text(c, o.Pay).SingleLine() })
				if b.Clicked() {
					*paid = i
					root.SetChanged()
					sound.Play(sound.Success)
					kit.Log("Bills", "paid", "index", i)
				}
			})
		}
	})
	return root
}
