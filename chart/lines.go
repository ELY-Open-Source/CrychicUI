package chart

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// look is how a line chart draws its series.
type look struct {
	curve                       mc.Curve
	smooth, fill, stack, center bool
	h                           float32
	format, comp                string
	// off hides series by index; hidden ones fade out.
	off []bool
	// also are values the axis must hold.
	also []float64
	// over adds marks over the plot.
	over func(area ui.Element, y mc.Linear, xs []float32, aims [][]float32, hover int)
}

// lines draws series over labels as lines or areas.
func lines(c *ui.Context, label string, labels []string, series []Series, at *int, l look) ui.Element {
	vals, names := table(l.comp, label, labels, series)
	for _, vs := range vals {
		for _, v := range vs {
			if l.stack && v < 0 {
				panic(fmt.Sprintf("chart: %s %q stacks values of zero or more", l.comp, label))
			}
		}
	}
	if l.off != nil && len(l.off) != len(series) {
		panic(fmt.Sprintf("chart: %s %q needs a shown flag per series", l.comp, label))
	}
	finite(l.comp, label, l.also...)
	chosen(l.comp, label, at, len(labels))
	def(&l.format, "%g")
	t := c.Theme()
	how := mc.StackStandard
	if l.center {
		how = mc.StackCenter
	}
	segs := mc.Stack(vals, how)
	lo, hi := 0.0, 0.0
	for _, v := range l.also {
		lo, hi = min(lo, v), max(hi, v)
	}
	for s, vs := range vals {
		for i, v := range vs {
			if l.stack {
				lo, hi = min(lo, segs[s][i].Lo), max(hi, segs[s][i].Hi)
			} else {
				lo, hi = min(lo, v), max(hi, v)
			}
		}
	}
	if hi == lo {
		hi = lo + 1
	}
	x := mc.NewBand(labels...)
	y := mc.NewLinear(lo, hi).Nice(4)
	ticks := y.Ticks(4)
	xs := make([]float32, len(labels))
	for i := range labels {
		xs[i] = x.At(i)
	}
	describe := func(i int) string {
		var parts []string
		for s, n := range names {
			parts = append(parts, n+" "+fmt.Sprintf(l.format, vals[s][i]))
		}
		return labels[i] + ": " + strings.Join(parts, ", ")
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	hover := ui.Local(root, "hover", func() int { return -1 })
	root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, l.h, func(area ui.Element) {
			g := grow(area)
			aims := make([][]float32, len(vals))
			tops := make([][]float32, len(vals))
			bots := make([][]float32, len(vals))
			shown := make([]float32, len(vals))
			for s := range vals {
				top, bot := vals[s], make([]float64, len(labels))
				if l.stack {
					top = make([]float64, len(labels))
					for i := range labels {
						top[i], bot[i] = segs[s][i].Hi, segs[s][i].Lo
					}
				}
				var floor []float32
				for i := range labels {
					aims[s] = append(aims[s], y.At(top[i]))
					floor = append(floor, y.At(bot[i]))
				}
				tops[s], bots[s] = eased(area, fmt.Sprint("hi", s), aims[s]), eased(area, fmt.Sprint("lo", s), floor)
				shown[s] = 1
				if l.off != nil {
					shown[s] = area.AnimateTo(fmt.Sprintf("off.%d", s), 1-unit(l.off[s]), ui.Snappy)
				}
			}
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for s := len(vals) - 1; s >= 0; s-- {
					ink := tone(t, s).Alpha(shown[s])
					if l.fill {
						p.FillPath(mc.Area(r, xs, bots[s], tops[s], l.curve), ink.Alpha(0.2*g))
					}
					path := mc.Line(r, xs, tops[s], l.curve)
					if l.smooth {
						path = catmull(r, xs, tops[s])
					}
					p.StrokePath(path.Trim(0, g), 2, ink)
				}
			}))
			if i := *at; i >= 0 {
				for s := range vals {
					if l.off != nil && l.off[s] {
						continue
					}
					c.Key(s)
					mc.Point(c, xs[i], aims[s][i], t.Space(3), t.Space(3)).Radius(999).Background(t.Background).
						Border(2, tone(t, s)).Transition(spring)
				}
			}
			follow(c, root, area, at, hover, xs, 1, describe, l.comp)
			if l.over != nil {
				l.over(area, y, xs, aims, *hover)
			}
		})
		if len(series) > 1 && l.off == nil {
			legend(c, names)
		}
	})
	return root
}

// catmull is a Catmull-Rom curve through the points.
func catmull(r ui.Rect, xs, ys []float32) *ui.Path {
	var path ui.Path
	at := func(i int) (float32, float32) {
		i = max(0, min(len(xs)-1, i))
		return mc.X(r, xs[i]), mc.Y(r, ys[i])
	}
	for i := range xs {
		x1, y1 := at(i)
		if i == 0 {
			path.MoveTo(x1, y1)
			continue
		}
		x0, y0 := at(i - 1)
		xp, yp := at(i - 2)
		x2, y2 := at(i + 1)
		path.CubeTo(x0+(x1-xp)/6, y0+(y1-yp)/6, x1-(x2-x0)/6, y1-(y2-y0)/6, x1, y1)
	}
	return &path
}
