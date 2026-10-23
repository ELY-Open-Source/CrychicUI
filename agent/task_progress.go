package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TaskProgressOptions word the step count.
type TaskProgressOptions struct {
	Format string
}

// TaskProgress shows the current step and a filling bar.
func TaskProgress(c *ui.Context, label string, step, total int, o TaskProgressOptions) ui.Element {
	if label == "" || total <= 0 || step < 0 || step > total {
		panic(fmt.Sprintf("agent: TaskProgress %q needs a label, a total and a step", label))
	}
	def(&o.Format, "Step %d of %d")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	return root.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, label).Bold().Grow(1).SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Format, step, total)).FontSize(t.Rem(0.85)).TextColor(t.TextMuted)
		})
		feedback.ProgressBar(c, label, float64(step)/float64(total))
	})
}
