package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Plan is one billing plan; Period and Note are optional.
type Plan struct {
	Name, Price, Period, Note string
	Features                  []string
}

// PlanCardOptions word the card.
type PlanCardOptions struct{ Current, Change string }

// PlanCard shows the current plan with change; cue Tap.
func PlanCard(c *ui.Context, plan Plan, current bool, change *bool, o PlanCardOptions) ui.Element {
	if plan.Name == "" || plan.Price == "" || change == nil {
		panic(fmt.Sprintf("account: PlanCard %q needs a name, price and state", plan.Name))
	}
	def(&o.Current, "Current plan")
	def(&o.Change, "Change plan")
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5))
	root.Border(1, root.AnimateColor("line", t.Border.Mix(t.Accent, unit(current)), ui.Snappy))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			tile(c, "sparkles", t.Accent, t.Space(9))
			ui.Text(c, plan.Name).FontSize(t.Rem(1.05)).Bold().Grow(1).MinWidth(0).SingleLine()
			if current {
				kit.Chip(c, o.Current, t.Accent)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
				kit.Numeric(c, plan.Price).FontSize(t.Rem(1.9)).Bold().Shrink(0)
				if plan.Period != "" {
					kit.Muted(c, plan.Period).FontSize(t.Rem(0.85)).Padding(0, 0, t.Space(1), 0).Shrink(0)
				}
			})
			if plan.Note != "" {
				kit.Muted(c, plan.Note).FontSize(t.Rem(0.78)).Wrap().FillWidth()
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(2))
		list.Children(func() {
			for i, feature := range plan.Features {
				if feature == "" {
					panic(fmt.Sprintf("account: PlanCard %q feature must not be empty", plan.Name))
				}
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Radius(999).Background(t.Success.Alpha(0.14)).
						TextColor(t.Success).Children(func() { glyph(c, "check", 12) })
					ui.Text(c, feature).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
				})
			}
		})
		b := outline(c, o.Change, o.Change, "arrow-right-left", t.Text).FillWidth()
		if b.Clicked() {
			*change = true
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("PlanCard", "change")
		}
	})
	return root
}
