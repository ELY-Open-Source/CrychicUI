package finance

import "github.com/egoist/mygo/ui"

// StatementOptions are the statement labels.
type StatementOptions struct{ Total string }

// Statement shows monthly statement rows and totals.
func Statement(c *ui.Context, month string, rows []Transaction, currency string, o StatementOptions) ui.Element {
	if month == "" || len(rows) == 0 || currency == "" {
		panic("finance: Statement has invalid data")
	}
	def(&o.Total, "Total")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	total := 0.0
	root.Children(func() {
		ui.Text(c, month).Bold()
		for i, row := range rows {
			c.Key(i)
			TransactionRow(c, row, currency)
			total += row.Amount
			if !finite(total) {
				panic("finance: Statement total overflow")
			}
		}
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1)
			money(c, total, currency).Bold()
		})
	})
	return root
}
