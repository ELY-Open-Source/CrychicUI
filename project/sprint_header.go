package project

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SprintHeaderOptions are the sprint labels.
type SprintHeaderOptions struct{ Left, DateFormat string }

// SprintHeader shows sprint name, dates, and days left.
func SprintHeader(c *ui.Context, name string, start, end time.Time, o SprintHeaderOptions) ui.Element {
	if name == "" || !end.After(start) {
		panic("project: SprintHeader has invalid dates")
	}
	def(&o.Left, "%d days left")
	def(&o.DateFormat, "Jan 02")
	t := c.Theme()
	now := c.Now()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	left := max(0, int(end.Sub(now).Hours()/24))
	if left > 0 {
		next := end.Add(-time.Duration(left) * 24 * time.Hour)
		c.After(min(time.Second, next.Sub(now)))
	}
	root.Children(func() {
		ui.Text(c, name).Bold()
		kit.Muted(c, start.Format(o.DateFormat)+" – "+end.Format(o.DateFormat)).FontSize(t.Rem(0.75))
		kit.Numeric(c, fmt.Sprintf(o.Left, left)).TextColor(t.Accent)
	})
	return root
}
