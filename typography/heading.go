package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Heading is a leveled heading fading in.
func Heading(c *ui.Context, text string, level int) ui.Element {
	if level < 1 || level > 6 {
		panic("typography: Heading level must be 1 through 6")
	}
	sizes := [...]float32{2, 1.5, 1.25, 1.1, 1, 0.85}
	e := ui.Text(c, text).Role(ui.RoleHeading).Level(level).Bold().FontSize(c.Theme().Rem(sizes[level-1]))
	return kit.Arrive(e, ui.Pose{Fade: 1}, 0)
}
