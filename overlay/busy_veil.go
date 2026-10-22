package overlay

import (
	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/egoist/mygo/ui"
)

// BusyVeil covers content with a spinner and message while busy.
func BusyVeil(c *ui.Context, busy bool, message string, content func()) ui.Element {
	if message == "" || content == nil {
		panic("overlay: BusyVeil needs a message and content")
	}
	t := c.Theme()
	root := ui.Box(c)
	cover(c, root, busy, 0, content, func() {
		ui.Column(c).FillWidth().FillHeight().Center().Gap(t.Space(2)).Padding(t.Space(4)).Radius(t.Radius).
			Background(t.Background.Alpha(0.85)).Role(ui.RoleStatus).Label(message).Children(func() {
			feedback.Spin(c, message, feedback.SpinOptions{})
			ui.Text(c, message).TextColor(t.TextMuted)
		})
	})
	moved(root, busy, "BusyVeil")
	return root
}
