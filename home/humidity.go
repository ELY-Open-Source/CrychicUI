package home

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HumidityOptions set the comfort band; zero means default.
type HumidityOptions struct {
	Title, Room, Unit, Dry, Comfortable, Humid string
	Low, High                                  float64
}

// Humidity gauges relative humidity against a comfort band.
func Humidity(c *ui.Context, percent float64, o HumidityOptions) ui.Element {
	def(&o.Title, "Humidity")
	def(&o.Unit, "%")
	def(&o.Dry, "Dry")
	def(&o.Comfortable, "Comfortable")
	def(&o.Humid, "Humid")
	if o.Low == 0 && o.High == 0 {
		o.Low, o.High = 40, 60
	}
	if !(percent >= 0 && percent <= 100) || !(o.Low >= 0 && o.Low < o.High && o.High <= 100) {
		panic(fmt.Sprintf("home: Humidity %q needs 0 to 100 percent and a band inside it", o.Title))
	}
	t := c.Theme()
	state, ink := o.Comfortable, t.Success
	switch {
	case percent < o.Low:
		state, ink = o.Dry, t.Warning
	case percent > o.High:
		state, ink = o.Humid, t.Accent
	}
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "droplets", o.Title, nil)
		ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(1)).Padding(t.Space(3), t.Space(3), t.Space(4), t.Space(3)).Children(func() {
			w, h := t.Space(48), t.Space(26)
			gauge := ui.Box(c).Size(w, h).Shrink(0)
			grow := gauge.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			k := grow * gauge.AnimateTo("share", float32(percent/100), ui.Smooth)
			col := gauge.AnimateColor("ink", ink, ui.Snappy)
			gauge.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				sw := t.Space(2.5)
				cx, cy := r.X+r.W/2, r.Y+r.H-sw/2
				rr := min(r.W/2, r.H) - sw
				p.StrokePath(arc(cx, cy, rr, -90, 90), sw, t.SurfaceHover)
				p.StrokePath(arc(cx, cy, rr, float32(-90+180*o.Low/100), float32(-90+180*o.High/100)), sw, t.Success.Alpha(0.35))
				a := float64(-90+180*k) * math.Pi / 180
				var needle ui.Path
				needle.MoveTo(cx, cy).LineTo(cx+(rr-sw)*float32(math.Sin(a)), cy-(rr-sw)*float32(math.Cos(a)))
				p.StrokePath(&needle, t.Space(0.75), col)
				p.FillPath(ring(cx, cy, t.Space(1.25)), col)
			}))
			ui.Row(c).AlignItems(ui.Start).Gap(t.Space(0.5)).Children(func() {
				kit.Numeric(c, fmt.Sprint(int(math.Round(percent)))).FontSize(t.Rem(1.8)).Bold().Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				kit.Muted(c, o.Unit).FontSize(t.Rem(0.8)).Shrink(0)
			})
			c.Key(state)
			chip := ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2.5)).Radius(999).
				Background(ink.Alpha(0.12)).TextColor(ink)
			kit.Arrive(chip, ui.Pose{Scale: 0.8, Fade: 1}, 0).Children(func() {
				ui.Text(c, state).FontSize(t.Rem(0.76)).Bold().Shrink(0)
			})
			if o.Room != "" {
				kit.Muted(c, o.Room).FontSize(t.Rem(0.72)).Shrink(0)
			}
		})
	})
	return root
}
