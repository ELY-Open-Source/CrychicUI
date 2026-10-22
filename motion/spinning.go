package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Spinning turns content continuously.
func Spinning(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: Spinning needs content")
	}
	root := ui.Box(c)
	return root.Rotate(root.Loop("crychic.spin", 1600*time.Millisecond, ui.Linear) * 360).Children(content)
}
