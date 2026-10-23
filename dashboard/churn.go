package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ChurnOptions are the churn labels.
type ChurnOptions struct{ Title string }

// Churn shows churn rate with trend.
func Churn(c *ui.Context, rate float64, points []Point, o ChurnOptions) ui.Element {
	if rate < 0 || rate > 100 || !finite(rate) || len(points) < 2 {
		panic("dashboard: Churn has invalid data")
	}
	for _, point := range points {
		if point.Value < 0 || point.Value > 100 {
			panic("dashboard: Churn history must be 0–100")
		}
	}
	def(&o.Title, "Churn rate")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Title).Bold()
		kit.Numeric(c, fmt.Sprintf("%.1f%%", rate)).FontSize(t.Rem(1.4))
		series(c, points, t.Danger)
	})
	return root
}
