package project

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BurndownPoint is one remaining and ideal pair.
type BurndownPoint struct {
	Day       string
	Remaining float64
	Ideal     float64
}

// BurndownOptions word the chart; Left takes the remaining points.
type BurndownOptions struct{ Title, Left, Actual, Ideal string }

// Burndown shows remaining against the ideal line.
func Burndown(c *ui.Context, points []BurndownPoint, o BurndownOptions) ui.Element {
	if len(points) < 2 {
		panic("project: Burndown needs points")
	}
	def(&o.Title, "Burndown")
	def(&o.Left, "%s pts left")
	def(&o.Actual, "Remaining")
	def(&o.Ideal, "Ideal")
	maxValue := 0.0
	for _, point := range points {
		if point.Day == "" || point.Remaining < 0 || point.Ideal < 0 || !finite(point.Remaining) || !finite(point.Ideal) {
			panic(fmt.Sprintf("project: Burndown %q has an invalid point %q", o.Title, point.Day))
		}
		maxValue = max(maxValue, point.Remaining, point.Ideal)
	}
	if maxValue == 0 {
		maxValue = 1
	}
	t := c.Theme()
	last := points[len(points)-1]
	behind := last.Remaining > last.Ideal
	ink := t.Accent
	if behind {
		ink = t.Warning
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "trending-down", o.Title, func() {
			kit.Chip(c, fmt.Sprintf(o.Left, strconv.FormatFloat(last.Remaining, 'f', -1, 64)), ink)
		})
		ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(3)).Children(func() {
			chart := ui.Box(c).FillWidth().Height(t.Space(36))
			draw := chart.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				for i := range 5 {
					y := r.Y + r.H*float32(i)/4
					p.Line(r.X, y, r.X+r.W, y, 1, t.Border.Alpha(0.6))
				}
				at := func(i int, v float64) (float32, float32) {
					return r.X + r.W*(float32(i)+0.5)/float32(len(points)), r.Y + r.H*(1-float32(v/maxValue))
				}
				var actual, ideal, area ui.Path
				for i, point := range points {
					xa, ya := at(i, point.Remaining)
					_, yi := at(i, point.Ideal)
					if i == 0 {
						actual.MoveTo(xa, ya)
						ideal.MoveTo(xa, yi)
						area.MoveTo(xa, r.Y+r.H).LineTo(xa, ya)
					} else {
						actual.LineTo(xa, ya)
						ideal.LineTo(xa, yi)
						area.LineTo(xa, ya)
					}
				}
				end, _ := at(len(points)-1, 0)
				area.LineTo(end, r.Y+r.H).Close()
				p.StrokePath(&ideal, 1.5, t.TextMuted.Alpha(0.45))
				p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W * draw, H: r.H}, 0, func() { p.FillPath(&area, ink.Alpha(0.12)) })
				p.StrokePath(actual.Trim(0, draw), 2, ink)
				for i, point := range points {
					x, y := at(i, point.Remaining)
					if x > r.X+r.W*draw+1 {
						break
					}
					ring, dot := t.Space(0.75), t.Space(0.5)
					p.Fill(ui.Rect{X: x - ring, Y: y - ring, W: 2 * ring, H: 2 * ring}, t.Background, ring)
					p.Fill(ui.Rect{X: x - dot, Y: y - dot, W: 2 * dot, H: 2 * dot}, ink, dot)
				}
			}))
			// Labels own spans around their points and truncate inside.
			stride, last := (len(points)+5)/6, len(points)-1
			var shown []int
			for i := range points {
				if i == 0 || i == last || (i%stride == 0 && i <= last-stride) {
					shown = append(shown, i)
				}
			}
			ui.Row(c).FillWidth().Children(func() {
				at := 0.0
				for k, i := range shown {
					mid := float64(i) + 0.5
					from, to, justify := 0.0, float64(len(points)), ui.Center
					switch {
					case len(shown) == 1:
					case k == 0:
						to, justify = (mid+float64(shown[1])+0.5)/2, ui.Start
					case k == len(shown)-1:
						from, justify = (mid+float64(shown[k-1])+0.5)/2, ui.End
					default:
						half := float64(min(i-shown[k-1], shown[k+1]-i)) / 2
						from, to = mid-half, mid+half
					}
					if from > at {
						ui.Box(c).Grow(float32(from - at)).Basis(0)
					}
					ui.Row(c).Grow(float32(to - from)).Basis(0).MinWidth(0).Justify(justify).Children(func() {
						kit.Muted(c, points[i].Day).FontSize(t.Rem(0.6)).MinWidth(0).SingleLine()
					})
					at = to
				}
			})
			ui.Row(c).Gap(t.Space(3)).Children(func() {
				legend(c, o.Actual, ink)
				legend(c, o.Ideal, t.TextMuted.Alpha(0.45))
			})
		})
	})
	return root
}

// legend is a colored swatch and its name.
func legend(c *ui.Context, name string, ink ui.Color) {
	t := c.Theme()
	ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
		ui.Box(c).Size(t.Space(3), t.Space(1)).Radius(999).Background(ink)
		kit.Muted(c, name).FontSize(t.Rem(0.7)).Shrink(0)
	})
}
