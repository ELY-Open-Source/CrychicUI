package devtools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Service is one map node; Icon is optional.
type Service struct {
	Name, Icon string
	State      HealthState
}

// Call is traffic from one service to another.
type Call struct {
	From, To, Rate string
	Failing        bool
}

// ServiceMapOptions word the panel; States names Up to Down.
type ServiceMapOptions struct {
	Title, Hint string
	States      [3]string
}

// ServiceMap draws services by call depth; picks one; cue Select.
func ServiceMap(c *ui.Context, services []Service, calls []Call, selected *string, o ServiceMapOptions) ui.Element {
	def(&o.Title, "Service map")
	if len(services) == 0 || selected == nil {
		panic(fmt.Sprintf("devtools: ServiceMap %q needs services and a selection", o.Title))
	}
	def(&o.Hint, "Pick a service to see its calls")
	for i, name := range [...]string{"Healthy", "Degraded", "Down"} {
		def(&o.States[i], name)
	}
	names := make([]string, len(services))
	index := map[string]int{}
	for i, s := range services {
		validHealth("ServiceMap", s.Name, s.State)
		names[i] = s.Name
		index[s.Name] = i
	}
	unique("ServiceMap", names)
	pairs := make([]string, len(calls))
	for i, k := range calls {
		_, a := index[k.From]
		_, b := index[k.To]
		if !a || !b || k.From == k.To {
			panic(fmt.Sprintf("devtools: ServiceMap call %q to %q needs two known services", k.From, k.To))
		}
		pairs[i] = strconv.Quote(k.From) + ">" + strconv.Quote(k.To)
	}
	if len(calls) > 0 {
		unique("ServiceMap calls", pairs)
	}
	layer := layers(o.Title, len(services), calls, index)
	cols, pos, rows := 0, make([]int, len(services)), map[int]int{}
	for i := range services {
		cols = max(cols, layer[i]+1)
		pos[i] = rows[layer[i]]
		rows[layer[i]]++
	}
	most := 0
	for _, n := range rows {
		most = max(most, n)
	}
	t := c.Theme()
	nodeH, rowGap, colGap, pad := t.Space(10), t.Space(3), t.Space(6), t.Space(3)
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "radar", o.Title, func() { kit.Chip(c, fmt.Sprint(len(services)), t.Accent) })
		view := ui.ScrollHorizontal(c).FillWidth()
		seen := view.LayoutBounds().W
		kit.Settle(c, view, []float32{seen}, func() []float32 { return []float32{view.LayoutBounds().W} })
		view.Children(func() {
			w := max(seen, 2*pad+float32(cols)*t.Space(18)+float32(cols-1)*colGap)
			nodeW := min(t.Space(32), (w-2*pad-float32(cols-1)*colGap)/float32(cols))
			at := func(i int) (float32, float32) {
				x := (w - nodeW) / 2
				if cols > 1 {
					x = pad + float32(layer[i])*(w-2*pad-nodeW)/float32(cols-1)
				}
				return x, pad + float32(pos[i])*(nodeH+rowGap)
			}
			canvas := ui.Box(c).Width(w).Height(2*pad + float32(most)*nodeH + float32(most-1)*rowGap)
			draw := canvas.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			canvas.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for _, k := range calls {
					x0, y0 := at(index[k.From])
					x1, y1 := at(index[k.To])
					x0, y0, x1, y1 = r.X+x0+nodeW, r.Y+y0+nodeH/2, r.X+x1, r.Y+y1+nodeH/2
					ink, width := t.TextMuted.Alpha(0.35), t.Space(0.25)
					if k.From == *selected || k.To == *selected {
						ink, width = t.Accent, t.Space(0.5)
					}
					if k.Failing {
						ink = t.Danger
					}
					var edge ui.Path
					mid := (x0 + x1) / 2
					edge.MoveTo(x0, y0).CubeTo(mid, y0, mid, y1, x1, y1)
					p.StrokePath(edge.Trim(0, draw), width, ink)
					if draw >= 1 {
						p.FillPath(ring(x1, y1, t.Space(0.75)), ink)
					}
				}
			}))
			canvas.Children(func() {
				for i, s := range services {
					x, y := at(i)
					serviceNode(c, root, s, x, y, nodeW, nodeH, i, selected)
				}
			})
		})
		serviceCalls(c, services, calls, *selected, o)
	})
	return root
}

