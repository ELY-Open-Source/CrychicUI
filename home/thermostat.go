package home

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThermostatOptions set the range; zero means default.
type ThermostatOptions struct {
	Min, Max, Step                                    float64
	Unit, Now, Heating, Cooling, Idle, Warmer, Cooler string
}

// Thermostat dials a target temperature; cue Tick.
func Thermostat(c *ui.Context, name string, target *float64, current float64, o ThermostatOptions) ui.Element {
	if o.Min == 0 && o.Max == 0 {
		o.Min, o.Max = 10, 30
	}
	if o.Step == 0 {
		o.Step = 0.5
	}
	if strings.TrimSpace(name) == "" || target == nil || !(o.Max > o.Min) || math.IsInf(o.Max-o.Min, 0) ||
		!(o.Step > 0 && o.Step <= o.Max-o.Min && (o.Max-o.Min)/o.Step <= 1e6) {
		panic(fmt.Sprintf("home: Thermostat %q needs a target and a finite range with a step of at most 1e6 per range", name))
	}
	if !(*target >= o.Min && *target <= o.Max) || !(current >= -60 && current <= 60) {
		panic(fmt.Sprintf("home: Thermostat %q target %v must sit in %v..%v, current in -60..60", name, *target, o.Min, o.Max))
	}
	def(&o.Unit, "°C")
	def(&o.Now, "Now %s°")
	def(&o.Heating, "Heating")
	def(&o.Cooling, "Cooling")
	def(&o.Idle, "Idle")
	def(&o.Warmer, "Warmer")
	def(&o.Cooler, "Cooler")
	t := c.Theme()
	deg := func(v float64) string { return strconv.FormatFloat(v, 'f', 1, 64) }
	mode, mark, ink := o.Idle, "pause", t.TextMuted
	switch {
	case *target > current+0.25:
		mode, mark, ink = o.Heating, "flame", t.Danger
	case *target < current-0.25:
		mode, mark, ink = o.Cooling, "snowflake", t.Accent
	}
	share := float32((*target - o.Min) / (o.Max - o.Min))
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, name).FontSize(t.Rem(0.92)).Bold().Grow(1).MinWidth(0).SingleLine()
			c.Key(mode)
			chip := ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).
				Background(ink.Alpha(0.12)).TextColor(ink)
			kit.Arrive(chip, ui.Pose{Scale: 0.8, Fade: 1}, 0).Children(func() {
				glyph(c, mark, 12)
				ui.Text(c, mode).FontSize(t.Rem(0.7)).Bold().Shrink(0)
			})
		})
		size := t.Space(44)
		dial := ui.Box(c).Size(size, size).Shrink(0).Center()
		k := dial.AnimateTo("share", share, ui.Smooth)
		grow := dial.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		col := dial.AnimateColor("ink", t.Accent.Mix(t.Danger, share), ui.Snappy)
		dial.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			w := t.Space(3)
			cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2-w
			p.StrokePath(arc(cx, cy, rr, -135, 135), w, t.SurfaceHover)
			end := -135 + 270*k*grow
			if end > -134 {
				p.StrokePath(arc(cx, cy, rr, -135, end), w, col)
			}
			a := float64(end) * math.Pi / 180
			kx, ky := cx+rr*float32(math.Sin(a)), cy-rr*float32(math.Cos(a))
			p.FillPath(ring(kx, ky, w*0.9), t.Background)
			p.StrokePath(ring(kx, ky, w*0.9), t.Space(0.5), col)
		}))
		dial.Children(func() {
			ui.Column(c).AlignItems(ui.Center).Gap(t.Space(0.5)).Children(func() {
				ui.Row(c).AlignItems(ui.Start).Gap(t.Space(0.5)).Children(func() {
					kit.Numeric(c, deg(*target)).FontSize(t.Rem(2.2)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Muted(c, o.Unit).FontSize(t.Rem(0.8)).Shrink(0)
				})
				kit.Muted(c, fmt.Sprintf(o.Now, deg(current))).FontSize(t.Rem(0.74)).Shrink(0)
			})
		})
		ui.Row(c).Gap(t.Space(6)).AlignItems(ui.Center).Children(func() {
			step := func(label, mark string, up bool) {
				// Steps count from Min; the ends stay reachable.
				k := (*target - o.Min) / o.Step
				n := math.Ceil(k-1e-9) - 1
				if up {
					n = math.Floor(k+1e-9) + 1
				}
				next := math.Max(o.Min, math.Min(o.Max, o.Min+n*o.Step))
				ok := next != *target
				c.Key(label)
				b := round(c, label, mark, t.Space(11), t.SurfaceHover, t.Text).Disabled(!ok)
				if b.Clicked() && ok {
					*target = next
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("Thermostat", "target", "tenths", int(math.Round(*target*10)))
				}
			}
			step(o.Cooler, "minus", false)
			step(o.Warmer, "plus", true)
		})
	})
	return root
}
