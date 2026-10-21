package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Eyebrow is an accent label sliding in.
func Eyebrow(c *ui.Context, text string) ui.Element {
	t := c.Theme()
	e := ui.Text(c, text).TextColor(t.Accent).Bold().FontSize(t.Rem(0.85))
	return kit.Arrive(e, ui.Pose{X: -8, Fade: 1}, 0)
}
