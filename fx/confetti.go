package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Confetti bursts falling, spinning pieces over content when run changes.
func Confetti(c *ui.Context, run uint64, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	played := ui.Local(root, "played", func() uint64 { return run })
	if run != *played {
		*played = run
		sound.Play(sound.Pop)
		kit.Log("Confetti", "burst", "run", run)
	}
	k := burst(root, run, 1600*time.Millisecond)
	cols := palette(t)
	w, h := t.Space(1.5), t.Space(0.75)
	const n = 40
	return root.Children(content).DrawOver(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		if k <= 0 || k >= 1 {
			return
		}
		for i := range n {
			vx, vy := (rnd(i, 1)-0.5)*r.W*1.2, -(0.3+0.5*rnd(i, 2))*r.H*2
			x := r.X + r.W/2 + vx*k
			y := r.Y + r.H*0.4 + vy*k + 2*r.H*k*k
			turn := 720 * k * (rnd(i, 3) - 0.5) * 2
			col := cols[i%len(cols)].Alpha(1 - k*k*k)
			p.Transform(ui.Rotation(turn), x, y, func() { p.Fill(ui.Rect{X: x - w/2, Y: y - h/2, W: w, H: h}, col, 0) })
		}
	}))
}
