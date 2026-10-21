package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Persona is an avatar, name and role fading in.
func Persona(c *ui.Context, name, role string, image *ui.Bitmap) ui.Element {
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
		avatarCircle(c, name, image, t.Space(8))
		ui.Column(c).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, name).Bold()
			kit.Muted(c, role).FontSize(t.Rem(0.85))
		})
	})
	return kit.Arrive(root, ui.Pose{Fade: 1}, 0)
}
