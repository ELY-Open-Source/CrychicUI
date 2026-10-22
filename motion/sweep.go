package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var sweep = replay(500*time.Millisecond, ui.EaseInOut)

// Sweep marks text on first build and each run change.
func Sweep(c *ui.Context, run uint64, text string) ui.Element {
	if text == "" {
		panic("motion: Sweep needs text")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(1).AlignSelf(ui.Start)
	k := root.Keyframes("crychic.sweep", sweep, ui.KeyframeOptions{Run: run})
	return root.Children(func() {
		ui.Box(c).ColumnStart(1).RowStart(1).JustifySelf(ui.Start).AlignSelf(ui.End).HeightPercent(45).
			WidthPercent(100 * k).Radius(t.Space(0.5)).Background(t.Warning.Alpha(0.4))
		ui.Text(c, text).ColumnStart(1).RowStart(1)
	})
}
