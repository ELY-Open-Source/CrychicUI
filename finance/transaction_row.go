package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TransactionRow shows one transaction; it fades in.
func TransactionRow(c *ui.Context, row Transaction, currency string) ui.Element {
	return kit.Arrive(txRow(c, row, currency), ui.Pose{Y: c.Theme().Space(1), Fade: 1}, 0)
}

// txRow builds a transaction row without its arrival.
func txRow(c *ui.Context, row Transaction, currency string) ui.Element {
	if row.Merchant == "" || row.Category == "" || !finite(row.Amount) || currency == "" {
		panic(fmt.Sprintf("finance: TransactionRow %q needs a category, a finite amount and a currency", row.Merchant))
	}
	t := c.Theme()
	glyph, ink := row.Icon, palette(t, len(row.Category)+int(row.Category[0]))
	switch {
	case row.Amount > 0:
		ink = t.Success
		if glyph == "" {
			glyph = "arrow-down-left"
		}
	case glyph == "":
		glyph = "shopping-bag"
	}
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0)
	return root.Children(func() {
		tile(c, glyph, ink)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, row.Merchant).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
			meta := row.Category
			if row.When != "" {
				meta += " · " + row.When
			}
			kit.Muted(c, meta).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
		})
		amount := money(c, row.Amount, currency).Bold().Shrink(0)
		if row.Amount > 0 {
			amount.TextColor(t.Success)
		}
	})
}
