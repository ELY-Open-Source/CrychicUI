package layout

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Ladder alternates rows around a center line; rows rise.
func Ladder(c *ui.Context, rows []func()) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Children(func() {
		for i, row := range rows {
			r := ui.Row(c).Key(i).FillWidth().AlignItems(ui.Stretch)
			kit.Arrive(r, ui.Pose{Y: 8, Fade: 1}, time.Duration(i)*60*time.Millisecond).Children(func() {
				side := func(key string, fill bool) {
					s := ui.Column(c).Key(key).Grow(1).Basis(0).Padding(t.Space(2), t.Space(3))
					if key == "start" {
						s.AlignItems(ui.End)
					}
					if fill {
						s.Children(row)
					}
				}
				side("start", i%2 == 0)
				ui.Box(c).Width(1).Background(t.Border)
				side("end", i%2 == 1)
			})
		}
	})
}
