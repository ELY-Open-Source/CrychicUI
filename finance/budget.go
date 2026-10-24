package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BudgetItem is one category budget.
type BudgetItem struct {
	Category    string
	Used, Limit float64
}

// BudgetOptions word the card; Of takes used and limit.
type BudgetOptions struct{ Title, Left, Over, Of string }

// Budget shows category budgets with bars warming as they fill.
func Budget(c *ui.Context, rows []BudgetItem, currency string, o BudgetOptions) ui.Element {
	if len(rows) == 0 || currency == "" {
		panic("finance: Budget needs rows and a currency")
	}
	for _, row := range rows {
		if row.Category == "" || row.Used < 0 || row.Limit <= 0 || !finite(row.Used) || !finite(row.Limit) {
			panic(fmt.Sprintf("finance: Budget row %q needs used of zero or more and a positive limit", row.Category))
		}
	}
	def(&o.Title, "Budgets")
	def(&o.Left, "%s left")
	def(&o.Over, "%s over")
	def(&o.Of, "%s of %s")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "piggy-bank", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(3.5)).Padding(t.Space(3))
		list.Children(func() {
			for i, row := range rows {
				c.Key(i)
				used := row.Used / row.Limit
				ink := t.Accent
				switch {
				case used > 1:
					ink = t.Danger
				case used >= 0.8:
					ink = t.Warning
				}
				line := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				kit.Stagger(c, list, line, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
							ui.Text(c, row.Category).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
							kit.Muted(c, fmt.Sprintf(o.Of, moneyText(c, row.Used, currency), moneyText(c, row.Limit, currency))).
								FontSize(t.Rem(0.75)).FillWidth().SingleLine()
						})
						text, chip := o.Left, t.Success
						if used > 1 {
							text, chip = o.Over, t.Danger
						}
						kit.Chip(c, fmt.Sprintf(text, moneyText(c, abs(row.Limit-row.Used), currency)), chip)
					})
					bar := ui.Box(c).FillWidth().Height(t.Space(2)).Radius(999).Background(t.Surface)
					fill := bar.AnimateTo("used", float32(min(1, used)), ui.Smooth)
					tint := bar.AnimateColor("ink", ink, ui.Snappy)
					bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
						if fill > 0 {
							p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*fill), H: r.H}, tint, 999)
						}
					}))
				})
			}
		})
	})
	return root
}
