package gallery

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestHomeDemosLog(t *testing.T) {
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
	hold := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			r, ok := tt.Find(name)
			if !ok {
				return fmt.Errorf("%q missing", name)
			}
			tt.Press(r.X+r.W/2, r.Y+r.H/2)
			tt.Frame()
			tt.Advance(800 * time.Millisecond)
			tt.Release(r.X+r.W/2, r.Y+r.H/2)
			tt.Frame()
			return nil
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Light":      {press("Switch Ceiling light"), "on false, 70%, tone 0"},
		"Thermostat": {press("Warmer"), "target 22.0"},
		"Lock":       {hold("Unlock Front door"), "locked false"},
		"Scene":      {press("Movie night", "Movie night"), "runs 1, scene Movie night\nruns 2, scene Movie night"},
		"Rooms":      {press("Kitchen"), "room Kitchen"},
		"Device":     {press("Switch Coffee maker"), "coffee maker on true"},
		"Fan":        {press("High"), "speed 3"},
		"Blinds":     {press("Position"), "open 50%"},
		"Plug":       {press("Switch Desk lamp plug"), "plug false"},
		"Speaker":    {press("Pause"), "playing false, volume 45%"},
		"Automation": {press("Switch Sunset lights"), "rule false"},
		"Vacuum":     {press("Start"), "vacuum 1"},
		"Garage":     {press("Open Garage"), "door true"},
		"Sprinklers": {press("Run Front lawn"), "running Front lawn"},
	}
	quiet := map[string]bool{"Camera": true, "Energy": true, "AirQuality": true, "Humidity": true, "Sensors": true, "HomeSummary": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "home" {
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
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
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
	if seen != 20 {
		t.Errorf("checked %d home demos", seen)
	}
}
