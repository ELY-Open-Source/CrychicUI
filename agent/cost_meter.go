package agent

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// CostMeterOptions tune a CostMeter.
type CostMeterOptions struct {
	Label string
}

// CostMeter shows the spend so far; digits roll.
func CostMeter(c *ui.Context, amount float64, currency string, o CostMeterOptions) ui.Element {
	if math.IsNaN(amount) || math.IsInf(amount, 0) || amount < 0 || currency == "" {
		panic(fmt.Sprintf("agent: CostMeter needs a finite amount and a currency, got %v", amount))
	}
	def(&o.Label, "Spend")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(o.Label)
	return root.Children(func() {
		kit.Muted(c, o.Label).FontSize(t.Rem(0.85))
		typography.Money(c, amount, currency).Bold()
	})
}
