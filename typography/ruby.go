package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Ruby places a fading annotation over base text.
func Ruby(c *ui.Context, base, annotation string) ui.Element {
	t := c.Theme()
	return ui.Column(c).AlignSelf(ui.Start).AlignItems(ui.Center).Children(func() {
		note := kit.Muted(c, annotation).FontSize(t.Rem(0.75))
		kit.Arrive(note, ui.Pose{Fade: 1}, 0)
		ui.Text(c, base)
	})
}
