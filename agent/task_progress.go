package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TaskProgressOptions word the step count and the finish.
type TaskProgressOptions struct {
	Format, Finished string
}

// TaskProgress shows the current step and a filling bar.
func TaskProgress(c *ui.Context, label string, steps []string, current int, o TaskProgressOptions) ui.Element {
	if label == "" || len(steps) == 0 || current < 0 || current > len(steps) {
		panic(fmt.Sprintf("agent: TaskProgress %q needs a label, steps and a current step", label))
	}
	for _, s := range steps {
		if s == "" {
			panic(fmt.Sprintf("agent: TaskProgress %q steps must not be empty", label))
		}
	}
	def(&o.Format, "Step %d of %d")
	def(&o.Finished, "Finished")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Role(ui.RoleGroup).Label(label)
	finished := current == len(steps)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			if finished {
				mark(c, StepDone)
			} else {
				mark(c, StepRunning)
			}
			ui.Text(c, label).Bold().Grow(1).SingleLine()
			kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*float64(current)/float64(len(steps)))).Bold().TextColor(t.Accent)
		})
		feedback.SegmentProgress(c, label, len(steps), current)
		now := o.Finished
		if !finished {
			now = steps[current]
		}
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			ui.Text(c, now).FontSize(t.Rem(0.85)).Grow(1).SingleLine().ContentTransition(ui.ContentTransition{})
			kit.Numeric(c, fmt.Sprintf(o.Format, min(current+1, len(steps)), len(steps))).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		})
	})
}
