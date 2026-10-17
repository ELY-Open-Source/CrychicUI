package layout

import (
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Centered centers content in its room; the content scales in.
func Centered(c *ui.Context, content func()) ui.Element {
	return ui.Column(c).Fill().Center().Children(func() {
		kit.Arrive(ui.Box(c), ui.Pose{Scale: 0.96, Fade: 1}, 0).Children(content)
	})
}
