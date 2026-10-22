package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// OrbitOptions tune an Orbit.
type OrbitOptions struct {
	// Radius is the path's radius; t.Space(12) when zero.
	Radius float32
}

// Orbit circles n items around a center.
func Orbit(c *ui.Context, n int, item func(i int), o OrbitOptions) ui.Element {
	if n < 1 || item == nil {
		panic("motion: Orbit needs n ≥ 1 and an item")
	}
	t := c.Theme()
	r, slot := positive("Orbit", "Radius", o.Radius, t.Space(12)), t.Space(8)
	root := ui.Box(c).Size(2*r+slot, 2*r+slot).Shrink(0)
	k := float64(root.Loop("crychic.orbit", 8*time.Second, ui.Linear))
	root.Draw(func(p *ui.Painter, b ui.Rect) {
		var ring ui.Path
		p.StrokePath(ring.Circle(b.X+b.W/2, b.Y+b.H/2, r), 1, t.Border)
	})
	return root.Children(func() {
		for i := range n {
			c.Key(i)
			a := 2 * math.Pi * (k + float64(i)/float64(n))
			ui.Box(c).Absolute().Top(r).Left(r).Size(slot, slot).Center().
				Translate(r*float32(math.Cos(a)), r*float32(math.Sin(a))).Children(func() { item(i) })
		}
	})
}
