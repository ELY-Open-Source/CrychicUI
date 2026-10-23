package finance

import "github.com/egoist/mygo/ui"

// Transaction is one merchant transaction.
type Transaction struct {
	Merchant, Category string
	Amount             float64
}

// Transactions shows merchant rows with signed amounts.
func Transactions(c *ui.Context, rows []Transaction, currency string) ui.Element {
	if len(rows) == 0 || currency == "" {
		panic("finance: Transactions needs rows and currency")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, row := range rows {
			c.Key(i)
			TransactionRow(c, row, currency)
		}
	})
	return root
}
