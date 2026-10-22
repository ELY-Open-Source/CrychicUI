package chart

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Node is a named share; children sum up.
type Node struct {
	Name     string
	Value    float64
	Children []Node
}

// SunburstOptions tune a Sunburst.
type SunburstOptions = PieOptions

// arcPart is one node's ring and angles.
type arcPart struct {
	name       string
	depth, top int
	value      float64
	start, end float64
}

// Sunburst rings nested shares; at indexes nodes in order.
func Sunburst(c *ui.Context, label string, roots []Node, at *int, o SunburstOptions) ui.Element {
	var parts []arcPart
	deepest := 0
	var weigh func(n Node) float64
	weigh = func(n Node) float64 {
		if n.Name == "" || math.IsNaN(n.Value) || math.IsInf(n.Value, 0) || n.Value < 0 {
			panic(fmt.Sprintf("chart: Sunburst %q nodes need names and finite values of zero or more", label))
		}
		if len(n.Children) == 0 {
			return n.Value
		}
		sum := 0.0
		for _, k := range n.Children {
			sum += weigh(k)
		}
		return sum
	}
	total := 0.0
	for _, n := range roots {
		total += weigh(n)
	}
	if label == "" || total == 0 {
		panic(fmt.Sprintf("chart: Sunburst %q needs nodes with value", label))
	}
	var lay func(ns []Node, depth, top int, from float64)
	lay = func(ns []Node, depth, top int, from float64) {
		for i, n := range ns {
			if depth == 0 {
				top = i
			}
			v := weigh(n)
			span := 2 * math.Pi * v / total
			parts = append(parts, arcPart{n.Name, depth, top, v, from, from + span})
			deepest = max(deepest, depth)
			lay(n.Children, depth+1, top, from)
			from += span
		}
	}
	lay(roots, 0, 0, 0)
	chosen("Sunburst", label, at, len(parts))
	def(&o.Format, "%g")
	t := c.Theme()
	size := o.Size
	if size == 0 {
		size = t.Space(56)
	}
	hole := size * 0.12
	width := (size/2 - t.Space(1) - hole) / float32(deepest+1)
	ring := func(d int) (float32, float32) { return hole + width*float32(d), hole + width*float32(d+1) - 1 }
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		starts, ends := make([]float64, len(parts)), make([]float64, len(parts))
		for i, pt := range parts {
			starts[i], ends[i] = pt.start, pt.end
		}
		starts, ends = springs(root, "start", starts), springs(root, "end", ends)
		wheel{label: label, comp: "Sunburst", n: len(parts), size: size,
			hit: func(x, y, g float32) int {
				d, a := radius(x, y), turn(x, y)
				for i, pt := range parts {
					if lo, hi := ring(pt.depth); d >= lo && d <= hi && a >= starts[i]*float64(g) && a < ends[i]*float64(g) {
						return i
					}
				}
				return -1
			},
			describe: func(i int) string {
				return fmt.Sprintf("%s: "+o.Format+" (%.0f%%)", parts[i].name, parts[i].value, 100*parts[i].value/total)
			},
			draw: func(p *ui.Painter, r ui.Rect, hot int, g float32) {
				cx, cy := r.X+r.W/2, r.Y+r.H/2
				for i, pt := range parts {
					lo, hi := ring(pt.depth)
					ink := tone(t, pt.top).Mix(t.Background, 0.25*float32(pt.depth))
					if *at >= 0 && *at != i {
						ink = ink.Alpha(0.45)
					}
					if i == hot {
						ink = ink.Mix(t.Text, 0.15)
					}
					s := mc.Slice{Start: starts[i] * float64(g), End: ends[i] * float64(g), Pad: 0.01}
					p.FillPath(mc.Sector(cx, cy, lo, hi, s), ink)
				}
			},
		}.build(c, root, at, nil)
	})
}
