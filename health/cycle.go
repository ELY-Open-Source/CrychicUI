package health

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CycleOptions word the ring; Period is bleeding days, default 5.
type CycleOptions struct {
	Title, Day, Next, Tomorrow string
	Phases                     [4]string
	Period                     int
}

// Cycle marks today on a ring of the cycle's days.
func Cycle(c *ui.Context, day, length int, o CycleOptions) ui.Element {
	def(&o.Title, "Cycle")
	if o.Period == 0 {
		o.Period = 5
	}
	if length < 21 || length > 45 || day < 1 || day > length || o.Period < 2 || o.Period > 8 {
		panic(fmt.Sprintf("health: Cycle %q needs 21 to 45 days, a day within it and 2 to 8 period days", o.Title))
	}
	def(&o.Day, "Day %d")
	def(&o.Next, "Period in %d days")
	def(&o.Tomorrow, "Period tomorrow")
	for i, p := range [...]string{"Period", "Follicular", "Fertile", "Luteal"} {
		def(&o.Phases[i], p)
	}
	t := c.Theme()
	ovulation := length - 14
	phase := func(d int) int {
		switch {
		case d <= o.Period:
			return 0
		case d >= ovulation-5 && d <= ovulation+1:
			return 2
		case d < ovulation-5:
			return 1
		}
		return 3
	}
	inks := [4]ui.Color{t.Danger, t.SurfaceHover, t.Success, t.SurfaceHover.Mix(t.Border, 0.5)}
	now := phase(day)
	ink := [4]ui.Color{t.Danger, t.Accent, t.Success, t.Warning}[now]
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Title).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine()
		size := t.Space(44)
		wheel := ui.Box(c).Size(size, size).Shrink(0).Center()
		grow := wheel.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		at := wheel.AnimateTo("day", float32(day), ui.Smooth)
		wheel.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			w := size * 0.1
			cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2-w
			gap := float32(0.004)
			for d := 1; d <= length; d++ {
				a, b := float32(d-1)/float32(length), float32(d)/float32(length)
				if a >= grow {
					break
				}
				p.StrokePath(ring(cx, cy, rr).Trim(a+gap, min(b, grow)-gap), w, inks[phase(d)])
			}
			th := 2 * math.Pi * (float64(at) - 0.5) / float64(length)
			mx, my := cx+rr*float32(math.Sin(th)), cy-rr*float32(math.Cos(th))
			p.FillPath(ring(mx, my, w*0.75), t.Background)
			p.StrokePath(ring(mx, my, w*0.75), t.Space(0.5), ink)
		}))
		wheel.Children(func() {
			ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				kit.Numeric(c, fmt.Sprintf(o.Day, day)).FontSize(t.Rem(1.4)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
				c.Key(now)
				kit.Arrive(kit.Chip(c, o.Phases[now], ink), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			})
		})
		left := length - day + 1
		next := fmt.Sprintf(o.Next, left)
		if left == 1 {
			next = o.Tomorrow
		}
		kit.Muted(c, next).FontSize(t.Rem(0.8)).FillWidth().TextAlign(ui.Center).Wrap()
	})
	return root
}
