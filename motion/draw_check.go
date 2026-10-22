package motion

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var redraw = replay(450*time.Millisecond, ui.EaseOut)

// DrawCheck draws a check stroke each time run changes.
func DrawCheck(c *ui.Context, label string, run uint64) ui.Element {
	if label == "" {
		panic("motion: DrawCheck needs a label")
	}
	t := c.Theme()
	d, w := t.Space(12), t.Space(1)
	root := ui.Box(c).Size(d, d).Shrink(0).Role(ui.RoleImage).Label(label)
	played := ui.Local(root, "played", func() uint64 { return run })
	if run != *played {
		*played = run
		sound.Play(sound.Success)
		kit.Log("DrawCheck", "drawn", "run", run)
	}
	k := onChange(root, run, redraw)
	return root.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, (r.W-w)/2
		var ring, mark ui.Path
		p.StrokePath(ring.Circle(cx, cy, rad), w, t.Success.Mix(t.Border, 1-k))
		if k > 0 {
			mark.MoveTo(r.X+0.28*r.W, r.Y+0.52*r.H).LineTo(r.X+0.44*r.W, r.Y+0.68*r.H).LineTo(r.X+0.72*r.W, r.Y+0.36*r.H)
			p.StrokePath(mark.Trim(0, k), w, t.Success)
		}
	})
}
