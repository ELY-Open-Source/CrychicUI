package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PlanStep is one step of a Plan.
type PlanStep struct {
	Title, Detail string
	State         StepState
}

// PlanOptions word the done count.
type PlanOptions struct {
	Format string
}

// Plan shows steps with pending, running and done marks.
func Plan(c *ui.Context, label string, steps []PlanStep, o PlanOptions) ui.Element {
	if label == "" || len(steps) == 0 {
		panic(fmt.Sprintf("agent: Plan %q needs a label and steps", label))
	}
	done := 0
	for _, s := range steps {
		if s.Title == "" || s.State > StepFailed {
			panic(fmt.Sprintf("agent: Plan %q step %q needs a title and a known state, got %d", label, s.Title, s.State))
		}
		if s.State == StepDone {
			done++
		}
	}
	def(&o.Format, "%d of %d done")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3), t.Space(4)).Radius(t.Radius).
		Background(t.Background).Border(1, t.Border).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, label).Bold().Grow(1).SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Format, done, len(steps))).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		})
		ui.Column(c).FillWidth().Children(func() {
			for i, s := range steps {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch).Role(ui.RoleListItem).Label(s.Title)
				kit.Arrive(row, ui.Pose{Y: t.Space(1.5), Fade: 1}, time.Duration(i)*40*time.Millisecond)
				row.Children(func() {
					ui.Column(c).AlignItems(ui.Center).Children(func() {
						mark(c, s.State)
						if i < len(steps)-1 {
							rail := ui.Box(c).Width(2).Grow(1).MinHeight(t.Space(3)).Radius(1)
							rail.Background(rail.AnimateColor("rail", t.Border.Mix(t.Success, unit(s.State == StepDone)), ui.Smooth))
						}
					})
					ui.Column(c).Grow(1).Gap(t.Space(0.5)).Padding(0, 0, t.Space(3), 0).Children(func() {
						title := ui.Text(c, s.Title).Wrap()
						switch s.State {
						case StepRunning:
							title.Bold()
						case StepDone:
							title.TextColor(t.TextMuted)
						case StepFailed:
							title.TextColor(t.Danger)
						}
						if s.Detail != "" {
							kit.Muted(c, s.Detail).FontSize(t.Rem(0.8)).Wrap()
						}
					})
				})
			}
		})
	})
}
