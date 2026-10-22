package feedback

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// BarsLoader is equalizer bars rising and falling.
func BarsLoader(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("feedback: BarsLoader needs a label")
	}
	t := c.Theme()
	w := t.Space(1)
	return ui.Box(c).Size(w*9, t.Space(6)).Shrink(0).Role(ui.RoleProgress).Label(label).Range(0, 1, -1).Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		s := float64(p.Now().UnixMilli()) / 1000
		for i := range 5 {
			k := float32(0.3 + 0.7*(0.5+0.5*math.Sin(s*2*math.Pi*(0.9+0.23*float64(i))+float64(i))))
			h := r.H * k
			p.Fill(ui.Rect{X: r.X + float32(i)*w*2, Y: r.Y + r.H - h, W: w, H: h}, t.Accent, w/2)
		}
	})
}
