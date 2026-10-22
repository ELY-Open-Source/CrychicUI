package chart

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// ZoomChartOptions tune a ZoomChart.
type ZoomChartOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
	// Brush names the range strip for assistive technology.
	Brush string
}

// ZoomChart shows from to to; a brush picks them.
func ZoomChart(c *ui.Context, label string, labels []string, series []Series, from, to *int, o ZoomChartOptions) ui.Element {
	vals, _ := table("ZoomChart", label, labels, series)
	n := len(labels)
	if n < 2 || from == nil || to == nil || *from < 0 || *from >= *to || *to >= n {
		panic(fmt.Sprintf("chart: ZoomChart %q needs from < to within %d labels", label, n))
	}
	def(&o.Brush, "Zoom range")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	at := ui.Local(root, "at", func() int { return -1 })
	set := func(a, b int) {
		if a == *from && b == *to {
			return
		}
		*from, *to, *at = a, b, -1
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("ZoomChart", "zoomed", "from", a, "to", b)
	}
	last := ui.Local(root, "range", func() [2]int { return [2]int{*from, *to} })
	if *last != [2]int{*from, *to} {
		*last, *at = [2]int{*from, *to}, -1
	}
	return root.Children(func() {
		a, b := *from, *to
		sub := make([]Series, len(series))
		for s, sr := range series {
			sub[s] = Series{Name: sr.Name, Values: sr.Values[a : b+1]}
		}
		lines(c, label, labels[a:b+1], sub, at, look{curve: mc.CurveMonotone, h: o.Height, format: o.Format, comp: "ZoomChart"})
		box := ui.Box(c).FillWidth().Height(t.Space(14)).Radius(t.Radius).Border(1, t.Border).Focusable().
			Role(ui.RoleSlider).Label(o.Brush).Value(labels[a] + "–" + labels[b])
		w := box.LayoutBounds().W
		kit.Settle(c, box, []float32{w}, func() []float32 { return []float32{box.LayoutBounds().W} })
		brush := ui.Local(box, "anchor", func() int { return -1 })
		x, _, _ := box.PointerPosition()
		if _, _, held := box.Dragged(); held && w > 0 {
			i := max(0, min(n-1, int(math.Round(float64(x/w*float32(n-1))))))
			if *brush < 0 {
				*brush = i
			}
			lo, hi := min(*brush, i), max(*brush, i)
			if lo == hi {
				hi = min(n-1, lo+1)
				lo = hi - 1
			}
			set(lo, hi)
		} else {
			*brush = -1
		}
		switch {
		case box.Shortcut(0, ui.KeyLeft) && a > 0:
			set(a-1, b-1)
		case box.Shortcut(0, ui.KeyRight) && b < n-1:
			set(a+1, b+1)
		case box.Shortcut(0, ui.KeyUp) && b-a > 1:
			set(a, b-1)
		case box.Shortcut(0, ui.KeyDown):
			if b < n-1 {
				set(a, b+1)
			} else if a > 0 {
				set(a-1, b)
			}
		}
		lo, hi := mc.Extent(vals[0]...)
		for _, vs := range vals[1:] {
			l, h := mc.Extent(vs...)
			lo, hi = min(lo, l), max(hi, h)
		}
		if hi == lo {
			hi = lo + 1
		}
		x0 := box.AnimateTo("from", float32(*from)/float32(n-1), ui.Snappy)
		x1 := box.AnimateTo("to", float32(*to)/float32(n-1), ui.Snappy)
		box.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			xs := make([]float32, n)
			for i := range xs {
				xs[i] = float32(i) / float32(n-1)
			}
			for s, vs := range vals {
				ys := make([]float32, n)
				for i, v := range vs {
					ys[i] = 0.1 + 0.8*float32((v-lo)/(hi-lo))
				}
				p.StrokePath(mc.Line(r, xs, ys, mc.CurveMonotone), 1, tone(t, s).Alpha(0.6))
			}
			p.Fill(ui.Rect{X: mc.X(r, x0), Y: r.Y, W: (x1 - x0) * r.W, H: r.H}, t.Accent.Alpha(0.15), 0)
			for _, e := range []float32{x0, x1} {
				p.Line(mc.X(r, e), r.Y, mc.X(r, e), r.Y+r.H, 2, t.Accent)
			}
		}))
	})
}
