package home

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Pollutant is one reading; Value shows as given.
type Pollutant struct{ Name, Value, Unit string }

// AirQualityOptions word the card; Bands are good, moderate, poor, hazardous.
type AirQualityOptions struct {
	Title, Index string
	Bands        [4]string
	Pollutants   []Pollutant
}

// AirQuality shows the air quality index on a ring.
func AirQuality(c *ui.Context, aqi int, o AirQualityOptions) ui.Element {
	def(&o.Title, "Air quality")
	def(&o.Index, "AQI")
	if aqi < 0 || aqi > 500 {
		panic(fmt.Sprintf("home: AirQuality %q needs an index of 0 to 500, got %d", o.Title, aqi))
	}
	names := make([]string, len(o.Pollutants))
	for i, p := range o.Pollutants {
		if strings.TrimSpace(p.Value) == "" {
			panic(fmt.Sprintf("home: AirQuality pollutant %q needs a value", p.Name))
		}
		names[i] = p.Name
	}
	unique("AirQuality", names)
	for i, s := range [...]string{"Good", "Moderate", "Poor", "Hazardous"} {
		def(&o.Bands[i], s)
	}
	t := c.Theme()
	band := 0
	switch {
	case aqi > 200:
		band = 3
	case aqi > 100:
		band = 2
	case aqi > 50:
		band = 1
	}
	ink := [4]ui.Color{t.Success, t.Warning, t.Danger, t.Danger.Mix(t.Text, 0.3)}[band]
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "wind", o.Title, nil)
		ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(3)).Children(func() {
			size := t.Space(26)
			dial := ui.Box(c).Size(size, size).Shrink(0).Center()
			grow := dial.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			k := grow * dial.AnimateTo("share", min(1, float32(aqi)/300), ui.Smooth)
			col := dial.AnimateColor("ink", ink, ui.Snappy)
			dial.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				w := t.Space(2)
				cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2-w/2
				p.StrokePath(ring(cx, cy, rr), w, t.SurfaceHover)
				if k > 0.005 {
					p.StrokePath(arc(cx, cy, rr, 0, 360*k), w, col)
				}
			}))
			dial.Children(func() {
				ui.Column(c).AlignItems(ui.Center).Children(func() {
					kit.Numeric(c, fmt.Sprint(aqi)).FontSize(t.Rem(1.5)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Muted(c, o.Index).FontSize(t.Rem(0.66)).Shrink(0)
				})
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(2)).Children(func() {
				ui.Row(c).Children(func() {
					c.Key(band)
					chip := ui.Row(c).Shrink(0).Padding(t.Space(0.5), t.Space(2.5)).Radius(999).Background(ink.Alpha(0.12)).TextColor(ink)
					kit.Arrive(chip, ui.Pose{Scale: 0.8, Fade: 1}, 0).Children(func() {
						ui.Text(c, o.Bands[band]).FontSize(t.Rem(0.78)).Bold().Shrink(0)
					})
				})
				grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(1.5))
				grid.Children(func() {
					for i, p := range o.Pollutants {
						cell := ui.Column(c).MinWidth(0).Gap(t.Space(0.25)).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
						kit.Stagger(c, grid, cell, i).Children(func() {
							kit.Muted(c, p.Name).FontSize(t.Rem(0.66)).FillWidth().MinWidth(0).SingleLine()
							ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
								ui.Text(c, p.Value).FontSize(t.Rem(0.82)).Bold().MinWidth(0).SingleLine()
								kit.Muted(c, p.Unit).FontSize(t.Rem(0.6)).Shrink(0).MaxWidth(t.Space(12)).SingleLine()
							})
						})
					}
				})
			})
		})
	})
	return root
}
