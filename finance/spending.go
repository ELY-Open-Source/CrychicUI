package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Spend is one category's spending.
type Spend struct {
	Category string
	Amount   float64
}

// SpendingOptions word the card.
type SpendingOptions struct{ Title string }

// Spending shows a split bar and each category's share.
func Spending(c *ui.Context, rows []Spend, currency string, o SpendingOptions) ui.Element {
	if len(rows) == 0 || currency == "" {
		panic("finance: Spending needs rows and a currency")
	}
	total := 0.0
	for _, row := range rows {
		if row.Category == "" || row.Amount < 0 || !finite(row.Amount) {
			panic(fmt.Sprintf("finance: Spending row %q needs a finite amount of zero or more", row.Category))
		}
		total += row.Amount
	}
	if !finite(total) {
		panic("finance: Spending total overflows")
	}
	def(&o.Title, "Spending")
	share := func(v float64) float64 {
		if total == 0 {
			return 0
		}
		return v / total
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "chart-pie", o.Title, func() { money(c, total, currency).Bold().Shrink(0) })
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			bar := ui.Box(c).FillWidth().Height(t.Space(3)).Radius(999).Background(t.Surface)
			grow := bar.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				x, gap := r.X, t.Space(0.5)
				for i, row := range rows {
					w := r.W * float32(share(row.Amount)) * grow
					if w > gap {
						p.Fill(ui.Rect{X: x, Y: r.Y, W: w - gap, H: r.H}, palette(t, i), 0)
					}
					x += w
				}
			}))
			list := ui.Column(c).FillWidth().Gap(t.Space(2))
			list.Children(func() {
				for i, row := range rows {
					c.Key(i)
					line := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center)
					kit.Stagger(c, list, line, i).Children(func() {
						ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(palette(t, i))
						ui.Text(c, row.Category).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*share(row.Amount))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Shrink(0)
						money(c, row.Amount, currency).FontSize(t.Rem(0.85)).Bold().Shrink(0)
					})
				}
			})
		})
	})
	return root
}
