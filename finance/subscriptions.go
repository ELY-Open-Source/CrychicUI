package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Subscription is one recurring charge; Plan and Icon are optional.
type Subscription struct {
	Name, Next string
	Amount     float64
	Plan, Icon string
}

// SubscriptionsOptions word the card; Renews takes the next date.
type SubscriptionsOptions struct{ Title, Renews string }

// Subscriptions lists recurring charges with their next renewal.
func Subscriptions(c *ui.Context, rows []Subscription, currency string, o SubscriptionsOptions) ui.Element {
	if len(rows) == 0 || currency == "" {
		panic("finance: Subscriptions needs rows and a currency")
	}
	for _, row := range rows {
		if row.Name == "" || row.Next == "" || row.Amount < 0 || !finite(row.Amount) {
			panic(fmt.Sprintf("finance: Subscriptions row %q needs a next date and an amount of zero or more", row.Name))
		}
	}
	def(&o.Title, "Subscriptions")
	def(&o.Renews, "Renews %s")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "repeat", o.Title, func() { kit.Chip(c, fmt.Sprint(len(rows)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3))
		list.Children(func() {
			for i, row := range rows {
				c.Key(i)
				glyph := row.Icon
				if glyph == "" {
					glyph = "repeat"
				}
				line := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(rows)-1)))
				kit.Stagger(c, list, line, i).Children(func() {
					tile(c, glyph, palette(t, i))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, row.Name).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
						meta := fmt.Sprintf(o.Renews, row.Next)
						if row.Plan != "" {
							meta = row.Plan + " · " + meta
						}
						kit.Muted(c, meta).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
					})
					money(c, row.Amount, currency).Bold().FontSize(t.Rem(0.85)).Shrink(0)
				})
			}
		})
	})
	return root
}
