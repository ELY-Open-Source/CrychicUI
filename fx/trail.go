package fx

import (
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// mark is where the pointer was, and when.
type mark struct {
	x, y float32
	at   time.Time
}

// Trail leaves dots fading behind the pointer over content.
func Trail(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	const life, most = 600 * time.Millisecond, 32
	marks := ui.Local(root, "crychic.marks", func() []mark { return nil })
	now := c.Now()
	last := ui.Local(root, "crychic.last", func() ui.Point { return ui.Point{X: -1, Y: -1} })
	x, y, over := root.PointerPosition()
	if at := (ui.Point{X: x, Y: y}); over && at != *last {
		*last = at
		*marks = append(*marks, mark{x, y, now})
	}
	*marks = slices.DeleteFunc(*marks, func(m mark) bool { return now.Sub(m.at) >= life })
	if n := len(*marks); n > most {
		*marks = (*marks)[n-most:]
	}
	if c.Preferences().ReduceMotion {
		*marks = nil
	}
	if len(*marks) > 0 {
		c.AnimationFrame()
	}
	shown := slices.Clone(*marks)
	return root.Children(content).DrawOver(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for _, m := range shown {
			f := 1 - float32(now.Sub(m.at))/float32(life)
			var dot ui.Path
			p.FillPath(dot.Circle(r.X+m.x, r.Y+m.y, t.Space(1.5)*f), t.Accent.Alpha(0.5*f))
		}
	}))
}
