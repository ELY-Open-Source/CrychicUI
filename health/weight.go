package health

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WeightPoint is a labelled weigh-in.
type WeightPoint struct {
	Label string
	Value float64
}

// WeightOptions word the card; Goal 0 draws no goal line.
type WeightOptions struct {
	Title, Unit, Change, Target string
	Goal                        float64
}

// Weight draws a weight trend with its change and goal.
func Weight(c *ui.Context, points []WeightPoint, o WeightOptions) ui.Element {
	def(&o.Title, "Weight")
	if len(points) < 2 || len(points) > 120 || !(o.Goal == 0 || o.Goal > 0 && o.Goal < 1000) {
		panic(fmt.Sprintf("health: Weight %q needs 2 to 120 points and a goal of 0 to 1000", o.Title))
	}
	values := make([]float64, len(points))
	lo, hi := math.Inf(1), math.Inf(-1)
	for i, p := range points {
		if strings.TrimSpace(p.Label) == "" || !(p.Value > 0 && p.Value < 1000) {
			panic(fmt.Sprintf("health: Weight %q point %d needs a label and a value of 0 to 1000", o.Title, i))
		}
		values[i] = p.Value
		lo, hi = min(lo, p.Value), max(hi, p.Value)
	}
	if o.Goal > 0 {
		lo, hi = min(lo, o.Goal), max(hi, o.Goal)
	}
	pad := max(0.5, (hi-lo)*0.15)
	def(&o.Unit, "kg")
	def(&o.Change, "%+.1f %s")
	def(&o.Target, "Goal %.1f %s")
	t := c.Theme()
	first, last := values[0], values[len(values)-1]
	delta := math.Round((last-first)*10) / 10
	ink := t.TextMuted
	if o.Goal > 0 && delta != 0 {
		ink = t.Warning
		if math.Abs(last-o.Goal) < math.Abs(first-o.Goal) {
			ink = t.Success
		}
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			kit.Arrive(tile(c, "scale", t.Accent, t.Space(11)), ui.Pose{Scale: 0.6, Fade: 1}, 0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.End).Children(func() {
					kit.Numeric(c, fmt.Sprintf("%.1f", last)).FontSize(t.Rem(1.6)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Muted(c, o.Unit).FontSize(t.Rem(0.78)).Shrink(0)
				})
			})
			kit.Chip(c, fmt.Sprintf(o.Change, delta, o.Unit), ink)
		})
		guide := math.NaN()
		if o.Goal > 0 {
			guide = o.Goal
		}
		trend(c, values, lo-pad, hi+pad, t.Accent, t.Space(24), guide)
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			kit.Muted(c, points[0].Label).FontSize(t.Rem(0.66)).Shrink(0).MaxWidth(t.Space(24)).SingleLine()
			ui.Box(c).Grow(1)
			if o.Goal > 0 {
				kit.Muted(c, fmt.Sprintf(o.Target, o.Goal, o.Unit)).FontSize(t.Rem(0.66)).Bold().Shrink(0).SingleLine()
				ui.Box(c).Grow(1)
			}
			kit.Muted(c, points[len(points)-1].Label).FontSize(t.Rem(0.66)).Shrink(0).MaxWidth(t.Space(24)).SingleLine()
		})
	})
	return root
}
