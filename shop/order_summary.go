package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SummaryLine is one order amount; negative ones are discounts.
type SummaryLine struct {
	Label  string
	Amount float64
}

// OrderSummaryOptions word the summary.
type OrderSummaryOptions struct{ Title, Total string }

// OrderSummary shows lines, green discounts and a rolling total.
func OrderSummary(c *ui.Context, lines []SummaryLine, currency string, o OrderSummaryOptions) ui.Element {
	if len(lines) == 0 || currency == "" {
		panic("shop: OrderSummary needs lines and a currency")
	}
	total := 0.0
	for _, line := range lines {
		if line.Label == "" || !finite(line.Amount) {
			panic(fmt.Sprintf("shop: OrderSummary line %q needs a label and a finite amount", line.Label))
		}
		total += line.Amount
	}
	def(&o.Title, "Order summary")
	def(&o.Total, "Total")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(o.Title)
	root.Children(func() {
		kit.CardHead(c, "receipt", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			for i, line := range lines {
				c.Key(i)
				row := ui.Row(c).FillWidth().AlignItems(ui.Center)
				kit.Stagger(c, root, row, i).Children(func() {
					ui.Text(c, line.Label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					amount := signedMoney(c, line.Amount, currency).FontSize(t.Rem(0.85)).ContentTransition(ui.ContentTransition{Numeric: true})
					if line.Amount < 0 {
						amount.TextColor(t.Success)
					}
				})
			}
		})
		ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(3)).Background(t.Surface).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1)
			signedMoney(c, total, currency).Bold().FontSize(t.Rem(1.2)).ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root
}
