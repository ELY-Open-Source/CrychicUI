package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SavingsGoalOptions are the savings labels.
type SavingsGoalOptions struct{ Target string }

// SavingsGoal shows saved progress and target date.
func SavingsGoal(c *ui.Context, name string, saved, target float64, date, currency string, o SavingsGoalOptions) ui.Element {
	if name == "" || saved < 0 || target <= 0 || !finite(saved) || !finite(target) || date == "" || currency == "" {
		panic("finance: SavingsGoal has invalid data")
	}
	def(&o.Target, "Target %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, name).Bold().Grow(1)
			kit.Muted(c, fmt.Sprintf(o.Target, date)).FontSize(t.Rem(0.75))
		})
		progress(c, min(1, saved/target), t.Success)
		ui.Row(c).FillWidth().Children(func() {
			money(c, saved, currency)
			kit.Muted(c, date).FontSize(t.Rem(0.75))
		})
	})
	return root
}
