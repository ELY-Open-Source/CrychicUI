package finance

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LoanCalcOptions are the loan labels.
type LoanCalcOptions struct{ Payment string }

// LoanCalc calculates a monthly payment; digits roll.
func LoanCalc(c *ui.Context, principal, annualRate float64, years int, currency string, o LoanCalcOptions) ui.Element {
	if principal <= 0 || annualRate < 0 || years <= 0 || !finite(principal) || !finite(annualRate) || currency == "" {
		panic("finance: LoanCalc has invalid data")
	}
	def(&o.Payment, "Monthly payment")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		kit.Muted(c, o.Payment).FontSize(t.Rem(0.8))
		money(c, loanPayment(principal, annualRate, years), currency).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
	})
	return root
}

func loanPayment(principal, annualRate float64, years int) float64 {
	n := float64(years) * 12
	monthly := annualRate / 12
	if monthly == 0 {
		return principal / n
	}
	denominator := -math.Expm1(-n * math.Log1p(monthly))
	payment := principal * monthly / denominator
	if payment <= 0 || !finite(payment) {
		panic("finance: loan payment overflow")
	}
	return payment
}
