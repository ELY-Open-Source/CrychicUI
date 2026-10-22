package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SkeletonCard is a sheen image block, title and lines.
func SkeletonCard(c *ui.Context) ui.Element {
	t := c.Theme()
	root := kit.Card(c).FillWidth().MaxWidth(t.Space(80)).Padding(t.Space(3)).Gap(t.Space(3)).Role(ui.RoleNone)
	return root.Children(func() {
		sheen(c, ui.Box(c).FillWidth().AspectRatio(16.0/9), t.Radius)
		sheen(c, ui.Box(c).WidthPercent(50).Height(t.Space(4)), t.Space(2))
		SkeletonText(c, 2)
	})
}
