package health

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CaloriesOptions word the card's ring and stats.
type CaloriesOptions struct{ Title, Left, Over, Eaten, Burned, Goal, Unit string }

// Calories rings what is left after eaten and burned.
func Calories(c *ui.Context, eaten, burned, goal int, o CaloriesOptions) ui.Element {
	def(&o.Title, "Calories")
	if eaten < 0 || burned < 0 || goal < 1 || eaten > 1e6 || burned > 1e6 || goal > 1e6 {
		panic(fmt.Sprintf("health: Calories %q needs eaten and burned of 0 to 1,000,000 and a goal of 1 or more", o.Title))
	}
	def(&o.Left, "%s left")
	def(&o.Over, "%s over")
	def(&o.Eaten, "Eaten")
	def(&o.Burned, "Burned")
	def(&o.Goal, "Goal")
	def(&o.Unit, "kcal")
	t := c.Theme()
	net := eaten - burned
	left := goal - net
	ink, word := t.Success, fmt.Sprintf(o.Left, number(c, left))
	if left < 0 {
		ink, word = t.Danger, fmt.Sprintf(o.Over, number(c, -left))
	}
	root := card(c).Row().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		dial(c, t.Space(30), float32(max(0, min(1, float64(net)/float64(goal)))), ink, func() {
			ui.Column(c).Gap(0).AlignItems(ui.Center).Children(func() {
				kit.Numeric(c, word).FontSize(t.Rem(0.95)).Bold().MaxWidth(t.Space(24)).SingleLine().TextColor(ink).
					ContentTransition(ui.ContentTransition{Numeric: true})
				kit.Muted(c, o.Unit).FontSize(t.Rem(0.66)).Shrink(0)
			})
		})
		stats := ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(2.5))
		stats.Children(func() {
			for i, s := range []struct {
				name, mark string
				n          int
				ink        ui.Color
			}{{o.Goal, "target", goal, t.TextMuted}, {o.Eaten, "utensils", eaten, t.Warning}, {o.Burned, "flame", burned, t.Danger}} {
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
				kit.Stagger(c, stats, row, i).Children(func() {
					glyph(c, s.mark, 16).TextColor(s.ink)
					kit.Muted(c, s.name).FontSize(t.Rem(0.76)).Grow(1).MinWidth(0).SingleLine()
					kit.Numeric(c, number(c, s.n)).FontSize(t.Rem(0.88)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
				})
			}
		})
	})
	return root
}
