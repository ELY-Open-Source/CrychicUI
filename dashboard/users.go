package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// UsersOptions are the user metric labels.
type UsersOptions struct{ Title string }

// Users shows active users with a chart.
func Users(c *ui.Context, points []Point, o UsersOptions) ui.Element {
	if len(points) < 2 {
		panic("dashboard: Users needs points")
	}
	for _, point := range points {
		if point.Value < 0 {
			panic("dashboard: Users count must not be negative")
		}
	}
	def(&o.Title, "Active users")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Title).Bold()
		kit.Numeric(c, fmt.Sprintf("%.0f", points[len(points)-1].Value)).FontSize(t.Rem(1.4))
		series(c, points, t.Accent)
	})
	return root
}
