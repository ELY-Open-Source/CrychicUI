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
	curve        mc.Curve
	fill, stack  bool
	h            float32
	format, comp string
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
	chosen(l.comp, label, at, len(labels))
	def(&l.format, "%g")
	t := c.Theme()
	lo, hi := 0.0, 0.0
	segs := mc.Stack(vals, mc.StackStandard)
	for s, vs := range vals {
		for i, v := range vs {
			if l.stack {
				v = segs[s][i].Hi
			}
			lo, hi = min(lo, v), max(hi, v)
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
			tops := make([][]float32, len(vals))
			bots := make([][]float32, len(vals))
			for s := range vals {
				top, bot := vals[s], make([]float64, len(labels))
				if l.stack {
					top = make([]float64, len(labels))
					for i := range labels {
						top[i], bot[i] = segs[s][i].Hi, segs[s][i].Lo
					}
				}
				top, bot = springs(area, fmt.Sprint("hi", s), top), springs(area, fmt.Sprint("lo", s), bot)
				for i := range labels {
					tops[s] = append(tops[s], y.At(top[i]))
					bots[s] = append(bots[s], y.At(bot[i]))
				}
			}
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for s := len(vals) - 1; s >= 0; s-- {
					ink := tone(t, s)
					if l.fill {
						p.FillPath(mc.Area(r, xs, bots[s], tops[s], l.curve), ink.Alpha(0.2*g))
					}
					p.StrokePath(mc.Line(r, xs, tops[s], l.curve).Trim(0, g), 2, ink)
				}
			}))
			if i := *at; i >= 0 {
				for s := range vals {
					c.Key(s)
					mc.Point(c, xs[i], tops[s][i], t.Space(3), t.Space(3)).Radius(999).Background(t.Background).
						Border(2, tone(t, s)).Transition(spring)
				}
			}
			follow(c, root, area, at, hover, xs, 1, describe, l.comp)
		})
		if len(series) > 1 {
			legend(c, names)
		}
	})
	return root
}
