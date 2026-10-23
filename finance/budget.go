package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BudgetItem is one category budget.
type BudgetItem struct {
	Category    string
	Used, Limit float64
}

// BudgetOptions are the budget labels.
type BudgetOptions struct{ Left, Over string }

// Budget shows category budgets with used bars.
func Budget(c *ui.Context, rows []BudgetItem, currency string, o BudgetOptions) ui.Element {
	if len(rows) == 0 || currency == "" {
		panic("finance: Budget needs rows and currency")
	}
	def(&o.Left, "%s left")
	def(&o.Over, "%s over")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, row := range rows {
			if row.Category == "" || row.Used < 0 || row.Limit <= 0 || !finite(row.Used) || !finite(row.Limit) {
				panic("finance: Budget has an invalid row")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, row.Category).Grow(1)
					text, ink := o.Left, t.Success
					if row.Used > row.Limit {
						text, ink = o.Over, t.Danger
					}
					kit.Muted(c, fmt.Sprintf(text, moneyText(c, abs(row.Limit-row.Used), currency))).TextColor(ink)
				})
				progress(c, min(1, row.Used/row.Limit), t.Accent)
			})
		}
	})
	return root
}
