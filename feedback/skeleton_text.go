package feedback

import "github.com/egoist/mygo/ui"

// SkeletonText is n sheen lines, the last one shorter.
func SkeletonText(c *ui.Context, lines int) ui.Element {
	if lines < 1 {
		panic("feedback: SkeletonText needs a line")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleNone)
	return root.Children(func() {
		for i := range lines {
			c.Key(i)
			line := ui.Box(c).Height(t.Space(3)).WidthPercent(100)
			if i == lines-1 && lines > 1 {
				line.WidthPercent(60)
			}
			sheen(c, line, t.Space(1.5))
		}
	})
}
