package layout

import "github.com/egoist/mygo/ui"

// Grouped joins rows with separators; row changes spring.
func Grouped(c *ui.Context, rows []func()) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Border(1, t.Border).Radius(t.Radius).Clip().Background(t.Background).
		Children(func() {
			for i, row := range rows {
				ui.Column(c).Key(i).FillWidth().Transition(collapseMotion).Children(func() {
					if i > 0 {
						ui.Box(c).FillWidth().Height(1).Background(t.Border)
					}
					ui.Column(c).Padding(t.Space(2), t.Space(3)).Children(row)
				})
			}
		})
}
