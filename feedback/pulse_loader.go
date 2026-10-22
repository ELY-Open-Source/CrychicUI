package feedback

import "github.com/egoist/mygo/ui"

// PulseLoader is a circle sending rings outward.
func PulseLoader(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("feedback: PulseLoader needs a label")
	}
	t := c.Theme()
	size := t.Space(10)
	return ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleProgress).Label(label).Range(0, 1, -1).Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		core := r.W * 0.18
		for i := range 2 {
			k := float32((p.Now().UnixMilli()+int64(i)*700)%1400) / 1400
			rad := core + (r.W/2-core)*k
			p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: rad * 2, H: rad * 2}, t.Accent.Alpha(0.35*(1-k)), rad)
		}
		p.Fill(ui.Rect{X: cx - core, Y: cy - core, W: core * 2, H: core * 2}, t.Accent, core)
	})
}
