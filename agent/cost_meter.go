package agent

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// CostMeterOptions word the label, budget and share.
type CostMeterOptions struct {
	Label, Budget, Used string
}

// CostMeter shows spend against a budget; digits roll.
func CostMeter(c *ui.Context, amount, budget float64, currency string, o CostMeterOptions) ui.Element {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount < 0 || math.IsNaN(budget) || math.IsInf(budget, 0) ||
		budget <= 0 || currency == "" {
		panic(fmt.Sprintf("agent: CostMeter needs a finite amount, a budget and a currency, got %v of %v", amount, budget))
	}
	def(&o.Label, "Spend this month")
	def(&o.Budget, "Budget")
	def(&o.Used, "%.0f%% used")
	t := c.Theme()
	share := amount / budget
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Role(ui.RoleGroup).Label(o.Label)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			tile(c, "coins", t.Accent, t.Space(7))
			kit.Muted(c, o.Label).FontSize(t.Rem(0.85)).Grow(1).SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Used, share*100)).FontSize(t.Rem(0.75)).TextColor(meterInk(c, share))
		})
		typography.Money(c, amount, currency).Bold().FontSize(t.Rem(1.5))
		bar(c, "spend", share, meterInk(c, share))
		ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
			kit.Muted(c, o.Budget).FontSize(t.Rem(0.75)).Grow(1)
			typography.Money(c, budget, currency).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		})
	})
}
