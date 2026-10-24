package finance

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MortgageOptions word the card; Year takes the year number.
type MortgageOptions struct{ Title, Payment, Year, Principal, Interest string }

// Mortgage charts principal against interest each year.
func Mortgage(c *ui.Context, principal, annualRate float64, years int, currency string, o MortgageOptions) ui.Element {
	if principal <= 0 || annualRate < 0 || years <= 0 || years > 50 || !finite(principal) || !finite(annualRate) || currency == "" {
		panic("finance: Mortgage needs a positive principal, a rate and 1–50 years")
	}
	def(&o.Title, "Mortgage")
	def(&o.Payment, "Monthly payment")
	def(&o.Year, "Year %d")
	def(&o.Principal, "Principal")
	def(&o.Interest, "Interest")
	payment := loanPayment(principal, annualRate, years)
	plan := amortize(payment, annualRate, years)
	interest, peak := 0.0, 0.0
	for _, s := range plan {
		interest += s.interest
		peak = max(peak, s.principal+s.interest)
	}
	if !finite(interest) || !finite(peak) {
		panic("finance: Mortgage yearly totals overflow")
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "house", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			ui.Column(c).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Payment).FontSize(t.Rem(0.8))
				money(c, payment, currency).Bold().FontSize(t.Rem(1.6)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			})
			chart := ui.Box(c).FillWidth().Height(t.Space(30))
			grow := chart.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				slot := r.W / float32(years)
				w := max(1, slot*0.7)
				for y, s := range plan {
					x := r.X + slot*float32(y) + (slot-w)/2
					hp := r.H * float32(s.principal/peak) * grow
					hi := r.H * float32(s.interest/peak) * grow
					p.Fill(ui.Rect{X: x, Y: r.Y + r.H - hp, W: w, H: hp}, palette(t, 0), 0)
					p.Fill(ui.Rect{X: x, Y: r.Y + r.H - hp - hi, W: w, H: hi}, palette(t, 1), 0)
				}
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, t.Border, 0)
			}))
			ui.Row(c).FillWidth().Children(func() {
				kit.Muted(c, fmt.Sprintf(o.Year, 1)).FontSize(t.Rem(0.7)).Grow(1)
				if years > 1 {
					kit.Muted(c, fmt.Sprintf(o.Year, years)).FontSize(t.Rem(0.7)).Shrink(0)
				}
			})
			split(c, []float64{principal, interest}, []string{o.Principal, o.Interest}, currency)
		})
	})
	return root
}

type year struct{ principal, interest float64 }

// amortize splits yearly payments; month k repays payment·(1+m)^-(n-k+1).
func amortize(payment, annualRate float64, years int) []year {
	plan := make([]year, years)
	n, m := years*12, annualRate/12
	for k := 1; k <= n; k++ {
		paid, owed := payment, 0.0
		if m > 0 {
			j := -float64(n-k+1) * math.Log1p(m)
			paid, owed = payment*math.Exp(j), payment*-math.Expm1(j)
		}
		plan[(k-1)/12].principal += paid
		plan[(k-1)/12].interest += owed
	}
	return plan
}
