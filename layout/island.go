package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var islandSpring = ui.Snappy

// Island is an inverted pill; its size springs.
func Island(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	return ui.Row(c).Gap(t.Space(2)).Padding(t.Space(2), t.Space(4)).Radius(999).ClipX().
		Background(t.Text).TextColor(t.Background).Shadow(0, 1, 3, 0, kit.Shade(0.06)).
		Transition(ui.ElementTransition{Size: true, Animation: &islandSpring}).Children(content)
}