// serviceNode is one pickable service box.
func serviceNode(c *ui.Context, root ui.Element, s Service, x, y, w, h float32, i int, selected *string) {
	t := c.Theme()
	on := *selected == s.Name
	ink := healthInk(t, s.State)
	c.Key(s.Name)
	b := kit.Press(ui.ButtonBase(c).Label(s.Name).Absolute().Left(x).Top(y).Size(w, h).Padding(0, t.Space(2)).
		Radius(t.Radius).Checked(on).AlignItems(ui.Stretch))
	b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
	b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
	kit.Arrive(b, ui.Pose{Scale: 0.9, Fade: 1}, time.Duration(i)*50*time.Millisecond).Children(func() {
		ui.Row(c).FillWidth().FillHeight().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(ink)
			if s.Icon != "" {
				glyph(c, s.Icon, 14).TextColor(t.TextMuted)
			}
			ui.Text(c, s.Name).Font("monospace").FontSize(t.Rem(0.7)).Bold().Grow(1).MinWidth(0).SingleLine()
		})
	})
	if b.Clicked() && !on {
		*selected = s.Name
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("ServiceMap", "selected", "index", i)
	}
}

// serviceCalls lists calls into and out of the pick.
func serviceCalls(c *ui.Context, services []Service, calls []Call, name string, o ServiceMapOptions) {
	t := c.Theme()
	ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2.5), t.Space(3)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).
		Children(func() {
			picked := -1
			for i, s := range services {
				if s.Name == name {
					picked = i
				}
			}
			if picked < 0 {
				kit.Muted(c, o.Hint).FontSize(t.Rem(0.75))
				return
			}
			s := services[picked]
			c.Key(strconv.Quote(s.Name))
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, 0, t.Space(1), 0).Children(func() {
				ui.Text(c, s.Name).Font("monospace").FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
				kit.Chip(c, o.States[s.State], healthInk(t, s.State))
			})
			for _, k := range calls {
				if k.From != name && k.To != name {
					continue
				}
				other, mark := k.To, "arrow-right"
				if k.To == name {
					other, mark = k.From, "arrow-left"
				}
				c.Key(strconv.Quote(k.From) + ">" + strconv.Quote(k.To))
				ink := t.TextMuted
				if k.Failing {
					ink = t.Danger
				}
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					glyph(c, mark, 14).TextColor(ink)
					ui.Text(c, other).Font("monospace").FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
					if k.Rate != "" {
						kit.Numeric(c, k.Rate).Font("monospace").FontSize(t.Rem(0.7)).TextColor(ink).Shrink(0)
					}
				})
			}
		})
}

// layers is each service's depth along calls; loops panic.
func layers(title string, n int, calls []Call, index map[string]int) []int {
	out := make([][]int, n)
	in := make([]int, n)
	for _, k := range calls {
		a, b := index[k.From], index[k.To]
		out[a] = append(out[a], b)
		in[b]++
	}
	depth := make([]int, n)
	var queue []int
	for i := range n {
		if in[i] == 0 {
			queue = append(queue, i)
		}
	}
	done := 0
	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]
		done++
		for _, j := range out[i] {
			depth[j] = max(depth[j], depth[i]+1)
			if in[j]--; in[j] == 0 {
				queue = append(queue, j)
			}
		}
	}
	if done < n {
		panic(fmt.Sprintf("devtools: ServiceMap %q calls form a loop", title))
	}
	return depth
}
