package fx

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ClickBurst bursts particles from each click over content.
func ClickBurst(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	const life = 600 * time.Millisecond
	list, now := taps(c, root, "ClickBurst", sound.Pop, life), c.Now()
	cols := palette(t)
	const n = 12
	return root.Children(content).DrawOver(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for _, tp := range list {
			k := float32(now.Sub(tp.at)) / float32(life)
			x, y := r.X+tp.x, r.Y+tp.y
			out := t.Space(10) * (1 - (1-k)*(1-k)*(1-k))
			for j := range n {
				a := 2 * math.Pi * (float64(j)/n + 0.04*float64(rnd(j, 1)))
				var dot ui.Path
				p.FillPath(dot.Circle(x+out*float32(math.Cos(a)), y+out*float32(math.Sin(a)), t.Space(0.75)*(1-k)),
					cols[j%len(cols)].Alpha(1-k))
			}
		}
	}))
}
