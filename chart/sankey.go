package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Flow is a value moving between two nodes.
type Flow struct {
	From, To int
	Value    float64
}

// SankeyOptions tune a Sankey.
type SankeyOptions = LineOptions

// Sankey links nodes by flows as wide as their value.
func Sankey(c *ui.Context, label string, nodes []string, flows []Flow, at *int, o SankeyOptions) ui.Element {
	named("Sankey", label, nodes)
	n := len(nodes)
	if len(flows) == 0 {
		panic(fmt.Sprintf("chart: Sankey %q needs flows", label))
	}
	in, out := make([]float64, n), make([]float64, n)
	next := make([][]int, n)
	deg := make([]int, n)
	for _, f := range flows {
		finite("Sankey", label, f.Value)
		if f.From < 0 || f.From >= n || f.To < 0 || f.To >= n || f.From == f.To || f.Value <= 0 {
			panic(fmt.Sprintf("chart: Sankey %q flows need two known nodes and a value above zero", label))
		}
		out[f.From] += f.Value
		in[f.To] += f.Value
		next[f.From] = append(next[f.From], f.To)
		deg[f.To]++
	}
	chosen("Sankey", label, at, n)
	def(&o.Format, "%g")
	layer := make([]int, n)
	var queue []int
	for i := range n {
		if deg[i] == 0 {
			queue = append(queue, i)
		}
	}
	for done := 0; done < n; done++ {
		if len(queue) == 0 {
			panic(fmt.Sprintf("chart: Sankey %q flows must not loop", label))
		}
		i := queue[0]
		queue = queue[1:]
		for _, j := range next[i] {
			layer[j] = max(layer[j], layer[i]+1)
			if deg[j]--; deg[j] == 0 {
				queue = append(queue, j)
			}
		}
	}
	cols := 0
	for _, l := range layer {
		cols = max(cols, l+1)
	}
	const gap, nw = 0.05, 0.02
	total := make([]float64, cols)
	count := make([]int, cols)
	for i := range n {
		total[layer[i]] += max(in[i], out[i])
		count[layer[i]]++
	}
	k := 1.0
	for l := range cols {
		k = min(k, (1-gap*float64(count[l]-1))/max(total[l], 1e-9))
	}
	xs, ys := make([]mc.Span, n), make([]mc.Span, n)
	tops := make([]float64, cols)
	for l := range tops {
		tops[l] = 1
	}
	for i := range n {
		l, size := layer[i], max(in[i], out[i])*k
		x0 := float32(float64(l) / float64(max(cols-1, 1)) * (1 - nw))
		xs[i] = mc.Span{Lo: x0, Hi: x0 + nw}
		ys[i] = mc.Span{Lo: float32(tops[l] - size), Hi: float32(tops[l])}
		tops[l] -= size + gap
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{}, o.Height, func(area ui.Element) {
			g := grow(area)
			his, ws := make([]float64, n), make([]float64, len(flows))
			for i := range n {
				his[i] = float64(ys[i].Hi)
			}
			for j, f := range flows {
				ws[j] = f.Value * k
			}
			his, ws = springs(area, "hi", his), springs(area, "w", ws)
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				sOff, tOff := make([]float64, n), make([]float64, n)
				for j, f := range flows {
					w := ws[j]
					s0, t0 := his[f.From]-sOff[f.From], his[f.To]-tOff[f.To]
					sOff[f.From] += w
					tOff[f.To] += w
					x1, x2 := mc.X(r, xs[f.From].Hi), mc.X(r, xs[f.To].Lo)
					mid := (x1 + x2) / 2
					a, b := mc.Y(r, float32(s0)), mc.Y(r, float32(t0))
					a2, b2 := mc.Y(r, float32(s0-w)), mc.Y(r, float32(t0-w))
					var ribbon ui.Path
					ribbon.MoveTo(x1, a).CubeTo(mid, a, mid, b, x2, b).LineTo(x2, b2).CubeTo(mid, b2, mid, a2, x1, a2).Close()
					alpha := float32(0.28)
					if *at == f.From || *at == f.To {
						alpha = 0.55
					}
					p.FillPath(&ribbon, tone(t, f.From).Alpha(alpha*g))
				}
			}))
			for i, name := range nodes {
				c.Key(name)
				tip := fmt.Sprintf("%s: "+o.Format, name, max(in[i], out[i]))
				m, half := ys[i].Mid(), (ys[i].Hi-ys[i].Lo)/2*g
				mark(c, root, at, i, xs[i], mc.Span{Lo: m - half, Hi: m + half}, tone(t, i), tip, "Sankey")
				h, x := ui.Start, xs[i].Hi+0.01
				if layer[i] == cols-1 {
					h, x = ui.End, xs[i].Lo-0.01
				}
				mc.Anchor(c, x, ys[i].Mid(), h, ui.Center, func() {
					ui.Text(c, name).NoWrap().FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
				})
			}
		})
	})
}
