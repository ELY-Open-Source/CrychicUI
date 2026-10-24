package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SLAOptions word the card; Budget takes the error budget left.
type SLAOptions struct{ Title, Actual, Target, Met, Missed, Budget string }

// SLA shows met percent against target and error budget.
func SLA(c *ui.Context, actual, target float64, o SLAOptions) ui.Element {
	if actual < 0 || actual > 100 || target <= 0 || target > 100 || !finite(actual) || !finite(target) {
		panic("dashboard: SLA needs actual in [0,100] and target in (0,100]")
	}
	def(&o.Title, "SLA")
	def(&o.Actual, "%.2f%%")
	def(&o.Target, "Target %.2f%%")
	def(&o.Met, "Met")
	def(&o.Missed, "Missed")
	def(&o.Budget, "%s of error budget left")
	met := actual >= target
	left := 0.0
	switch {
	case target == 100 && actual == 100:
		left = 1
	case target < 100:
		left = max(0, 1-(100-actual)/(100-target))
	}
	t := c.Theme()
	ink := t.Success
	switch {
	case !met:
		ink = t.Danger
	case left < 0.25:
		ink = t.Warning
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "shield-check", o.Title, func() {
			if met {
				kit.Chip(c, o.Met, t.Success)
			} else {
				kit.Chip(c, o.Missed, t.Danger)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
				kit.Numeric(c, fmt.Sprintf(o.Actual, actual)).Bold().FontSize(t.Rem(1.8)).TextColor(ink).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				kit.Muted(c, fmt.Sprintf(o.Target, target)).FontSize(t.Rem(0.8)).Padding(0, 0, t.Space(1), 0).Grow(1).MinWidth(0).
					TextAlign(ui.End).SingleLine()
			})
			progress(c, left, ink)
			kit.Muted(c, fmt.Sprintf(o.Budget, fmt.Sprintf("%.0f%%", 100*left))).FontSize(t.Rem(0.75)).Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root
}
