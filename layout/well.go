package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Well is a sunken box; its content fades in.
func Well(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	return ui.Column(c).Padding(t.Space(4)).Background(t.Surface).Radius(t.Radius).
		InsetShadow(0, 1, 4, 0, kit.Shade(0.12)).Children(func() {
		kit.Arrive(ui.Column(c).Gap(t.Space(2)), ui.Pose{Fade: 1}, 0).Children(content)
	})
}
