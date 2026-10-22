package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FiscalYearOptions tune a FiscalYear.
type FiscalYearOptions struct {
	// Quarters name the four quarters.
	Quarters [4]string
}

// FiscalYear fills four quarter bars by the time elapsed.
func FiscalYear(c *ui.Context, label string, start time.Time, o FiscalYearOptions) ui.Element {
	if label == "" || start.IsZero() || !start.Equal(monthOf(start)) {
		panic("datetime: FiscalYear needs a label and a start at a month's first midnight")
	}
	if o.Quarters == ([4]string{}) {
		o.Quarters = [4]string{"Q1", "Q2", "Q3", "Q4"}
	}
	t := c.Theme()
	now := c.Now()
	if now.Before(addMonths(start, 12)) {
		c.After(time.Hour)
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return kit.Arrive(root, ui.Pose{Y: 4, Fade: 1}, 0).Children(func() {
		ui.Text(c, label).Bold()
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for q, name := range o.Quarters {
				from, to := addMonths(start, 3*q), addMonths(start, 3*q+3)
				f := min(max(float64(now.Sub(from))/float64(to.Sub(from)), 0), 1)
				c.Key(q)
				ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Children(func() {
					feedback.ProgressBar(c, name, f)
					kit.Muted(c, name).FontSize(t.Rem(0.75))
				})
			}
		})
	})
}
