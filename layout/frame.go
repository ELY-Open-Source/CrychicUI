package layout

import (
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Frame clips content to the theme radius; content fades in.
func Frame(c *ui.Context, content func()) ui.Element {
	return ui.Box(c).Clip().Radius(c.Theme().Radius).Children(func() {
		kit.Arrive(ui.Box(c), ui.Pose{Fade: 1}, 0).Children(content)
	})
}
