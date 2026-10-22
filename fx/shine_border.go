package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ShineBorder passes a sheen along content's border.
func ShineBorder(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("fx: ShineBorder needs content")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(1).AlignSelf(ui.Start).Radius(t.Radius).Border(1, t.Border)
	mid := -0.2 + 1.4*phase(c, root, 3*time.Second)
	none, w := t.Accent.Alpha(0), float32(0.15)
	inside := func(v float32) float32 { return min(max(v, 0), 1) }
	return root.Children(func() {
		ui.Box(c).ColumnStart(1).RowStart(1).Padding(t.Space(4)).Children(content)
		ui.Box(c).ColumnStart(1).RowStart(1).Radius(t.Radius).Border(t.Space(0.5), ui.Color{}).PassThrough().
			BorderLinearGradient(ui.LinearGradient{Angle: 120, Stops: []ui.GradientStop{at(none, 0), at(none, inside(mid-w)),
				at(t.Accent, inside(mid)), at(none, inside(mid+w)), at(none, 1)}})
	})
}
