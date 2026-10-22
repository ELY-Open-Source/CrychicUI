package motion

import (
	"strconv"

	"github.com/egoist/mygo/ui"
)

// RollingDigits shows value in large digits that roll on change.
func RollingDigits(c *ui.Context, value int64) ui.Element {
	t := c.Theme()
	root := ui.Box(c).FontSize(t.Rem(2)).Bold()
	prev := ui.Local(root, "prev", func() int64 { return value })
	down := ui.Local(root, "down", func() bool { return false })
	if value != *prev {
		*down, *prev = value < *prev, value
	}
	return root.Children(func() {
		ui.Text(c, strconv.FormatInt(value, 10)).FontFeatures("tnum").ContentTransition(ui.ContentTransition{Numeric: true, Down: *down})
	})
}
