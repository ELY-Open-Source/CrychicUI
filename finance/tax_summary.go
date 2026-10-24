package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TaxSummaryOptions word the card; Effective takes the rate.
type TaxSummaryOptions struct{ Title, Income, Deductions, Taxable, Tax, Effective, Keep string }

// TaxSummary shows income, deductions, taxable income and tax due.
func TaxSummary(c *ui.Context, income, deductions, rate float64, currency string, o TaxSummaryOptions) ui.Element {
	if income < 0 || deductions < 0 || rate < 0 || rate > 1 || !finite(income) || !finite(deductions) || !finite(rate) || currency == "" {
		panic("finance: TaxSummary has invalid data")
	}
	def(&o.Title, "Tax summary")
	def(&o.Income, "Income")
	def(&o.Deductions, "Deductions")
	def(&o.Taxable, "Taxable income")
	def(&o.Tax, "Tax due")
	def(&o.Effective, "%s effective")
	def(&o.Keep, "Take-home")
	taxable := max(0, income-deductions)
	tax := taxable * rate
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "landmark", o.Title, func() {
			if income > 0 {
				kit.Chip(c, fmt.Sprintf(o.Effective, fmt.Sprintf("%.1f%%", 100*(tax/income))), t.Accent)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		list.Children(func() {
			for i, line := range []string{o.Income, o.Deductions, o.Taxable} {
				c.Key(i)
				value, glyph := []float64{income, deductions, taxable}[i], []string{"arrow-down-left", "minus", "equal"}[i]
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
				kit.Stagger(c, list, row, i).Children(func() {
					tile(c, glyph, palette(t, i)).Size(t.Space(7), t.Space(7))
					ui.Text(c, line).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					money(c, value, currency).FontSize(t.Rem(0.85)).Shrink(0)
				})
			}
			ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(3), 0, 0, 0).Margin(t.Space(1), 0, 0, 0).
				BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
				ui.Text(c, o.Tax).Bold().Grow(1).MinWidth(0).SingleLine()
				money(c, tax, currency).Bold().FontSize(t.Rem(1.3)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			})
			if income > 0 {
				split(c, []float64{max(0, income-tax), tax}, []string{o.Keep, o.Tax}, currency)
			}
		})
	})
	return root
}
