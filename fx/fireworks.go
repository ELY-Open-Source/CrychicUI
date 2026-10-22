package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Fireworks rises shells that open into sparks when run changes.
func Fireworks(c *ui.Context, run uint64, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	played := ui.Local(root, "played", func() uint64 { return run })
	if run != *played {
		*played = run
		sound.Play(sound.Pop)
		kit.Log("Fireworks", "launched", "run", run)
	}
	k := burst(root, run, 1800*time.Millisecond)
	cols := palette(t)
	dot, reach := t.Space(0.75), t.Space(14)
	const shells, sparks, climb = 3, 16, 0.35
	return root.Children(content).DrawOver(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		if k <= 0 || k >= 1 {
			return
		}
		for i := range shells {
			u := min(max((k-float32(i)*0.15)/0.7, 0), 1)
			if u == 0 || u == 1 {
				continue
			}
			x, top := r.X+r.W*(0.25+0.25*float32(i)), r.Y+r.H*(0.25+0.1*rnd(i, 1))
			col := cols[i%len(cols)]
			if u < climb {
				f := ui.EaseOut(u / climb)
				var shell ui.Path
				p.FillPath(shell.Circle(x, r.Y+r.H-(r.Y+r.H-top)*f, dot), col)
				continue
			}
			e := ui.EaseOut((u - climb) / (1 - climb))
			for j := range sparks {
				a := 2 * math.Pi * float64(j) / sparks
				sx := x + reach*e*float32(math.Cos(a))
				sy := top + reach*e*float32(math.Sin(a)) + reach*0.5*e*e
				var spark ui.Path
				p.FillPath(spark.Circle(sx, sy, dot*(1-e*0.5)), col.Alpha(1-e))
			}
		}
	}))
}
