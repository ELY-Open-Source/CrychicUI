package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Card is the shared card surface; it lifts on hover.
func Card(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	return kit.Lift(kit.Card(c).Gap(t.Space(2)).Padding(t.Space(4))).Children(content)
}
