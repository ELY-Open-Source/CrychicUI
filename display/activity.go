package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Activity shows actor, action, object and time fading in.
func Activity(c *ui.Context, actor, action, object, when string) ui.Element {
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
		ui.Text(c, actor).Bold()
		kit.Muted(c, action)
		ui.Text(c, object)
		ui.Box(c).Grow(1)
		kit.Muted(c, when).FontSize(t.Rem(0.85))
	})
	return kit.Arrive(root, ui.Pose{Fade: 1}, 0)
}
