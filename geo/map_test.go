package geo

import (
	"math"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var studio = Point{Lat: 37.7599, Lng: -122.4148}

func TestMapCanvasZoomAndPan(t *testing.T) {
	view := View{Center: studio, Zoom: 17}
	changes := 0
	tt := ui.NewTester(padded(within(400, func(c *ui.Context) {
		if MapCanvas(c, &view, []MapPin{{Name: "Studio", At: studio}}, MapCanvasOptions{}).Changed() {
			changes++
		}
	})), 520, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Zoom in")
	if view.Zoom != 18 || changes != 1 {
		t.Fatalf("zoom %v after one press, %d changes", view.Zoom, changes)
	}
	cues(t, before, sound.Tick)
	before = baseline()
	click(t, tt, "Zoom in")
	if view.Zoom != 18 || changes != 1 {
		t.Errorf("zoom passed its maximum: %v, %d changes", view.Zoom, changes)
	}
	cues(t, before)
	r := rectOf(t, tt, "Map")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Move(x+40, y+20)
	tt.Frame()
	tt.Release(x+40, y+20)
	tt.Frame()
	if !(view.Center.Lng < studio.Lng && view.Center.Lat > studio.Lat) {
		t.Errorf("dragging right and down moved the center to %v", view.Center)
	}
}

func TestMapCanvasRejects(t *testing.T) {
	panics(t, "MapCanvas", func(c *ui.Context) {
		MapCanvas(c, &View{Center: studio, Zoom: 2}, nil, MapCanvasOptions{})
	})
	panics(t, "repeated", func(c *ui.Context) {
		MapCanvas(c, &View{Center: studio, Zoom: 14}, []MapPin{{Name: "A1", At: studio}, {Name: "A1", At: studio}}, MapCanvasOptions{})
	})
	panics(t, "latitude", func(c *ui.Context) {
		MapCanvas(c, &View{Center: Point{Lat: 89, Lng: 0}, Zoom: 14}, nil, MapCanvasOptions{})
	})
}

func TestPinDrops(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Pin(c, "Studio", PinOptions{Tint: 1}).Clicked() {
			clicks++
		}
	}), 400, 300)
	settles(t, tt)
	if tt.FramePending() {
		t.Fatal("a resting pin keeps animating")
	}
	// top is the highest row of the pin's green.
	box := rectOf(t, tt, "Studio")
	top := func() int {
		img := tt.Image()
		for y := int(box.Y); y < int(box.Y+box.H); y++ {
			for x := int(box.X); x < int(box.X+box.W); x++ {
				if p := img.RGBAAt(x, y); p.G > 120 && p.R < 90 && p.B < 140 {
					return y
				}
			}
		}
		return -1
	}
	rest := top()
	before := baseline()
	click(t, tt, "Studio")
	cues(t, before, sound.Drop)
	tt.Advance(40 * 1e6)
	if up := top(); up < 0 || up > rest-10 {
		t.Errorf("a pressed pin did not rise to drop again: head at %d, resting %d", up, rest)
	}
	if clicks != 1 {
		t.Errorf("%d clicks", clicks)
	}
	panics(t, "tint", func(c *ui.Context) { Pin(c, "Studio", PinOptions{Tint: 4}) })
}

func TestClusterCount(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Cluster(c, 27, ClusterOptions{})
		Cluster(c, 1840, ClusterOptions{})
	}), 400, 300)
	settles(t, tt)
	if !exact(tt, "27") || !exact(tt, "999+") || exact(tt, "1840") {
		t.Errorf("cluster texts %q", tt.Texts())
	}
	if _, ok := tt.Find("1840 places"); !ok {
		t.Error("cluster is not named by its full count")
	}
	panics(t, "at least 2", func(c *ui.Context) { Cluster(c, 1, ClusterOptions{}) })
}

func searchPlaces() []Place {
	return []Place{
		{Name: "Riverbend Clay Co.", Address: "2125 Bryant St", At: studio},
		{Name: "Mission Clay Supply", Address: "2414 Mission St", At: Point{Lat: 37.7650, Lng: -122.4190}},
		{Name: "Dolores Park Café", Address: "501 Dolores St", At: Point{Lat: 37.7610, Lng: -122.4260}},
	}
}

func TestLocationSearchFilters(t *testing.T) {
	query, picked := "", ""
	places := searchPlaces()
	tt := ui.NewTester(padded(func(c *ui.Context) {
		LocationSearch(c, &query, places, &picked, LocationSearchOptions{From: &studio})
	}), 520, 700)
	settles(t, tt)
	if !exact(tt, "0 m") {
		t.Errorf("distance from the studio to itself: %q", tt.Texts())
	}
	click(t, tt, "Search places")
	tt.Type("dolores")
	tt.Frame()
	settles(t, tt)
	if _, ok := tt.Find("Mission Clay Supply"); ok {
		t.Error("a place matching neither name nor address stayed")
	}
	before := baseline()
	click(t, tt, "Dolores Park Café")
	if picked != "Dolores Park Café" {
		t.Fatalf("picked %q", picked)
	}
	cues(t, before, sound.Select)
	before = baseline()
	click(t, tt, "Dolores Park Café")
	cues(t, before)
	click(t, tt, "Search places")
	tt.Type("zzz")
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "No places found") {
		t.Errorf("no empty state: %q", tt.Texts())
	}
}

