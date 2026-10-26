package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestGeoDemosLog(t *testing.T) {
	press := func(names ...string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, name := range names {
				if err := tt.Click(name); err != nil {
					return err
				}
				tt.Frame()
			}
			return nil
		}
	}
	// slide presses the Geofence slider halfway between its end labels.
	slide := func(tt *ui.Tester) error {
		lo, ok1 := tt.Find("100 m")
		hi, ok2 := tt.Find("2.0 km")
		if !ok1 || !ok2 {
			t.Fatal("no Geofence slider labels")
		}
		x, y := (lo.X+lo.W+hi.X)/2, lo.Y+lo.H/2
		tt.Press(x, y)
		tt.Release(x, y)
		return nil
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"MapCanvas":       {press("Zoom in"), "zoom 15"},
		"Pin":             {press("Clay supply"), "dropped Clay supply"},
		"LocationSearch":  {press("Mission Clay Supply"), "picked Mission Clay Supply"},
		"Compass":         {press("Turn 90°"), "heading 22"},
		"Coordinates":     {press("Copy"), "copied"},
		"RegionSelect":    {press("Oregon"), "region Oregon"},
		"StoreLocator":    {press("Riverbend Dogpatch"), "picked Riverbend Dogpatch"},
		"DeliveryTracker": {press("Advance"), "progress 60%"},
		"Geofence":        {slide, "radius "},
		"AddressLookup":   {press("2125 Bryant St, San Francisco, CA 94110"), "chose 2125 Bryant St"},
	}
	quiet := map[string]bool{"Cluster": true, "PlaceCard": true, "Directions": true, "Distance": true, "Elevation": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "geo" {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 900)
		tt.Advance(2 * time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 15 {
		t.Errorf("checked %d geo demos", seen)
	}
}

func TestGeoCoordinatesLogsNoText(t *testing.T) {
	for _, e := range Entries() {
		if e.Category != "geo" || e.Name != "Coordinates" {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 900)
		tt.Advance(2 * time.Second)
		tt.Key(0, ui.KeyTab)
		tt.Frame()
		tt.Key(0, ui.KeyEnd)
		tt.Type("54321")
		tt.Frame()
		got := strings.Join(logs, "\n")
		if !strings.Contains(got, "coordinates edited") || strings.Contains(got, "54321") || strings.Contains(got, "37.") {
			t.Errorf("Coordinates logs %q", got)
		}
	}
}
