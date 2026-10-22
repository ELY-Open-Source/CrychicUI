package feedback

import "github.com/egoist/mygo/ui"

// IndeterminateBar is a segment sliding along a track.
func IndeterminateBar(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("feedback: IndeterminateBar needs a label")
	}
	t := c.Theme()
	rad := t.Space(0.75)
	rtl := c.Direction() == ui.RightToLeft
	return ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(rad).Background(t.Border).
		Role(ui.RoleProgress).Label(label).Range(0, 1, -1).Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		phase := float32(p.Now().UnixMilli()%1400) / 1400
		w := r.W * 0.3
		x := r.X - w + (r.W+w)*phase
		if rtl {
			x = r.X + r.W - (r.W+w)*phase
		}
		p.Clip(r, rad, func() { p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, t.Accent, rad) })
	})
}
