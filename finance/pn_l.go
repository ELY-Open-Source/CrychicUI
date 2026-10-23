package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PnLOptions are the profit and loss labels.
type PnLOptions struct{ Label string }

// PnL shows profit and loss with sign color; rolls.
func PnL(c *ui.Context, value float64, currency string, o PnLOptions) ui.Element {
	if !finite(value) || currency == "" {
		panic("finance: PnL has invalid data")
	}
	def(&o.Label, "Profit and loss")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		kit.Muted(c, o.Label).FontSize(t.Rem(0.75))
		ink := t.Success
		if value < 0 {
			ink = t.Danger
		}
		money(c, value, currency).Bold().TextColor(ink).ContentTransition(ui.ContentTransition{Numeric: true})
	})
	return root
}
