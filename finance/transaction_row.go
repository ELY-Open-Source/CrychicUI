package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TransactionRow shows one transaction; it fades in.
func TransactionRow(c *ui.Context, row Transaction, currency string) ui.Element {
	if row.Merchant == "" || row.Category == "" || !finite(row.Amount) || currency == "" {
		panic("finance: TransactionRow has invalid data")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Radius(t.Radius)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	root.Children(func() {
		ui.Icon(c, icon.Must("shopping-bag")).Size(16, 16).TextColor(t.TextMuted)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, row.Merchant).SingleLine()
			kit.Muted(c, row.Category).FontSize(t.Rem(0.75))
		})
		ink := t.Success
		if row.Amount < 0 {
			ink = t.Danger
		}
		money(c, row.Amount, currency).TextColor(ink)
	})
	return root
}
