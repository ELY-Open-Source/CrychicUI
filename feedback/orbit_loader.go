package feedback

import "github.com/egoist/mygo/ui"

// OrbitLoader is a dot orbiting a ring.
func OrbitLoader(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("feedback: OrbitLoader needs a label")
	}
	t := c.Theme()
	size := t.Space(8)
	return ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleProgress).Label(label).Range(0, 1, -1).Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		cx, cy, d := r.X+r.W/2, r.Y+r.H/2, t.Space(2)
		rad := r.W/2 - d/2
		var ring ui.Path
		p.StrokePath(ring.Circle(cx, cy, rad), 1, t.Border)
		x, y := polar(cx, cy, rad, float64(p.Now().UnixMilli()%1100)/1100*360-90)
		p.Fill(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, t.Accent, d/2)
	})
}
