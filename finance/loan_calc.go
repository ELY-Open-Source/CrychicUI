package finance

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LoanCalcOptions word the card; Year and Years take the term.
type LoanCalcOptions struct{ Title, Payment, Amount, Rate, Term, Year, Years, Principal, Interest string }

// LoanCalc shows the monthly payment and cost split; rolls.
func LoanCalc(c *ui.Context, principal, annualRate float64, years int, currency string, o LoanCalcOptions) ui.Element {
	if principal <= 0 || annualRate < 0 || years <= 0 || years > 50 || !finite(principal) || !finite(annualRate) || currency == "" {
		panic("finance: LoanCalc needs a positive principal, a rate, 1–50 years and a currency")
	}
	def(&o.Title, "Loan")
	def(&o.Payment, "Monthly payment")
	def(&o.Amount, "Amount")
	def(&o.Rate, "Rate")
	def(&o.Term, "Term")
	def(&o.Year, "%d year")
	def(&o.Years, "%d years")
	def(&o.Principal, "Principal")
	def(&o.Interest, "Interest")
	payment := loanPayment(principal, annualRate, years)
	interest := 0.0
	for _, y := range amortize(payment, annualRate, years) {
		interest += y.interest
	}
	if !finite(interest) {
		panic("finance: LoanCalc interest overflows")
	}
	if !finite(annualRate * 100) {
		panic("finance: LoanCalc rate overflows")
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "calculator", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			ui.Column(c).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Payment).FontSize(t.Rem(0.8))
				money(c, payment, currency).Bold().FontSize(t.Rem(1.8)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				fact(c, o.Amount, moneyText(c, principal, currency))
				fact(c, o.Rate, fmt.Sprintf("%.2f%%", annualRate*100))
				term := o.Years
				if years == 1 {
					term = o.Year
				}
				fact(c, o.Term, fmt.Sprintf(term, years))
			})
			split(c, []float64{principal, interest}, []string{o.Principal, o.Interest}, currency)
		})
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
	payment := principal * (monthly / denominator)
	if payment <= 0 || !finite(payment) {
		panic("finance: loan payment overflow")
	}
	return payment
}
