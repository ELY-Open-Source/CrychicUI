package finance

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ExchangeRateOptions add a trend line; Per takes code, rate, code.
type ExchangeRateOptions struct {
	Points []float64
	Per    string
}

// ExchangeRate shows a pair "BASE/QUOTE" with rate, change and trend.
func ExchangeRate(c *ui.Context, pair string, rate, change float64, o ExchangeRateOptions) ui.Element {
	base, quote, ok := strings.Cut(pair, "/")
	if !ok || base == "" || quote == "" || strings.Contains(quote, "/") || rate <= 0 || !finite(rate) || !finite(change) || !finite(1/rate) {
		panic(fmt.Sprintf("finance: ExchangeRate %q needs a BASE/QUOTE pair, a positive rate and a finite change", pair))
	}
	if len(o.Points) == 1 {
		panic(fmt.Sprintf("finance: ExchangeRate %q needs two or more points", pair))
	}
	def(&o.Per, "1 %s = %s %s")
	t := c.Theme()
	coin := func(code string, i int) ui.Element {
		return ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).Background(palette(t, i).Mix(t.Background, 0.82)).
			Border(2, t.Background).Children(func() {
			ui.Text(c, code).FontSize(t.Rem(0.6)).Bold().TextColor(palette(t, i)).Shrink(0)
		})
	}
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Shrink(0).Children(func() {
				coin(base, 0)
				coin(quote, 1).MarginStart(-t.Space(1.5))
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, pair).Bold().FillWidth().SingleLine()
				kit.Muted(c, fmt.Sprintf(o.Per, quote, fmt.Sprintf("%.4f", 1/rate), base)).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
			})
			ui.Column(c).Shrink(0).AlignItems(ui.End).Gap(t.Space(1)).Children(func() {
				kit.Numeric(c, fmt.Sprintf("%.4f", rate)).Bold().FontSize(t.Rem(1.2)).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				changePill(c, change)
			})
		})
		if len(o.Points) > 1 {
			sparkline(c, o.Points, changeInk(t, change).Mix(t.Accent, unit(change == 0))).Height(t.Space(12)).FillWidth()
		}
	})
	return root
}
