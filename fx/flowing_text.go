package fx

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// FlowingText fills text with gradient stops that keep moving.
func FlowingText(c *ui.Context, text string) ui.Element {
	s := headline(c, text)
	t := c.Theme()
	e := filled(c, s)
	o := shader.Options{Uniforms: map[string]any{"A": t.Accent, "B": t.Success, "Shift": phase(c, e, 4*time.Second), "Flow": 1}}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		paint.Fill(p, r, o, func() { p.RichText(r.X, r.Y, 0, s) })
	})
}
