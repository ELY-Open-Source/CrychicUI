package schedule

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NowLine shows a line at the current time.
func NowLine(c *ui.Context, now time.Time) ui.Element {
	if now.IsZero() {
		panic("schedule: NowLine needs current time")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().AlignItems(ui.Center).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{X: -t.Space(4), Fade: 1}, 0)
	c.After(time.Second - time.Duration(now.Nanosecond()))
	root.Children(func() {
		ui.Row(c).Shrink(0).Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.Danger).Children(func() {
			kit.Numeric(c, now.Format("15:04:05")).FontSize(t.Rem(0.7)).Bold().TextColor(t.AccentText)
		})
		ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(t.Danger).
			Border(t.Space(0.5), t.Background).MarginStart(-t.Space(0.5))
		line := ui.Box(c).Grow(1).Height(t.Space(0.5)).MarginStart(-t.Space(0.5)).Background(t.Danger)
		kit.Arrive(line, ui.Pose{X: -t.Space(8), Fade: 1}, 120*time.Millisecond)
	})
	return root
}
