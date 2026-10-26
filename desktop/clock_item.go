package desktop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ClockItemOptions set the time layout.
type ClockItemOptions struct{ Format string }

// ClockItem is a menu bar clock; it turns each minute.
func ClockItem(c *ui.Context, o ClockItemOptions) ui.Element {
	def(&o.Format, "Mon Jan 2  15:04")
	t := c.Theme()
	now := minute(c)
	root := ui.Row(c).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Role(ui.RoleStatus)
	kit.Arrive(root, ui.Pose{Y: -t.Space(1), Fade: 1}, 0)
	return root.Children(func() {
		kit.Numeric(c, now.Format(o.Format)).FontSize(t.Rem(0.78)).Shrink(0).MaxWidth(t.Space(60)).SingleLine()
	})
}
