package chart

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Edge links two nodes by index.
type Edge struct{ From, To int }

// NetworkGraphOptions tune a NetworkGraph.
type NetworkGraphOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
}

// NetworkGraph rings nodes joined by edges; nodes drag.
func NetworkGraph(c *ui.Context, label string, nodes []string, edges []Edge, at *int, o NetworkGraphOptions) ui.Element {
	named("NetworkGraph", label, nodes)
	n := len(nodes)
	for _, e := range edges {
		if e.From < 0 || e.From >= n || e.To < 0 || e.To >= n || e.From == e.To {
			panic(fmt.Sprintf("chart: NetworkGraph %q edges need two known nodes", label))
		}
	}
	chosen("NetworkGraph", label, at, n)
	t := c.Theme()
	ring := func() [][2]float32 {
		out := make([][2]float32, n)
		for i := range out {
			x, y := mc.Polar(0.5, 0.5, 0.38*unit(n > 1), 2*math.Pi*float64(i)/float64(n))
			out[i] = [2]float32{x, y}
		}
		return out
	}
	root := ui.Column(c).FillWidth()
	pos := ui.Local(root, "pos", ring)
	if len(*pos) != n {
		*pos = ring()
	}
	dropped := func(i int) {
		sound.Play(sound.Drop)
		kit.Log("NetworkGraph", "moved", "index", i)
		root.SetChanged()
	}
	return root.Children(func() {
		plot(c, label, mc.Axes{}, o.Height, func(area ui.Element) {
			g := grow(area)
			b := area.LayoutBounds()
			kit.Settle(c, area, []float32{b.W, b.H}, func() []float32 {
				r := area.LayoutBounds()
				return []float32{r.W, r.H}
			})
			ps := *pos
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for _, e := range edges {
					ink := t.Border
					if *at == e.From || *at == e.To {
						ink = t.Accent
					}
					var line ui.Path
					line.MoveTo(mc.X(r, ps[e.From][0]), mc.Y(r, ps[e.From][1])).LineTo(mc.X(r, ps[e.To][0]), mc.Y(r, ps[e.To][1]))
					p.StrokePath(line.Trim(0, g), 1.5, ink)
				}
			}))
			for i, name := range nodes {
				c.Key(name)
				d := t.Space(4)
				node := mc.Point(c, ps[i][0], ps[i][1], d*g, d*g).Radius(999).Focusable().Label(name).Tooltip(name)
				move := func(fx, fy float32) {
					ps[i] = [2]float32{max(0.04, min(0.96, ps[i][0]+fx)), max(0.04, min(0.96, ps[i][1]+fy))}
				}
				moved := ui.Local(node, "moved", func() bool { return false })
				dx, dy, held := node.Dragged()
				if held && (dx != 0 || dy != 0) && b.W > 0 && b.H > 0 {
					move(dx/b.W, -dy/b.H)
					*moved = true
				}
				for _, k := range []struct {
					key    ui.Key
					fx, fy float32
				}{{ui.KeyLeft, -0.03, 0}, {ui.KeyRight, 0.03, 0}, {ui.KeyUp, 0, 0.03}, {ui.KeyDown, 0, -0.03}} {
					if node.Shortcut(0, k.key) {
						move(k.fx, k.fy)
						dropped(i)
					}
				}
				if node.Clicked() && !*moved {
					pick(root, at, i, "NetworkGraph")
				}
				if !held && *moved {
					*moved = false
					dropped(i)
				}
				ink := tone(t, i)
				if *at >= 0 && *at != i {
					ink = ink.Alpha(0.4)
				}
				if node.Hovered() || *at == i {
					node.Border(2, t.Text.Alpha(0.6))
				}
				node.Background(node.AnimateColor("ink", ink, ui.Snappy))
				c.Key(name + ".label")
				mc.Anchor(c, ps[i][0], ps[i][1], ui.Center, ui.Start, func() {
					kit.Muted(c, name).NoWrap().FontSize(t.Rem(0.75)).Margin(d/2+t.Space(1), 0, 0, 0).Opacity(g)
				})
			}
		})
	})
}
