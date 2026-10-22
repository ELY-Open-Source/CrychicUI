package motion

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// RingsOptions tune a Rings.
type RingsOptions struct {
	// Count is how many rings travel at once; 3.
	Count int
}

// Rings sends fading rings out from a center.
func Rings(c *ui.Context, o RingsOptions) ui.Element {
	if o.Count < 0 {
		panic("motion: Rings needs a positive Count")
	}
	if o.Count == 0 {
		o.Count = 3
	}
	t := c.Theme()
	d := t.Space(16)
	root := ui.Box(c).Size(d, d).Shrink(0)
	k := float64(root.Loop("crychic.rings", 2400*time.Millisecond, ui.Linear))
	return root.Draw(func(p *ui.Painter, r ui.Rect) {
		for i := range o.Count {
			f := float32(math.Mod(k+float64(i)/float64(o.Count), 1))
			var ring ui.Path
			p.StrokePath(ring.Circle(r.X+r.W/2, r.Y+r.H/2, f*r.W/2), t.Space(0.5), t.Accent.Alpha(1-f))
		}
	})
}
