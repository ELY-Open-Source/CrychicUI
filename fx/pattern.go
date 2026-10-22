package fx

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PatternKind is the figure a Pattern repeats.
type PatternKind int

const (
	Stripes PatternKind = iota
	Checks
)

// PatternOptions are a Pattern's optional knobs.
type PatternOptions struct {
	Kind PatternKind // Stripes when zero
}

// Pattern drifts stripes or checks in theme tints behind content.
func Pattern(c *ui.Context, o PatternOptions, content func()) ui.Element {
	if o.Kind != Stripes && o.Kind != Checks {
		panic(fmt.Sprintf("fx: unknown PatternKind %d", o.Kind))
	}
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 6*time.Second)
	col, size := t.Accent.Alpha(0.1), t.Space(4)
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		if o.Kind == Checks {
			off := 2 * size * k
			for j := -2; float32(j)*size < r.H; j++ {
				for i := -2; float32(i)*size < r.W; i++ {
					if (i+j)%2 == 0 {
						p.Fill(ui.Rect{X: r.X + float32(i)*size + off, Y: r.Y + float32(j)*size + off, W: size, H: size}, col, 0)
					}
				}
			}
			return
		}
		for x := r.X - r.H - size + size*k; x < r.X+r.W; x += size {
			var band ui.Path
			band.MoveTo(x, r.Y+r.H).LineTo(x+r.H, r.Y).LineTo(x+r.H+size/2, r.Y).LineTo(x+size/2, r.Y+r.H)
			p.FillPath(band.Close(), col)
		}
	})).Children(content)
}
