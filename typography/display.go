package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Display is hero text fading up on arrival.
func Display(c *ui.Context, text string) ui.Element {
	e := ui.Text(c, text).Bold().FontSize(c.Theme().Rem(3))
	return kit.Arrive(e, ui.Pose{Y: 8, Fade: 1}, 0)
}
