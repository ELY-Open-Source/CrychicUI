package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Caret is a blinking text caret one line tall.
func Caret(c *ui.Context) ui.Element {
	t := c.Theme()
	e := ui.Box(c).Size(t.Space(0.5), t.Rem(1.25)).Shrink(0).Radius(t.Space(0.25)).Background(t.Accent)
	return e.Opacity(unit(phase(c, e, 1060*time.Millisecond) < 0.55))
}
