package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Lead is large introductory text fading in.
func Lead(c *ui.Context, text string) ui.Element {
	e := ui.Text(c, text).FontSize(c.Theme().Rem(1.25))
	return kit.Arrive(e, ui.Pose{Fade: 1}, 0)
}
