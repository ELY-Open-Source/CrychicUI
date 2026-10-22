package fx

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Beam runs light between two box-fraction anchors behind content.
func Beam(c *ui.Context, from, to ui.Point, content func()) ui.Element {
	for _, a := range []ui.Point{from, to} {
		if !(a.X >= 0 && a.X <= 1 && a.Y >= 0 && a.Y <= 1) {
			panic(fmt.Sprintf("fx: Beam anchor %v outside 0 to 1", a))
		}
	}
	root := scene(c, content)
	t := c.Theme()
	head := 1.2 * phase(c, root, 2500*time.Millisecond)
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		ax, ay, bx, by := r.X+from.X*r.W, r.Y+from.Y*r.H, r.X+to.X*r.W, r.Y+to.Y*r.H
		lift := float32(math.Hypot(float64(bx-ax), float64(by-ay))) * 0.25
		curve := func() *ui.Path {
			var l ui.Path
			return l.MoveTo(ax, ay).QuadTo((ax+bx)/2, max(r.Y, min(ay, by)-lift), bx, by)
		}
		p.StrokePath(curve(), 1, t.Border)
		p.StrokePath(curve().Trim(head-0.2, head), t.Space(0.5), t.Accent)
		for _, x := range [][2]float32{{ax, ay}, {bx, by}} {
			var dot ui.Path
			dot.Circle(x[0], x[1], t.Space(1.5))
			p.FillPath(&dot, t.Surface)
			p.StrokePath(&dot, 1, t.Border)
		}
	})).Children(content)
}
