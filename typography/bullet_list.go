package typography

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BulletList staggers bulleted items into view.
func BulletList(c *ui.Context, items []string) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
		for i, item := range items {
			row := ui.Row(c).Key(i).Gap(t.Space(2))
			kit.Arrive(row, ui.Pose{Y: 4, Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				kit.Muted(c, "•")
				ui.Text(c, item).ContentTransition(ui.ContentTransition{})
			})
		}
	})
}
