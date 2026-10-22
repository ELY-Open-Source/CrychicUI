package chart

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// TreeDiagramOptions tune a TreeDiagram.
type TreeDiagramOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
}

// branch is one TreeDiagram node, flattened in order.
type branch struct {
	name               string
	depth, top, parent int
	x                  float32
}

// TreeDiagram links nodes as boxes; at indexes nodes.
func TreeDiagram(c *ui.Context, label string, roots []Node, at *int, o TreeDiagramOptions) ui.Element {
	if label == "" || len(roots) == 0 {
		panic("chart: TreeDiagram needs a label and nodes")
	}
	var nodes []branch
	leaves, deepest := 0, 0
	var lay func(ns []Node, depth, top, parent int) []float32
	lay = func(ns []Node, depth, top, parent int) []float32 {
		var xs []float32
		for i, n := range ns {
			if n.Name == "" {
				panic(fmt.Sprintf("chart: TreeDiagram %q nodes need names", label))
			}
			if depth == 0 {
				top = i
			}
			id := len(nodes)
			nodes = append(nodes, branch{name: n.Name, depth: depth, top: top, parent: parent})
			deepest = max(deepest, depth)
			if kids := lay(n.Children, depth+1, top, id); len(kids) > 0 {
				nodes[id].x = (kids[0] + kids[len(kids)-1]) / 2
			} else {
				nodes[id].x = float32(leaves) + 0.5
				leaves++
			}
			xs = append(xs, nodes[id].x)
		}
		return xs
	}
	lay(roots, 0, 0, -1)
	chosen("TreeDiagram", label, at, len(nodes))
	t := c.Theme()
	y := func(depth int) float32 { return 1 - (float32(depth)+0.5)/float32(deepest+1) }
	bw, bh := t.Space(22), t.Space(7)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{}, o.Height, func(area ui.Element) {
			g := grow(area)
			fx, fy := make([]float32, len(nodes)), make([]float32, len(nodes))
			for i, b := range nodes {
				fx[i], fy[i] = b.x/float32(leaves), y(b.depth)
			}
			ex, ey := eased(area, "x", fx), eased(area, "y", fy)
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for i, b := range nodes {
					if b.parent < 0 {
						continue
					}
					x0, y0 := mc.X(r, ex[b.parent]), mc.Y(r, ey[b.parent])+bh/2
					x1, y1 := mc.X(r, ex[i]), mc.Y(r, ey[i])-bh/2
					var link ui.Path
					link.MoveTo(x0, y0).CubeTo(x0, (y0+y1)/2, x1, (y0+y1)/2, x1, y1)
					p.StrokePath(link.Trim(0, g), 1.5, t.Border)
				}
			}))
			for i, b := range nodes {
				c.Key(i)
				box := mc.Point(c, fx[i], fy[i], bw, bh).Radius(t.Radius).Center().Clip().
					Padding(0, t.Space(1.5)).Focusable().Label(b.name).Tooltip(b.name).Transition(spring)
				ink := tone(t, b.top)
				fill := ink.Mix(t.Background, 0.8)
				if *at >= 0 && *at != i {
					fill = fill.Alpha(0.5)
				}
				if box.Hovered() || *at == i {
					fill = ink.Mix(t.Background, 0.55)
				}
				box.Background(box.AnimateColor("fill", fill, ui.Snappy)).Border(1, ink)
				kit.Arrive(box, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(b.depth)*90*time.Millisecond)
				box.Children(func() {
					ui.Text(c, b.name).NoWrap().FontSize(t.Rem(0.75))
				})
				if box.Clicked() {
					pick(root, at, i, "TreeDiagram")
				}
			}
		})
	})
}
