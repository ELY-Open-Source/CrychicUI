package dashboard

import (
	"fmt"
	"math/big"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// MRROptions word the card; Net takes the signed net new.
type MRROptions struct{ Title, New, Expansion, Churn, Net string }

// MRR shows recurring revenue split into new, expansion and churn.
func MRR(c *ui.Context, current, new, expansion, churn float64, currency string, o MRROptions) ui.Element {
	def(&o.Title, "MRR")
	if current < 0 || new < 0 || expansion < 0 || churn < 0 || !finite(current) || !finite(new) || !finite(expansion) || !finite(churn) || currency == "" {
		panic(fmt.Sprintf("dashboard: MRR %q needs finite amounts of zero or more and a currency", o.Title))
	}
	// Exact sum, rounded once, so cancellation keeps small nets.
	sum := big.NewFloat(new).SetPrec(2200)
	sum.Add(sum, big.NewFloat(expansion)).Sub(sum, big.NewFloat(churn))
	net, _ := sum.Float64()
	if !finite(net) {
		panic(fmt.Sprintf("dashboard: MRR %q net overflows", o.Title))
	}
	def(&o.New, "New")
	def(&o.Expansion, "Expansion")
	def(&o.Churn, "Churn")
	def(&o.Net, "Net new %s")
	peak := max(new, expansion, churn)
	if peak == 0 {
		peak = 1
	}
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "repeat", o.Title, func() {
			sign := ""
			if net > 0 {
				sign = "+"
			}
			kit.Chip(c, fmt.Sprintf(o.Net, sign+moneyText(c, net, currency)), changeOf(t, net))
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3))
		list.Children(func() {
			typography.Money(c, current, currency).Bold().FontSize(t.Rem(1.8)).ContentTransition(ui.ContentTransition{Numeric: true})
			for i, row := range []struct {
				name, glyph string
				value       float64
				ink         ui.Color
			}{{o.New, "plus", new, t.Success}, {o.Expansion, "arrow-up-right", expansion, t.Accent}, {o.Churn, "minus", churn, t.Danger}} {
				c.Key(i)
				line := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				kit.Stagger(c, list, line, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999).Background(row.ink.Alpha(0.12)).
							TextColor(row.ink).Children(func() { ui.Icon(c, icon.Must(row.glyph)).Size(12, 12) })
						ui.Text(c, row.name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						typography.Money(c, row.value, currency).Bold().FontSize(t.Rem(0.85)).Shrink(0)
					})
					progress(c, row.value/peak, row.ink)
				})
			}
		})
	})
	return root
}

// changeOf is green above zero, red below, else muted.
func changeOf(t *ui.Theme, v float64) ui.Color {
	switch {
	case v > 0:
		return t.Success
	case v < 0:
		return t.Danger
	}
	return t.TextMuted
}
