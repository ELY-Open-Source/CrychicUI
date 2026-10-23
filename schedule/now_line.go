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
	root := ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	c.After(time.Second)
	root.Children(func() {
		ui.Box(c).Grow(1).Height(1).Background(t.Danger)
		kit.Numeric(c, now.Format("15:04:05")).FontSize(t.Rem(0.75)).TextColor(t.Danger)
	})
	return root
}
