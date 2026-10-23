package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MarketStatusOptions are the market labels.
type MarketStatusOptions struct{ Open, Closed, Next, DateFormat string }

// MarketStatus shows open or closed with time to next.
func MarketStatus(c *ui.Context, open bool, next time.Time, o MarketStatusOptions) ui.Element {
	if next.IsZero() {
		panic("finance: MarketStatus needs a next transition")
	}
	def(&o.Open, "Market open")
	def(&o.Closed, "Market closed")
	def(&o.Next, "Next %s")
	def(&o.DateFormat, "Jan 02 15:04")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	left := max(0, next.Sub(c.Now()))
	if left > 0 {
		c.After(min(time.Second, left))
	}
	root.Children(func() {
		text, ink := o.Closed, t.Danger
		if open {
			text, ink = o.Open, t.Success
		}
		ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
		ui.Text(c, text).Bold().Grow(1)
		kit.Muted(c, fmt.Sprintf(o.Next, next.Format(o.DateFormat))).FontSize(t.Rem(0.75))
		if left > 0 {
			kit.Numeric(c, left.Round(time.Second).String()).FontSize(t.Rem(0.75))
		}
	})
	return root
}
