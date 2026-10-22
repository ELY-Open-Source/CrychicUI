package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Plan is a priced plan card's copy.
type Plan struct {
	Name, Price, Period, Note string
}

// PlanChoice is plan radio cards with prices.
func PlanChoice(c *ui.Context, label string, plans []Plan, selected *string) ui.Element {
	names := make([]string, len(plans))
	for i, p := range plans {
		if p.Price == "" {
			panic("choice: PlanChoice plans need a price")
		}
		names[i] = p.Name
	}
	distinct("PlanChoice", names)
	if label == "" || selected == nil || !slices.Contains(names, *selected) {
		panic("choice: PlanChoice needs a label and a valid selection")
	}
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for _, p := range plans {
			c.Key(p.Name)
			b := ui.RadioBase(c, selected, p.Name).Label(p.Name).Column().Grow(1).Basis(0).MinWidth(t.Space(40)).
				AlignItems(ui.Stretch).Gap(t.Space(2)).Padding(t.Space(4)).Radius(t.Radius + 2)
			on := *selected == p.Name
			edge, fill := t.Border, t.Background
			if on {
				edge, fill = t.Accent, t.Accent.Alpha(0.05)
			} else if b.Hovered() {
				edge = t.Border.Mix(t.Text, 0.25)
			}
			b.Background(b.AnimateColor("fill", fill, ui.Snappy)).Border(1, b.AnimateColor("edge", edge, ui.Snappy)).
				Shadow(0, 1, 3, 0, kit.Shade(0.06))
			kit.Lift(b).Children(func() {
				ui.Row(c).Gap(t.Space(2)).Children(func() {
					dot(c, b, on)
					ui.Text(c, p.Name).Bold()
				})
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
					ui.Text(c, p.Price).FontSize(t.Rem(1.8)).Bold()
					if p.Period != "" {
						kit.Muted(c, p.Period).Padding(0, 0, t.Space(1), 0)
					}
				})
				if p.Note != "" {
					kit.Muted(c, p.Note).FontSize(t.Rem(0.85))
				}
			})
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().Wrap().FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
	if changed {
		root.SetChanged()
		kit.Log("PlanChoice", "selected", "index", slices.Index(names, *selected))
	}
	return root
}
