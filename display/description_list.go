package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DescriptionList fades in key and value rows with dividers.
func DescriptionList(c *ui.Context, items []KeyValuePair) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Border(1, t.Border).Radius(t.Radius).Children(func() {
		for i, item := range items {
			row := ui.Column(c).Key(i).FillWidth()
			kit.Arrive(row, ui.Pose{Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				if i > 0 {
					ui.Box(c).FillWidth().Height(1).Background(t.Border)
				}
				ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(3)).Children(func() {
					kit.Muted(c, item.Key).Width(t.Rem(10))
					ui.Text(c, item.Value).Grow(1)
				})
			})
		}
	})
}
