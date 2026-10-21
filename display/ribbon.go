package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Ribbon slides a corner label over its content.
func Ribbon(c *ui.Context, label string, content func()) ui.Element {
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Clip()
	return root.Children(func() {
		content()
		ribbon := ui.Row(c).Absolute().Top(t.Space(2)).Right(-t.Space(4)).Padding(t.Space(0.5), t.Space(4)).
			Background(t.Accent).TextColor(t.AccentText).Rotate(45)
		kit.Arrive(ribbon, ui.Pose{X: 12, Fade: 1}, 0).Children(func() { ui.Text(c, label).FontSize(t.Rem(0.75)).Bold() })
	})
}
