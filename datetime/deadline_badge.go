package datetime

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// DeadlineBadgeOptions tune a DeadlineBadge.
type DeadlineBadgeOptions struct {
	// Overdue and Today are copy; Later formats days left.
	Overdue, Today, Later string
}

// DeadlineBadge says how far off due is; overdue pulses.
func DeadlineBadge(c *ui.Context, due time.Time, o DeadlineBadgeOptions) ui.Element {
	if due.IsZero() {
		panic("datetime: DeadlineBadge needs a due date")
	}
	def(&o.Overdue, "Overdue")
	def(&o.Today, "Due today")
	def(&o.Later, "In %d d")
	t := c.Theme()
	now := c.Now()
	c.After(tomorrow(now).Sub(now))
	n := days(now, due)
	text, bg := fmt.Sprintf(o.Later, n), t.Accent
	switch {
	case n < 0:
		text, bg = o.Overdue, t.Danger
	case n == 0:
		text, bg = o.Today, t.Warning
	}
	root := ui.Row(c).AlignSelf(ui.Start).Padding(t.Space(0.5), t.Space(2)).Radius(999).TextColor(t.AccentText)
	root.Background(root.AnimateColor("bg", bg, ui.Smooth))
	if n < 0 {
		k := 1 + 0.06*root.Loop("pulse", 1200*time.Millisecond, ui.Bounce(ui.EaseInOut))
		root.Scale(k, k)
	}
	return root.Children(func() {
		ui.Text(c, text).FontSize(t.Rem(0.85)).Bold().ContentTransition(ui.ContentTransition{})
	})
}
