package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// TimeLeftBar is a bar draining until ends, total long.
func TimeLeftBar(c *ui.Context, label string, ends time.Time, total time.Duration) ui.Element {
	if label == "" || total <= 0 {
		panic("feedback: TimeLeftBar needs a label and a positive total")
	}
	t := c.Theme()
	rad := t.Space(0.75)
	left := min(max(ends.Sub(c.Now()), 0), total)
	rtl := c.Direction() == ui.RightToLeft
	return ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(rad).Background(t.Border).Role(ui.RoleProgress).Label(label).
		Range(0, total.Seconds(), left.Seconds()).Draw(func(p *ui.Painter, r ui.Rect) {
		rest := min(max(ends.Sub(p.Now()), 0), total)
		if rest > 0 {
			p.AnimationFrame()
		}
		k := float32(rest) / float32(total)
		col := t.Accent
		if k < 0.2 {
			col = t.Accent.Mix(t.Danger, 1-k/0.2)
		}
		w := r.W * k
		x := r.X
		if rtl {
			x = r.X + r.W - w
		}
		p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, col, rad)
	})
}
