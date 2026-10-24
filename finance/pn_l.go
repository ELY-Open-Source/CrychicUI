package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PnLOptions word the card; a positive Basis adds percent.
type PnLOptions struct {
	Label, Period string
	Basis         float64
}

// PnL shows profit and loss in sign color; rolls.
func PnL(c *ui.Context, value float64, currency string, o PnLOptions) ui.Element {
	if !finite(value) || currency == "" || o.Basis < 0 || !finite(o.Basis) {
		panic("finance: PnL needs a finite value, a currency and a basis of zero or more")
	}
	change := 0.0
	if o.Basis > 0 {
		change = 100 * (value / o.Basis)
		if !finite(change) {
			panic("finance: PnL change overflows")
		}
	}
	def(&o.Label, "Profit and loss")
	t := c.Theme()
	ink := changeInk(t, value)
	glyph := "trending-up"
	if value < 0 {
		glyph = "trending-down"
	}
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, glyph, ink)
			kit.Muted(c, o.Label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
			if o.Period != "" {
				kit.Chip(c, o.Period, t.TextMuted)
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			amount := ui.Row(c).Grow(1).MinWidth(0)
			amount.TextColor(amount.AnimateColor("ink", ink, ui.Snappy)).Children(func() {
				sign := ""
				if value > 0 {
					sign = "+"
				}
				ui.Text(c, sign).Bold().FontSize(t.Rem(1.8)).Shrink(0)
				money(c, value, currency).Bold().FontSize(t.Rem(1.8)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			})
			if o.Basis > 0 {
				changePill(c, change)
			}
		})
	})
	return root
}
