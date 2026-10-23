package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SLAOptions are the SLA labels.
type SLAOptions struct{ Actual, Target string }

// SLA shows met percentage against target.
func SLA(c *ui.Context, actual, target float64, o SLAOptions) ui.Element {
	if actual < 0 || actual > 100 || target <= 0 || target > 100 || !finite(actual) || !finite(target) {
		panic("dashboard: SLA has invalid values")
	}
	def(&o.Actual, "%.2f%%")
	def(&o.Target, "Target %.2f%%")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Actual, actual)).Bold().Grow(1)
			kit.Muted(c, fmt.Sprintf(o.Target, target)).FontSize(t.Rem(0.75))
		})
		progress(c, min(1, actual/target), t.Success)
	})
	return root
}
