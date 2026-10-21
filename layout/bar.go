package layout

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Bar rows start, center and end slots; nil stays empty.
func Bar(c *ui.Context, start, center, end func()) ui.Element {
	t := c.Theme()
	return ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(3)).Children(func() {
		for i, fill := range []func(){start, center, end} {
			slot := ui.Row(c).Key(i).Gap(t.Space(2))
			switch i {
			case 0:
				slot.Grow(1).Basis(0)
			case 2:
				slot.Grow(1).Basis(0).Justify(ui.End)
			}
			kit.Arrive(slot, ui.Pose{Fade: 1}, time.Duration(i)*60*time.Millisecond)
			if fill != nil {
				slot.Children(fill)
			}
		}
	})
}
