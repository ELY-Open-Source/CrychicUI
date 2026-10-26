package geo

import (
	"image"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestGeofenceRejectsSteps(t *testing.T) {
	for _, step := range []float64{math.Inf(1), 5000} {
		panics(t, "steps", func(c *ui.Context) { r := 500.0; Geofence(c, studio, &r, GeofenceOptions{Step: step}) })
	}
}

func TestMapCanvasPanWrapsAndClamps(t *testing.T) {
	view := View{Center: Point{Lat: 85, Lng: 0}, Zoom: 3}
	changes := 0
	tt := ui.NewTester(padded(within(400, func(c *ui.Context) {
		if MapCanvas(c, &view, nil, MapCanvasOptions{MinZoom: 3}).Changed() {
			changes++
		}
	})), 520, 500)
	settles(t, tt)
	r := rectOf(t, tt, "Map")
	x, y := r.X+r.W/2, r.Y+r.H/2
	drag := func(dx, dy float32) {
		tt.Press(x, y)
		tt.Move(x+dx, y+dy)
		tt.Frame()
		tt.Release(x+dx, y+dy)
		tt.Frame()
	}
	drag(0, 60)
	if changes != 0 || view.Center != (Point{Lat: 85, Lng: 0}) {
		t.Errorf("a pan clamped at the pole reported %d changes, center %v", changes, view.Center)
	}
	// One long drag passes 540 degrees here.
	x = r.X + 10
	drag(400, 0)
	if !(view.Center.Lng >= -180 && view.Center.Lng <= 180) || changes == 0 {
		t.Errorf("a long pan left longitude %v, %d changes", view.Center.Lng, changes)
	}
}

func TestAntimeridian(t *testing.T) {
	r := ui.Rect{X: 0, Y: 0, W: 400, H: 300}
	x, _ := project(Point{Lng: -179.99}, Point{Lng: 179.99}, 12, r)
	if x < 0 || x > 400 {
		t.Errorf("a pin 0.02° away across the antimeridian projects to x %v", x)
	}
	mid, zoom := fit([]Point{{Lng: 179.99}, {Lng: -179.99}}, r)
	if math.Abs(math.Abs(mid.Lng)-180) > 0.001 || zoom < 10 {
		t.Errorf("fitting across the antimeridian gave center %v zoom %v", mid, zoom)
	}
}

func TestLongCopyStaysInside(t *testing.T) {
	u := Units{Kilometers: "Approximately %s kilometers", Meters: "Approximately %s meters"}
	walk := "Estimated walking time %d minutes"
	for _, c := range []struct {
		name  string
		view  func(c *ui.Context)
		texts []string
	}{
		{"Directions", func(c *ui.Context) {
			Directions(c, []Step{{Text: "Head north", Turn: "arrow-up", Meters: 2000, Minutes: 7}}, DirectionsOptions{Time: walk, Units: u})
		}, []string{"Directions", "Approximately 2.0 kilometers", "Estimated walking time 7 minutes"}},
		{"Distance", func(c *ui.Context) {
			Distance(c, 2000, DistanceOptions{Minutes: 7, Time: walk, Units: u})
		}, []string{"Approximately 2.0 kilometers", "Estimated walking time 7 minutes"}},
		{"Elevation", func(c *ui.Context) {
			Elevation(c, []float64{100, 180}, 2000, ElevationOptions{High: "Highest point on the whole route %s", Low: "Lowest point on the whole route %s", Units: u})
		}, []string{"Highest point on the whole route Approximately 180 meters", "Lowest point on the whole route Approximately 100 meters", "Approximately 2.0 kilometers"}},
		{"Geofence", func(c *ui.Context) {
			r := 500.0
			Geofence(c, studio, &r, GeofenceOptions{Units: Units{Meters: "Approximately %s meters from the studio", Kilometers: "Approximately %s kilometers from the studio"}})
		}, []string{"Approximately 100 meters from the studio", "Approximately 2.0 kilometers from the studio"}},
	} {
		tt := ui.NewTester(padded(within(320, c.view)), 520, 900)
		settles(t, tt)
		for _, s := range c.texts {
			r, ok := tt.Find(s)
			if !ok {
				t.Errorf("%s: %q missing: %q", c.name, s, tt.Texts())
				continue
			}
			inside(t, c.name+" "+s, r, 10, 330)
		}
		if c.name == "Geofence" {
			lo, _ := tt.Find(c.texts[0])
			hi, _ := tt.Find(c.texts[1])
			if room := hi.X - (lo.X + lo.W); room < 80 {
				t.Errorf("Geofence leaves the slider %v DIP", room)
			}
		}
	}
}

func TestTerrainSeamless(t *testing.T) {
	shot := func(lat, lng float64) *image.RGBA {
		view := View{Center: Point{Lat: lat, Lng: lng}, Zoom: 14}
		tt := ui.NewTester(padded(within(400, func(c *ui.Context) { MapCanvas(c, &view, nil, MapCanvasOptions{}) })), 420, 400)
		settles(t, tt)
		return tt.Image()
	}
	// 37.75 shows parks; 37.72 shows a river band.
	for _, lat := range []float64{37.75, 37.72} {
		a, b := shot(lat, 180), shot(lat, -180)
		diff := 0
		for i := 0; i < len(a.Pix); i += 4 {
			if a.Pix[i] != b.Pix[i] || a.Pix[i+1] != b.Pix[i+1] || a.Pix[i+2] != b.Pix[i+2] {
				diff++
			}
		}
		if diff > 200 {
			t.Errorf("at latitude %v the map at +180° and -180° differs in %d pixels", lat, diff)
		}
	}
}

func TestPinPaintsInsideItsBox(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Pin(c, "Studio", PinOptions{}) }), 400, 300)
	settles(t, tt)
	b := rectOf(t, tt, "Studio")
	img := tt.Image()
	bg := img.RGBAAt(2, 2)
	// The box sits under 4 DIP of padding, 64 tall.
	bottom := int(b.Y + 4 + 64)
	for y := bottom; y < bottom+2; y++ {
		for x := int(b.X); x < int(b.X+b.W); x++ {
			if img.RGBAAt(x, y) != bg {
				t.Fatalf("pin paint at %d,%d below its box", x, y)
			}
		}
	}
}
