package feedback

import "github.com/egoist/mygo/ui"

// SkeletonList is rows of a sheen avatar and two lines.
func SkeletonList(c *ui.Context, rows int) ui.Element {
	if rows < 1 {
		panic("feedback: SkeletonList needs a row")
	}
	t := c.Theme()
	d := t.Space(10)
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleNone)
	return root.Children(func() {
		for i := range rows {
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(3)).Children(func() {
				sheen(c, ui.Box(c).Size(d, d).Shrink(0), d/2)
				ui.Column(c).Grow(1).Gap(t.Space(2)).Children(func() {
					sheen(c, ui.Box(c).WidthPercent(40).Height(t.Space(3)), t.Space(1.5))
					sheen(c, ui.Box(c).WidthPercent(80).Height(t.Space(3)), t.Space(1.5))
				})
			})
		}
	})
}
