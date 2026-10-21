package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Showcase centers content on a surface stage; content scales in.
func Showcase(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Padding(t.Space(10)).Center().Radius(t.Radius).Background(t.Surface).
		Children(func() { kit.Arrive(ui.Box(c), ui.Pose{Scale: 0.96, Fade: 1}, 0).Children(content) })
}
