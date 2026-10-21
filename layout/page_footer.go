package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PageFooter is a muted footer row; it fades in.
func PageFooter(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	foot := ui.Row(c).FillWidth().Gap(t.Space(4)).Padding(t.Space(4), 0).TextColor(t.TextMuted).FontSize(t.Rem(0.85))
	return kit.Arrive(foot, ui.Pose{Fade: 1}, 0).Children(content)
}
