package geo

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// GeofenceOptions word the fence; radii run Min to Max meters.
type GeofenceOptions struct {
	Label          string
	Min, Max, Step float64
	Imperial       bool
	Units          Units
}

// Geofence sizes a radius with a slider; cue Tick.
func Geofence(c *ui.Context, center Point, radius *float64, o GeofenceOptions) ui.Element {
	def(&o.Label, "Radius")
	if o.Min == 0 {
		o.Min = 100
	}
	if o.Max == 0 {
		o.Max = 2000
	}
	if o.Step == 0 {
		o.Step = 50
	}
	o.Units.fill()
	if !(o.Min > 0 && o.Max > o.Min && o.Max <= 1e6 && o.Step > 0 && o.Step <= o.Max-o.Min && (o.Max-o.Min)/o.Step <= 1e4) {
		panic(fmt.Sprintf("geo: Geofence %q needs 0 < Min < Max <= 1,000 km and 1 to 10,000 steps", o.Label))
	}
	if radius == nil || !(*radius >= o.Min && *radius <= o.Max) {
		panic(fmt.Sprintf("geo: Geofence %q needs a radius from %v to %v", o.Label, o.Min, o.Max))
	}
	valid("Geofence", o.Label, center)
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "locate-fixed", t.Accent, t.Space(9))
			ui.Text(c, o.Label).FontSize(t.Rem(0.9)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Numeric(c, span(*radius, o.Imperial, o.Units)).FontSize(t.Rem(0.9)).Bold().TextColor(t.Accent).Shrink(0).MaxWidth(t.Space(40)).SingleLine()
		})
		c.Key("map")
		m := ui.Box(c).FillWidth().Height(t.Space(48)).Radius(t.Radius).Border(1, t.Border)
		rr := m.AnimateTo("radius", float32(*radius), ui.Smooth)
		m.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			// The largest fence fills 45% of the shorter side.
			k := float64(min(r.W, r.H)) * 0.45 / (o.Max / metersPerDegree)
			zoom := math.Log2(k * 360 / 256)
			terrain(p, r, t, center, zoom)
			x, y := project(center, center, zoom, r)
			px := float32(float64(rr) / metersPerDegree * k)
			var fence ui.Path
			fence.Circle(x, y, px)
			p.FillPath(&fence, t.Accent.Alpha(0.16))
			p.StrokePath(&fence, t.Space(0.5), t.Accent)
			pin(p, x, y, t.Space(3.5), t.Accent, t.Background)
		}))
		c.Key("slider")
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			kit.Muted(c, span(o.Min, o.Imperial, o.Units)).FontSize(t.Rem(0.7)).Shrink(0).MaxWidth(t.Space(16)).SingleLine()
			sl := ui.Slider(c, radius, o.Min, o.Max).Step(o.Step).Label(o.Label).Grow(1).MinWidth(0)
			if sl.Changed() {
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("Geofence", "radius", "meters", int(math.Round(*radius)))
			}
			kit.Muted(c, span(o.Max, o.Imperial, o.Units)).FontSize(t.Rem(0.7)).Shrink(0).MaxWidth(t.Space(16)).SingleLine()
		})
	})
}
