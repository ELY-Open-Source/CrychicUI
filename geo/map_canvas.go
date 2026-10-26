package geo

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Point is a place in degrees.
type Point struct{ Lat, Lng float64 }

// View is a map's center and zoom.
type View struct {
	Center Point
	Zoom   float64
}

// MapPin is a named place painted on a map.
type MapPin struct {
	Name string
	At   Point
}

// MapCanvasOptions size and word the map; zoom 10 to 18.
type MapCanvasOptions struct {
	Label, ZoomIn, ZoomOut string
	Height                 float32
	MinZoom, MaxZoom       float64
	Units                  Units
}

// MapCanvas is a drawn map; drag pans; cue Tick.
func MapCanvas(c *ui.Context, view *View, pins []MapPin, o MapCanvasOptions) ui.Element {
	def(&o.Label, "Map")
	def(&o.ZoomIn, "Zoom in")
	def(&o.ZoomOut, "Zoom out")
	if o.Height == 0 {
		o.Height = 64
	}
	if o.MinZoom == 0 {
		o.MinZoom = 10
	}
	if o.MaxZoom == 0 {
		o.MaxZoom = 18
	}
	o.Units.fill()
	if view == nil || !(o.MinZoom >= 3 && o.MinZoom < o.MaxZoom && o.MaxZoom <= 20) || !(view.Zoom >= o.MinZoom && view.Zoom <= o.MaxZoom) {
		panic(fmt.Sprintf("geo: MapCanvas %q needs a view zoomed within 3 <= MinZoom < MaxZoom <= 20", o.Label))
	}
	valid("MapCanvas", o.Label, view.Center)
	names := make([]string, len(pins))
	for i, m := range pins {
		valid("MapCanvas", m.Name, m.At)
		names[i] = m.Name
	}
	unique("MapCanvas", names)
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(t.Space(o.Height)).Radius(t.Radius*1.5).Clip().Border(1, t.Border).
		Cursor(ui.CursorGrab).Role(ui.RoleGroup).Label(o.Label)
	kit.Arrive(root, ui.Pose{Scale: 0.98, Fade: 1}, 0)
	if dx, dy, held := root.Dragged(); held && (dx != 0 || dy != 0) {
		k := scale(view.Zoom)
		lat := min(85, max(-85, view.Center.Lat+float64(dy)/k))
		lng := wrap(view.Center.Lng - float64(dx)/(k*math.Cos(view.Center.Lat*math.Pi/180)))
		if moved := (Point{lat, lng}); moved != view.Center {
			view.Center = moved
			root.SetChanged()
			kit.Log("MapCanvas", "panned")
		}
	}
	zoom := root.AnimateTo("zoom", float32(view.Zoom), ui.Smooth)
	center := view.Center
	root.Draw(kit.Inside(t.Radius*1.5, func(p *ui.Painter, r ui.Rect) {
		terrain(p, r, t, center, float64(zoom))
		for _, m := range pins {
			x, y := project(center, m.At, float64(zoom), r)
			pin(p, x, y, t.Space(3), t.Danger, t.Background)
		}
	}))
	root.Children(func() {
		c.Key("zoom")
		ui.Column(c).Absolute().Top(t.Space(2)).Right(t.Space(2)).Gap(t.Space(1)).Children(func() {
			for i, step := range [2]float64{1, -1} {
				label, mark := o.ZoomIn, "plus"
				if i == 1 {
					label, mark = o.ZoomOut, "minus"
				}
				next := min(o.MaxZoom, max(o.MinZoom, view.Zoom+step))
				b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius).
					Border(1, t.Border).Disabled(next == view.Zoom))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				b.Shadow(0, 1, 3, 0, kit.Shade(0.1)).Children(func() { glyph(c, mark, 16) })
				if b.Clicked() && next != view.Zoom {
					view.Zoom = next
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("MapCanvas", "zoomed", "zoom", next)
				}
			}
		})
		// The scale bar spans a round distance under 96 DIP.
		mpp := metersPerDegree / scale(view.Zoom)
		long := 10.0
		for _, n := range []float64{20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000, 50000} {
			if n/mpp <= float64(t.Space(24)) {
				long = n
			}
		}
		c.Key("scale")
		ui.Column(c).Absolute().Left(t.Space(2)).Bottom(t.Space(2)).Gap(t.Space(0.5)).Padding(t.Space(1), t.Space(1.5)).
			Radius(t.Radius).Background(t.Background.Alpha(0.85)).Children(func() {
			ui.Text(c, span(long, false, o.Units)).FontSize(t.Rem(0.66)).Bold().Shrink(0).MaxWidth(t.Space(30)).SingleLine().ContentTransition(ui.ContentTransition{})
			bar := ui.Box(c).Height(t.Space(0.75)).Radius(999).Background(t.Text)
			bar.Width(bar.AnimateTo("w", float32(long/mpp), ui.Smooth))
		})
	})
	return root
}
