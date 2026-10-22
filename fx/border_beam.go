package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// BorderBeam runs a light segment around content's border.
func BorderBeam(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("fx: BorderBeam needs content")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(1).AlignSelf(ui.Start).Radius(t.Radius).Border(1, t.Border)
	turn := 360 * phase(c, root, 4*time.Second)
	none := t.Accent.Alpha(0)
	return root.Children(func() {
		ui.Box(c).ColumnStart(1).RowStart(1).Padding(t.Space(4)).Children(content)
		ui.Box(c).ColumnStart(1).RowStart(1).Radius(t.Radius).Border(t.Space(0.5), ui.Color{}).PassThrough().
			BorderConicGradient(ui.ConicGradient{Angle: turn, Stops: []ui.GradientStop{at(none, 0), at(none, 0.8),
				at(t.Accent, 0.95), at(none, 1)}})
	})
}
