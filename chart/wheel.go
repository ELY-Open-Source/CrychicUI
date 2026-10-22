package chart

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// wheel is a drawing whose parts hover and pick.
type wheel struct {
	label, comp string
	n           int
	// size is the width; high the height, size when zero.
	size, high float32
	// hit finds the drawn part at x, y, or -1.
	hit      func(x, y, g float32) int
	describe func(i int) string
	draw     func(p *ui.Painter, r ui.Rect, hot int, g float32)
}

// build lays the wheel in root and returns its box.
func (w wheel) build(c *ui.Context, root ui.Element, at *int, center func()) ui.Element {
	if w.high == 0 {
		w.high = w.size
	}
	box := ui.Box(c).Size(w.size, w.high).Shrink(0).Center().Focusable().Role(ui.RoleGroup).Label(w.label)
	hot := ui.Local(box, "hot", func() int { return -1 })
	under := ui.Local(box, "under", func() int { return -1 })
	if *hot >= w.n {
		*hot = -1
	}
	g := grow(box)
	pointed := -1
	x, y, over := box.PointerPosition()
	if over {
		pointed = w.hit(x-w.size/2, y-w.high/2, g)
		// A resting pointer leaves keyboard steps alone.
		if pointed != *under {
			*hot, *under = pointed, pointed
		}
	} else {
		*under = -1
		if !box.Focused() {
			*hot = -1
		}
	}
	switch {
	case box.Shortcut(0, ui.KeyRight), box.Shortcut(0, ui.KeyDown):
		*hot = (*hot + 1) % w.n
	case box.Shortcut(0, ui.KeyLeft), box.Shortcut(0, ui.KeyUp):
		*hot = (*hot - 1 + w.n) % w.n
	case box.Shortcut(0, ui.KeyEscape):
		*hot = -1
	}
	// Keys pick the highlight; a pointer picks what it hits.
	if box.Shortcut(0, ui.KeyEnter) || box.Shortcut(0, ui.KeySpace) {
		if *hot >= 0 {
			pick(root, at, *hot, w.comp)
		}
	} else if box.Clicked() && pointed >= 0 {
		*hot = pointed
		pick(root, at, pointed, w.comp)
	}
	if *hot >= 0 {
		box.Tooltip(w.describe(*hot)).Description(w.describe(*hot))
	}
	shown := *hot
	box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) { w.draw(p, r, shown, g) }))
	if center != nil {
		box.Children(center)
	}
	return box
}

// pies draws shares as slices around a hole.
func pies(c *ui.Context, label string, names []string, vs []float64, at *int, size, hole float32, format, comp string, center func()) ui.Element {
	named(comp, label, names)
	if len(vs) != len(names) {
		panic(fmt.Sprintf("chart: %s %q needs a value per name", comp, label))
	}
	finite(comp, label, vs...)
	sum := 0.0
	for _, v := range vs {
		if v < 0 {
			panic(fmt.Sprintf("chart: %s %q values must not be negative", comp, label))
		}
		sum += v
	}
	if sum == 0 {
		panic(fmt.Sprintf("chart: %s %q needs a value above zero", comp, label))
	}
	chosen(comp, label, at, len(names))
	def(&format, "%g")
	t := c.Theme()
	size = side(comp, label, size, t.Space(56))
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center)
	root.Children(func() {
		slices := mc.Pie(springs(root, "v", vs), 0.02)
		outer := size/2 - t.Space(2)
		inner := outer * hole
		lifts := make([]float32, len(vs))
		for i := range vs {
			lifts[i] = root.AnimateTo(fmt.Sprint("lift", i), unit(i == *at), ui.Bouncy)
		}
		wheel{label: label, comp: comp, n: len(vs), size: size,
			hit: func(x, y, g float32) int {
				// Later slices draw on top, so they hit first.
				for i := len(slices) - 1; i >= 0; i-- {
					s := slices[i]
					s.Start, s.End = s.Start*float64(g), s.End*float64(g)
					dx, dy := mc.Polar(0, 0, t.Space(1.5)*lifts[i], s.Mid())
					if _, ok := mc.SliceAt([]mc.Slice{s}, x-dx, y-dy, inner, outer); ok {
						return i
					}
				}
				return -1
			},
			describe: func(i int) string {
				return fmt.Sprintf("%s: "+format+" (%.0f%%)", names[i], vs[i], 100*vs[i]/sum)
			},
			draw: func(p *ui.Painter, r ui.Rect, hot int, g float32) {
				for i, s := range slices {
					s.Start, s.End = s.Start*float64(g), s.End*float64(g)
					cx, cy := mc.Polar(r.X+r.W/2, r.Y+r.H/2, t.Space(1.5)*lifts[i], s.Mid())
					ink := tone(t, i)
					if i == hot {
						ink = ink.Mix(t.Text, 0.15)
					}
					p.FillPath(mc.Sector(cx, cy, inner, outer, s), ink)
				}
			},
		}.build(c, root, at, center)
		legend(c, names)
	})
	return root
}

// unit is 1 for true.
func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// turn is the angle of x, y clockwise from twelve.
func turn(x, y float32) float64 {
	a := math.Atan2(float64(x), float64(-y))
	if a < 0 {
		a += 2 * math.Pi
	}
	return a
}

// radius is the distance of x, y from the center.
func radius(x, y float32) float32 { return float32(math.Hypot(float64(x), float64(y))) }
