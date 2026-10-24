package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Plan is one billing plan.
type Plan struct {
	Name, Price string
	Features    []string
}

// PlanCardOptions are the plan labels.
type PlanCardOptions struct{ Current, Change string }

// PlanCard shows the current plan with change; cue Tap.
func PlanCard(c *ui.Context, plan Plan, current bool, change *bool, o PlanCardOptions) ui.Element {
	if plan.Name == "" || plan.Price == "" || change == nil {
		panic("account: PlanCard needs plan and state")
	}
	def(&o.Current, "Current plan")
	def(&o.Change, "Change plan")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).Children(func() {
			ui.Text(c, plan.Name).Bold().Grow(1)
			if current {
				kit.Muted(c, o.Current)
			}
		})
		ui.Text(c, plan.Price)
		ui.Column(c).Gap(t.Space(1)).Children(func() {
			for _, feature := range plan.Features {
				if feature == "" {
					panic("account: plan feature must not be empty")
				}
				ui.Text(c, feature).Wrap()
			}
		})
		b := button(c, o.Change)
		b.Children(func() { ui.Text(c, o.Change) })
		if b.Clicked() {
			*change = true
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("PlanCard", "change")
		}
	})
	return root
}
