package feedback

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// DotsLoader is three dots bouncing in turn.
func DotsLoader(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("feedback: DotsLoader needs a label")
	}
	t := c.Theme()
	d := t.Space(2)
	return ui.Box(c).Size(d*5, d*3).Shrink(0).Role(ui.RoleProgress).Label(label).Range(0, 1, -1).Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		ms := float64(p.Now().UnixMilli() % 1200)
		for i := range 3 {
			phase := math.Mod(ms/1200-float64(i)*0.15+1, 1)
			lift := float32(math.Max(0, math.Sin(phase*2*math.Pi))) * d
			x := r.X + float32(i)*d*2
			p.Fill(ui.Rect{X: x, Y: r.Y + d*2 - lift, W: d, H: d}, t.Accent.Alpha(0.5+0.5*lift/d), d/2)
		}
	})
}
