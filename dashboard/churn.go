package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ChurnOptions word the card; Change takes the point change.
type ChurnOptions struct{ Title, Change, Note string }

// Churn shows the rate and trend; a rise reads red.
func Churn(c *ui.Context, rate float64, points []Point, o ChurnOptions) ui.Element {
	if rate < 0 || rate > 100 || !finite(rate) || len(points) < 2 {
		panic("dashboard: Churn needs a rate in [0,100] and two points")
	}
	for _, point := range points {
		if point.Value < 0 || point.Value > 100 {
			panic(fmt.Sprintf("dashboard: Churn point %q must be in [0,100]", point.Label))
		}
	}
	def(&o.Title, "Churn rate")
	def(&o.Change, "%+.1f pts")
	t := c.Theme()
	change := points[len(points)-1].Value - points[len(points)-2].Value
	ink, glyph := t.Success, "trending-down"
	if change > 0 {
		ink, glyph = t.Danger, "trending-up"
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "user-minus", o.Title, func() { tone(c, fmt.Sprintf(o.Change, change), ink, glyph) })
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			ui.Column(c).Gap(t.Space(0.5)).Children(func() {
				kit.Numeric(c, fmt.Sprintf("%.1f%%", rate)).Bold().FontSize(t.Rem(1.8)).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				if o.Note != "" {
					kit.Muted(c, o.Note).FontSize(t.Rem(0.75)).Shrink(0)
				}
			})
			series(c, "Churn", points, t.Danger).Height(t.Space(20))
			ui.Row(c).FillWidth().Children(func() {
				kit.Muted(c, points[0].Label).FontSize(t.Rem(0.7)).Grow(1)
				kit.Muted(c, points[len(points)-1].Label).FontSize(t.Rem(0.7)).Shrink(0)
			})
		})
	})
	return root
}
