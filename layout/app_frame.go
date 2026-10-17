package layout

import (
	"time"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AppFrame stacks top bar, content, status bar; bars slide.
func AppFrame(c *ui.Context, top, content, status func()) ui.Element {
	t := c.Theme()
	return ui.Column(c).Fill().Children(func() {
		bar := ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(4)).AlignItems(ui.Center)
		kit.Arrive(bar, ui.Pose{Y: -8, Fade: 1}, 0).Children(top)
		ui.Divider(c)
		ui.Column(c).Grow(1).Children(content)
		ui.Divider(c)
		foot := ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(1.5), t.Space(4)).AlignItems(ui.Center)
		kit.Arrive(foot, ui.Pose{Y: 8, Fade: 1}, 80*time.Millisecond).Children(status)
	})
}
