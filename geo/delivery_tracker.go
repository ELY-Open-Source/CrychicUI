package geo

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DeliveryTrackerOptions word the courier line; ETA formats minutes.
type DeliveryTrackerOptions struct {
	Courier, Role, ETA, Arrived string
	Photo                       *ui.Bitmap
	Minutes                     int
}

// DeliveryTracker moves a courier along a route by progress.
func DeliveryTracker(c *ui.Context, route []Point, progress float64, o DeliveryTrackerOptions) ui.Element {
	def(&o.Role, "Your courier")
	def(&o.ETA, "Arriving in %d min")
	def(&o.Arrived, "Delivered")
	if len(route) < 2 || strings.TrimSpace(o.Courier) == "" || !(progress >= 0 && progress <= 1) || o.Minutes < 0 {
		panic(fmt.Sprintf("geo: DeliveryTracker %q needs a courier, two route points, progress 0 to 1 and minutes >= 0", o.Courier))
	}
	for _, q := range route {
		valid("DeliveryTracker", o.Courier, q)
	}
	t := c.Theme()
	root := card(c).Clip().Role(ui.RoleGroup)
	return root.Children(func() {
		c.Key("map")
		m := ui.Box(c).FillWidth().Height(t.Space(44))
		done := m.AnimateTo("done", float32(progress), ui.Smooth)
		m.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			center, zoom := fit(route, r)
			terrain(p, r, t, center, zoom)
			xs, ys := make([]float32, len(route)), make([]float32, len(route))
			var line ui.Path
			total := float32(0)
			for i, q := range route {
				xs[i], ys[i] = project(center, q, zoom, r)
				if i == 0 {
					line.MoveTo(xs[i], ys[i])
				} else {
					line.LineTo(xs[i], ys[i])
					total += float32(math.Hypot(float64(xs[i]-xs[i-1]), float64(ys[i]-ys[i-1])))
				}
			}
			p.StrokePath(&line, t.Space(1.5), t.Text.Alpha(0.22))
			p.StrokePath(line.Trim(0, done), t.Space(1.5), t.Accent)
			// The courier sits done of the way along.
			x, y, left := xs[0], ys[0], done*total
			for i := 1; i < len(route); i++ {
				seg := float32(math.Hypot(float64(xs[i]-xs[i-1]), float64(ys[i]-ys[i-1])))
				if left <= seg || i == len(route)-1 {
					f := min(1, left/max(seg, 1e-6))
					x, y = xs[i-1]+(xs[i]-xs[i-1])*f, ys[i-1]+(ys[i]-ys[i-1])*f
					break
				}
				left -= seg
			}
			var start ui.Path
			p.FillPath(start.Circle(xs[0], ys[0], t.Space(1.25)), t.Success)
			pin(p, xs[len(xs)-1], ys[len(ys)-1], t.Space(4), t.Danger, t.Background)
			var halo, dot ui.Path
			p.FillPath(halo.Circle(x, y, t.Space(4)), t.Accent.Alpha(0.2))
			p.FillPath(dot.Circle(x, y, t.Space(2.25)), t.Background)
			var core ui.Path
			p.FillPath(core.Circle(x, y, t.Space(1.6)), t.Accent)
		}))
		c.Key("courier")
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3), t.Space(4)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
				display.Avatar(c, o.Courier, o.Photo).Size(t.Space(10), t.Space(10)).Shrink(0)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
					ui.Text(c, o.Courier).FontSize(t.Rem(0.86)).Bold().FillWidth().MinWidth(0).SingleLine()
					kit.Muted(c, o.Role).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				})
				eta, ink := fmt.Sprintf(o.ETA, o.Minutes), t.Accent
				if progress == 1 {
					eta, ink = o.Arrived, t.Success
				}
				kit.Numeric(c, eta).FontSize(t.Rem(0.78)).Bold().TextColor(ink).Shrink(0).MaxWidth(t.Space(40)).SingleLine()
			})
			bar := ui.Box(c).FillWidth().Height(t.Space(1.25)).Radius(999).Background(t.SurfaceHover)
			bar.Draw(func(p *ui.Painter, r ui.Rect) {
				if done > 0 {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * done, H: r.H}, t.Accent.Mix(t.Success, unit(progress == 1)), r.H/2)
				}
			})
		})
	})
}
