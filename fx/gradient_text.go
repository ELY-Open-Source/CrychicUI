package fx

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// GradientText fills text from accent to success; it rises in.
func GradientText(c *ui.Context, text string) ui.Element {
	s := headline(c, text)
	t := c.Theme()
	o := shader.Options{Uniforms: map[string]any{"A": t.Accent, "B": t.Success}}
	e := filled(c, s).Draw(func(p *ui.Painter, r ui.Rect) {
		paint.Fill(p, r, o, func() { p.RichText(r.X, r.Y, 0, s) })
	})
	return kit.Arrive(e, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}
