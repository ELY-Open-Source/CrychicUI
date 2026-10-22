package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// TreemapOptions tune a Treemap.
type TreemapOptions = LineOptions

// cell is one Treemap node, flattened in order.
type cell struct {
	name       string
	value      float64
	depth, top int
	kids       []int
}

// Treemap nests rectangles by value; at indexes nodes.
func Treemap(c *ui.Context, label string, roots []Node, at *int, o TreemapOptions) ui.Element {
	if label == "" || len(roots) == 0 {
		panic("chart: Treemap needs a label and nodes")
	}
	var cells []cell
	var add func(ns []Node, depth, top int) []int
	add = func(ns []Node, depth, top int) []int {
		var ids []int
		for i, n := range ns {
			if depth == 0 {
				top = i
			}
			finite("Treemap", label, n.Value)
			if n.Name == "" {
				panic(fmt.Sprintf("chart: Treemap %q nodes need names", label))
			}
			id := len(cells)
			cells = append(cells, cell{name: n.Name, depth: depth, top: top})
			kids := add(n.Children, depth+1, top)
			v := n.Value
			if len(kids) > 0 {
				v = 0
				for _, k := range kids {
					v += cells[k].value
				}
			}
			if v <= 0 {
				panic(fmt.Sprintf("chart: Treemap %q leaves need values above zero", label))
			}
			cells[id].value, cells[id].kids = v, kids
			ids = append(ids, id)
		}
		return ids
	}
	tops := add(roots, 0, 0)
	chosen("Treemap", label, at, len(cells))
	def(&o.Format, "%g")
	t := c.Theme()
	vs := make([]float64, len(cells))
	for i, cl := range cells {
		vs[i] = cl.value
	}
	pad, head := t.Space(1), t.Rem(0.8)*1.4+t.Space(2)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		shown := springs(root, "v", vs)
		plot(c, label, mc.Axes{}, o.Height, func(area ui.Element) {
			g := grow(area)
			b := area.LayoutBounds()
			kit.Settle(c, area, []float32{b.W, b.H}, func() []float32 {
				r := area.LayoutBounds()
				return []float32{r.W, r.H}
			})
			frac := func(px, size float32) float32 {
				if size <= 0 {
					return 0
				}
				return px / size
			}
			boxes := make([][2]mc.Span, len(cells))
			var lay func(ids []int, x, y mc.Span)
			lay = func(ids []int, x, y mc.Span) {
				part := make([]float64, len(ids))
				for j, id := range ids {
					part[j] = shown[id]
				}
				for j, tl := range squarify(part) {
					id := ids[j]
					bx := mc.Span{Lo: x.Lo + tl[0].Lo*(x.Hi-x.Lo), Hi: x.Lo + tl[0].Hi*(x.Hi-x.Lo)}
					by := mc.Span{Lo: y.Lo + tl[1].Lo*(y.Hi-y.Lo), Hi: y.Lo + tl[1].Hi*(y.Hi-y.Lo)}
					boxes[id] = [2]mc.Span{bx, by}
					if kids := cells[id].kids; len(kids) > 0 {
						ix := mc.Span{Lo: bx.Lo + frac(pad, b.W), Hi: bx.Hi - frac(pad, b.W)}
						iy := mc.Span{Lo: by.Lo + frac(pad, b.H), Hi: by.Hi - frac(head, b.H)}
						ix.Hi, iy.Hi = max(ix.Hi, ix.Lo), max(iy.Hi, iy.Lo)
						lay(kids, ix, iy)
					}
				}
			}
			lay(tops, mc.Span{Hi: 1}, mc.Span{Hi: 1})
			for i, cl := range cells {
				c.Key(i)
				x, y := boxes[i][0], boxes[i][1]
				x = mc.Span{Lo: x.Lo * g, Hi: x.Hi * g}
				tip := fmt.Sprintf("%s: "+o.Format, cl.name, cl.value)
				ink, text := tone(t, cl.top).Mix(t.Background, 0.15*float32(cl.depth)), t.AccentText
				if len(cl.kids) > 0 {
					ink, text = tone(t, cl.top).Mix(t.Background, 0.75), t.Text
				}
				m := mark(c, root, at, i, x, y, ink, tip, "Treemap").Radius(t.Radius / 2).Padding(t.Space(1)).Clip()
				m.Children(func() {
					ui.Text(c, cl.name).FontSize(t.Rem(0.8)).Bold().TextColor(text).NoWrap()
				})
			}
		})
	})
}
