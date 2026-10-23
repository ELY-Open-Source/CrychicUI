package finance

import "github.com/egoist/mygo/ui"

// FlowMonth is one month's incoming and outgoing cash.
type FlowMonth struct {
	Label   string
	In, Out float64
}

// CashFlow shows in and out bars by month.
func CashFlow(c *ui.Context, months []FlowMonth, currency string) ui.Element {
	if len(months) == 0 || currency == "" {
		panic("finance: CashFlow needs months and currency")
	}
	t := c.Theme()
	maxValue := 0.0
	for _, month := range months {
		if month.Label == "" || month.In < 0 || month.Out < 0 || !finite(month.In) || !finite(month.Out) {
			panic("finance: CashFlow has an invalid month")
		}
		maxValue = max(maxValue, month.In, month.Out)
	}
	if maxValue == 0 {
		maxValue = 1
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, month := range months {
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Text(c, month.Label).FontSize(t.Rem(0.75))
				ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
					ui.Column(c).Grow(1).Basis(0).Children(func() {
						progress(c, month.In/maxValue, t.Success)
						money(c, month.In, currency).FontSize(t.Rem(0.75))
					})
					ui.Column(c).Grow(1).Basis(0).Children(func() {
						progress(c, month.Out/maxValue, t.Danger)
						money(c, month.Out, currency).FontSize(t.Rem(0.75))
					})
				})
			})
		}
	})
	return root
}
