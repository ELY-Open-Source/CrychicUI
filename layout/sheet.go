package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Sheet is a centered paper page; it rises in.
func Sheet(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	page := kit.Card(c).FillWidth().MaxWidth(t.Rem(48)).Margin(0, ui.Auto).Gap(t.Space(3)).Padding(t.Space(10))
	return kit.Arrive(page, ui.Pose{Y: 12, Fade: 1}, 0).Children(content)
}
