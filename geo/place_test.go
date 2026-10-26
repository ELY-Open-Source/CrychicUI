package geo

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCoordinatesEdits(t *testing.T) {
	at := studio
	changes := 0
	var copied bool
	tt := ui.NewTester(padded(func(c *ui.Context) {
		r, copy := Coordinates(c, &at, CoordinatesOptions{})
		if r.Changed() {
			changes++
		}
		copied = copied || copy.Clicked()
	}), 520, 500)
	settles(t, tt)
	if !exact(tt, "37°45′36″N") || !exact(tt, "122°24′53″W") {
		t.Fatalf("coordinate texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Type("95")
	tt.Frame()
	if at != studio || !exact(tt, "Enter a number from -85 to 85") {
		t.Errorf("an out-of-range latitude changed the point to %v: %q", at, tt.Texts())
	}
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Type("-33.8688")
	tt.Frame()
	settles(t, tt)
	if at.Lat != -33.8688 || at.Lng != studio.Lng || !exact(tt, "33°52′08″S") {
		t.Errorf("valid latitude gave %v, %q", at, tt.Texts())
	}
	if exact(tt, "Enter a number from -85 to 85") || changes == 0 {
		t.Errorf("error stayed or no change reported (%d)", changes)
	}
	before := baseline()
	click(t, tt, "Copy")
	cues(t, before, sound.Pop)
	if tt.Clipboard() != "-33.8688,-122.4148" || !exact(tt, "Copied") || !copied {
		t.Errorf("clipboard %q, texts %q", tt.Clipboard(), tt.Texts())
	}
	// An outside point rewrites the fields and clears Copied.
	at = Point{Lat: 51.5074, Lng: -0.1278}
	tt.Frame()
	tt.Frame()
	if exact(tt, "Copied") || !exact(tt, "51°30′27″N") {
		t.Errorf("outside change left %q", tt.Texts())
	}
	// The field holds the new latitude; a digit extends it.
	lat := rectOf(t, tt, "Latitude")
	tt.Press(lat.X+20, lat.Y+lat.H+18)
	tt.Release(lat.X+20, lat.Y+lat.H+18)
	tt.Frame()
	tt.Key(0, ui.KeyEnd)
	tt.Type("1")
	tt.Frame()
	if at.Lat != 51.50741 {
		t.Errorf("typing after an outside change gave latitude %v", at.Lat)
	}
}

func regionCountries() []Country {
	return []Country{
		{Name: "United States", Regions: []string{"California", "Oregon", "Washington"}},
		{Name: "Canada", Regions: []string{"British Columbia", "Ontario"}},
	}
}

func TestRegionSelectPicks(t *testing.T) {
	country, region := "United States", "California"
	countries := regionCountries()
	tt := ui.NewTester(padded(func(c *ui.Context) {
		RegionSelect(c, countries, &country, &region, RegionSelectOptions{})
	}), 520, 600)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Oregon")
	if region != "Oregon" {
		t.Fatalf("region %q", region)
	}
	cues(t, before, sound.Select)
	before = baseline()
	click(t, tt, "Oregon")
	cues(t, before)
	before = baseline()
	click(t, tt, "United States")
	settles(t, tt)
	click(t, tt, "Canada")
	if country != "Canada" || region != "" {
		t.Errorf("country %q region %q after switching", country, region)
	}
	cues(t, before, sound.Open, sound.Select)
	settles(t, tt)
	if !exact(tt, "Ontario") || exact(tt, "Oregon") {
		t.Errorf("regions did not follow the country: %q", tt.Texts())
	}
}

func TestRegionSelectRejects(t *testing.T) {
	country, region := "United States", "Ontario"
	panics(t, "not in", func(c *ui.Context) { RegionSelect(c, regionCountries(), &country, &region, RegionSelectOptions{}) })
	country, region = "Atlantis", ""
	panics(t, "listed country", func(c *ui.Context) { RegionSelect(c, regionCountries(), &country, &region, RegionSelectOptions{}) })
	country = "Canada"
	panics(t, "repeated", func(c *ui.Context) {
		RegionSelect(c, []Country{{Name: "Canada", Regions: []string{"Ontario", "Ontario"}}}, &country, &region, RegionSelectOptions{})
	})
}

func TestRegionSelectKeepsIdentity(t *testing.T) {
	country, region := "United States", ""
	countries := regionCountries()
	tt := ui.NewTester(padded(func(c *ui.Context) {
		RegionSelect(c, countries, &country, &region, RegionSelectOptions{})
	}), 520, 600)
	settles(t, tt)
	regions := countries[0].Regions
	swapUnderPress(t, tt, "Oregon", func() { regions[0], regions[1] = regions[1], regions[0] })
	if region == "California" {
		t.Error("the region that slid under the pointer was picked")
	}
}

func locatorStores() []Store {
	return []Store{
		{Name: "Riverbend Mission", Address: "2125 Bryant St", At: studio, Open: true},
		{Name: "Riverbend Dogpatch", Address: "901 Minnesota St", At: Point{Lat: 37.7576, Lng: -122.3880}},
	}
}

func TestStoreLocatorPicks(t *testing.T) {
	picked := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		StoreLocator(c, locatorStores(), &picked, StoreLocatorOptions{From: &studio})
	}), 520, 700)
	settles(t, tt)
	if !exact(tt, "0 m") || !exact(tt, "Closed") {
		t.Errorf("store texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Riverbend Dogpatch")
	if picked != "Riverbend Dogpatch" {
		t.Fatalf("picked %q", picked)
	}
	cues(t, before, sound.Select)
	if !tt.FramePending() {
		t.Error("the map did not move to the picked store")
	}
	before = baseline()
	settles(t, tt)
	click(t, tt, "Riverbend Dogpatch")
	cues(t, before)
	picked = "Riverbend Oakland"
	panics(t, "no store", func(c *ui.Context) { StoreLocator(c, locatorStores(), &picked, StoreLocatorOptions{}) })
}