func TestLocationSearchKeepsIdentity(t *testing.T) {
	query, picked := "", ""
	places := searchPlaces()
	tt := ui.NewTester(padded(func(c *ui.Context) {
		LocationSearch(c, &query, places, &picked, LocationSearchOptions{})
	}), 520, 700)
	settles(t, tt)
	swapUnderPress(t, tt, "Mission Clay Supply", func() { places[0], places[1] = places[1], places[0] })
	if picked == "Riverbend Clay Co." {
		t.Error("the place that slid under the pointer was picked")
	}
}

func TestPlaceCard(t *testing.T) {
	p := Place{Name: "Dolores Park Café", Address: "501 Dolores St", Hours: "7:00–19:00", At: studio, Rating: 4.6, Reviews: 1284, Open: true}
	tt := ui.NewTester(padded(func(c *ui.Context) { PlaceCard(c, p, PlaceCardOptions{}) }), 520, 600)
	settles(t, tt)
	for _, s := range []string{"4.6", "(1,284)", "Open", "7:00–19:00"} {
		if !exact(tt, s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	p.Open, p.Rating = false, 0
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "Closed") || exact(tt, "4.6") {
		t.Errorf("closed, unrated card shows %q", tt.Texts())
	}
	panics(t, "rating", func(c *ui.Context) { PlaceCard(c, Place{Name: "Café", At: studio, Rating: 6}, PlaceCardOptions{}) })
}

func TestDirectionsTotals(t *testing.T) {
	steps := []Step{
		{Text: "Head north on Bryant St", Turn: "arrow-up", Meters: 350, Minutes: 30},
		{Text: "Turn left onto 20th St", Turn: "corner-up-left", Meters: 1650, Minutes: 38},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { Directions(c, steps, DirectionsOptions{}) }), 520, 700)
	settles(t, tt)
	for _, s := range []string{"2.0 km", "1 h 8 min", "350 m", "You have arrived"} {
		if !exact(tt, s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	panics(t, "text", func(c *ui.Context) { Directions(c, []Step{{Text: " ", Turn: "arrow-up"}}, DirectionsOptions{}) })
	panics(t, "minutes", func(c *ui.Context) {
		Directions(c, []Step{{Text: "Go", Turn: "arrow-up", Minutes: 100001}}, DirectionsOptions{})
	})
}

func TestSpanRounds(t *testing.T) {
	var u Units
	u.fill()
	for _, c := range []struct {
		m        float64
		imperial bool
		want     string
	}{
		{850, false, "850 m"}, {999.6, false, "1.0 km"}, {4200, false, "4.2 km"}, {9960, false, "10 km"}, {12400, false, "12 km"},
		{100, true, "328 ft"}, {4200, true, "2.6 mi"}, {16090, true, "10 mi"},
	} {
		if got := span(c.m, c.imperial, u); got != c.want {
			t.Errorf("span(%v, %v) = %q, want %q", c.m, c.imperial, got, c.want)
		}
	}
	panics(t, "Distance", func(c *ui.Context) { Distance(c, math.NaN(), DistanceOptions{}) })
	panics(t, "Distance", func(c *ui.Context) { Distance(c, -1, DistanceOptions{}) })
}

func TestDistanceBadge(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Distance(c, 4200, DistanceOptions{Minutes: 68})
		Distance(c, 850, DistanceOptions{})
	}), 400, 300)
	settles(t, tt)
	if !exact(tt, "4.2 km") || !exact(tt, "1 h 8 min") || !exact(tt, "850 m") {
		t.Errorf("badge texts %q", tt.Texts())
	}
	if exact(tt, "0 min") {
		t.Error("zero minutes showed a time")
	}
}

func TestCompassTurnsShortWay(t *testing.T) {
	heading := 350.0
	tt := ui.NewTester(padded(func(c *ui.Context) { Compass(c, heading, CompassOptions{}) }), 400, 400)
	settles(t, tt)
	if !exact(tt, "350°") || !exact(tt, "N") {
		t.Fatalf("compass texts %q", tt.Texts())
	}
	r := rectOf(t, tt, "Heading")
	size := float32(160)
	x0, y0 := r.X+(r.W-size)/2, r.Y
	heading = 10
	tt.Frame()
	red := 0
	for range 6 {
		tt.Advance(40 * 1e6)
		img := tt.Image()
		for y := int(y0 + size*0.6); y < int(y0+size); y++ {
			for x := int(x0); x < int(x0+size); x++ {
				p := img.RGBAAt(x, y)
				if p.R > 180 && p.G < 110 && p.B < 110 {
					red++
				}
			}
		}
	}
	if red > 0 {
		t.Errorf("the needle swept through the south: %d red pixels low on the dial", red)
	}
	settles(t, tt)
	if !exact(tt, "10°") {
		t.Errorf("compass texts %q", tt.Texts())
	}
	panics(t, "finite", func(c *ui.Context) { Compass(c, math.Inf(1), CompassOptions{}) })
}
