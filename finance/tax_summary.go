package finance

import "github.com/egoist/mygo/ui"

// TaxSummaryOptions are the tax labels.
type TaxSummaryOptions struct{ Income, Deductions, Tax string }

// TaxSummary shows income, deductions, and tax due.
func TaxSummary(c *ui.Context, income, deductions, rate float64, currency string, o TaxSummaryOptions) ui.Element {
	if income < 0 || deductions < 0 || rate < 0 || rate > 1 || !finite(income) || !finite(deductions) || !finite(rate) || currency == "" {
		panic("finance: TaxSummary has invalid data")
	}
	def(&o.Income, "Income")
	def(&o.Deductions, "Deductions")
	def(&o.Tax, "Tax due")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, row := range []struct {
			name  string
			value float64
		}{{o.Income, income}, {o.Deductions, deductions}, {o.Tax, max(0, (income-deductions)*rate)}} {
			ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, row.name).Grow(1)
				money(c, row.value, currency)
			})
		}
	})
	return root
}