func TestDeliveryTrackerProgress(t *testing.T) {
	progress := 0.3
	route := []Point{{Lat: 37.7650, Lng: -122.4190}, {Lat: 37.7625, Lng: -122.4150}, studio}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		DeliveryTracker(c, route, progress, DeliveryTrackerOptions{Courier: "Lucas Martin", Minutes: 10})
	}), 520, 600)
	settles(t, tt)
	if !exact(tt, "Arriving in 10 min") {
		t.Fatalf("tracker texts %q", tt.Texts())
	}
	progress = 1
	tt.Frame()
	tt.Advance(100 * 1e6)
	mid := tt.Image()
	settles(t, tt)
	if !exact(tt, "Delivered") {
		t.Errorf("finished tracker texts %q", tt.Texts())
	}
	// Mid-glide, the map differs from where the courier ends.
	end := tt.Image()
	same := true
	for y := 12; y < 180 && same; y++ {
		for x := 12; x < 500; x++ {
			if mid.RGBAAt(x, y) != end.RGBAAt(x, y) {
				same = false
				break
			}
		}
	}
	if same {
		t.Error("the courier jumped instead of moving")
	}
	panics(t, "progress", func(c *ui.Context) {
		DeliveryTracker(c, route, 1.5, DeliveryTrackerOptions{Courier: "Lucas Martin"})
	})
}

func TestGeofenceSlides(t *testing.T) {
	radius := 500.0
	changes := 0
	tt := ui.NewTester(padded(within(420, func(c *ui.Context) {
		if Geofence(c, studio, &radius, GeofenceOptions{}).Changed() {
			changes++
		}
	})), 520, 600)
	settles(t, tt)
	lo := rectOf(t, tt, "100 m")
	hi := rectOf(t, tt, "2.0 km")
	before := baseline()
	x, y := (lo.X+lo.W+hi.X)/2, lo.Y+lo.H/2
	tt.Press(x, y)
	tt.Release(x, y)
	tt.Frame()
	if radius <= 900 || radius >= 1200 || changes == 0 {
		t.Errorf("halfway press gave radius %v, %d changes", radius, changes)
	}
	cues(t, before, sound.Tick)
	var u Units
	u.fill()
	if !exact(tt, span(radius, false, u)) {
		t.Errorf("radius text missing: %q", tt.Texts())
	}
	panics(t, "radius", func(c *ui.Context) { r := 50.0; Geofence(c, studio, &r, GeofenceOptions{}) })
	panics(t, "steps", func(c *ui.Context) { r := 500.0; Geofence(c, studio, &r, GeofenceOptions{Step: 0.01}) })
}

func TestElevationTotalsAndHover(t *testing.T) {
	samples := []float64{100, 150, 120, 180}
	tt := ui.NewTester(padded(within(420, func(c *ui.Context) {
		Elevation(c, samples, 5400, ElevationOptions{})
	})), 520, 600)
	settles(t, tt)
	for _, s := range []string{"+110 m", "−30 m", "High 180 m", "Low 100 m", "5.4 km"} {
		if !exact(tt, s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	read := rectOf(t, tt, "High 180 m")
	head := rectOf(t, tt, "Elevation")
	// The chart sits between the heading and the readout.
	y := (head.Y + head.H + read.Y) / 2
	tt.Move(read.X+(420-32)*1.5/4, y)
	tt.Frame()
	if !exact(tt, "150 m at 1.8 km") {
		t.Errorf("hovering the second sample shows %q", tt.Texts())
	}
	panics(t, "samples", func(c *ui.Context) { Elevation(c, []float64{1}, 100, ElevationOptions{}) })
	panics(t, "heights", func(c *ui.Context) { Elevation(c, []float64{1, 10000}, 100, ElevationOptions{}) })
}

func TestAddressLookupSuggests(t *testing.T) {
	query, chosen := "", ""
	addresses := []Address{
		{Line: "2125 Bryant St", City: "San Francisco, CA 94110"},
		{Line: "2140 Folsom St", City: "San Francisco, CA 94110"},
		{Line: "501 Dolores St", City: "San Francisco, CA 94114"},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		AddressLookup(c, &query, addresses, &chosen, AddressLookupOptions{})
	}), 520, 600)
	settles(t, tt)
	click(t, tt, "Start typing an address")
	tt.Type("21")
	tt.Frame()
	settles(t, tt)
	if exact(tt, "501 Dolores St") || !exact(tt, "San Francisco, CA 94110") {
		t.Fatalf("suggestions %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "2140 Folsom St, San Francisco, CA 94110")
	if chosen != "2140 Folsom St, San Francisco, CA 94110" || query != chosen {
		t.Fatalf("chosen %q query %q", chosen, query)
	}
	cues(t, before, sound.Select)
	click(t, tt, "Start typing an address")
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Type("zzz")
	tt.Frame()
	settles(t, tt)
	if chosen != "" || !exact(tt, "No matching addresses") {
		t.Errorf("editing kept %q, texts %q", chosen, tt.Texts())
	}
	panics(t, "repeated", func(c *ui.Context) {
		AddressLookup(c, &query, []Address{addresses[0], addresses[0]}, &chosen, AddressLookupOptions{})
	})
}
