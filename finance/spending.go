package finance

import "github.com/egoist/mygo/ui"

// Spend is one category's spending.
type Spend struct {
	Category string
	Amount   float64
}

// Spending shows category bars of spend.
func Spending(c *ui.Context, rows []Spend, currency string) ui.Element {
	if len(rows) == 0 || currency == "" {
		panic("finance: Spending needs rows and currency")
	}
	t := c.Theme()
	total := 0.0
	for _, row := range rows {
		if row.Category == "" || row.Amount < 0 || !finite(row.Amount) {
			panic("finance: Spending has an invalid row")
		}
		total += row.Amount
	}
	if !finite(total) {
		panic("finance: Spending total overflow")
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, row := range rows {
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, row.Category).Grow(1)
					money(c, row.Amount, currency)
				})
				fraction := 0.0
				if total > 0 {
					fraction = row.Amount / total
				}
				progress(c, fraction, t.Accent)
			})
		}
	})
	return root
}
