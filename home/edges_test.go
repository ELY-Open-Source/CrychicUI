package home

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestLockOnePressOneTurn(t *testing.T) {
	locked := true
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Lock(c, "Front door", &locked, LockOptions{}).Changed() {
			changes++
		}
	}), 400, 400)
	settles(t, tt)
	r := rectOf(t, tt, "Unlock Front door")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	tt.Advance(time.Second)
	tt.Move(5, 5)
	tt.Frame()
	tt.Move(x, y)
	tt.Frame()
	tt.Advance(2 * time.Second)
	tt.Release(x, y)
	tt.Frame()
	if locked || changes != 1 {
		t.Errorf("one press turned the lock %d times, locked %v", changes, locked)
	}
	tt.Press(x, y)
	tt.Frame()
	tt.Advance(200 * time.Millisecond)
	tt.Move(5, 5)
	tt.Frame()
	tt.Move(x, y)
	tt.Frame()
	tt.Advance(2 * time.Second)
	tt.Release(x, y)
	tt.Frame()
	if locked || changes != 1 {
		t.Errorf("a hold that left the button turned it: %d, locked %v", changes, locked)
	}
}

func TestSprinklersExternalStop(t *testing.T) {
	zones := []Zone{{"Front lawn", 1}, {"Hedges", 8}}
	running := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Sprinklers(c, zones, &running, SprinklersOptions{}) }), 420, 500)
	settles(t, tt)
	click(t, tt, "Run Front lawn")
	tt.Advance(30 * time.Second)
	running = ""
	tt.Frame()
	tt.Advance(70 * time.Second)
	running = "Front lawn"
	tt.Frame()
	if running != "Front lawn" || !exact(tt, "1:00 left") {
		t.Errorf("restart after an outside stop: %q %q", running, tt.Texts())
	}
}

func TestGarageCopyBesideButton(t *testing.T) {
	open := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Garage(c, "Garage", &open, GarageOptions{Closed: "door"}) }), 420, 500)
	settles(t, tt)
	click(t, tt, "Open Garage")
	if !open {
		t.Error("door did not open")
	}
}

func TestThermostatCustomRange(t *testing.T) {
	target := 10.2
	changes := 0
	o := ThermostatOptions{Min: 10.2, Max: 11.2, Step: 0.5}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Thermostat(c, "Hall", &target, 20, o).Changed() {
			changes++
		}
	}), 420, 500)
	settles(t, tt)
	click(t, tt, "Warmer")
	click(t, tt, "Warmer")
	if math.Abs(target-11.2) > 1e-9 {
		t.Errorf("warmer stops at %v, not the top", target)
	}
	click(t, tt, "Cooler")
	click(t, tt, "Cooler")
	if math.Abs(target-10.2) > 1e-9 || changes != 4 {
		t.Errorf("cooler stops at %v after %d changes", target, changes)
	}
	before := baseline()
	click(t, tt, "Cooler")
	cues(t, before)
	if changes != 4 {
		t.Errorf("a step past the bottom reported: %d", changes)
	}
	for _, step := range []float64{math.Inf(1), 1e-16, math.SmallestNonzeroFloat64} {
		panics(t, "finite range with a step", func(c *ui.Context) {
			v := 20.0
			Thermostat(c, "Hall", &v, 20, ThermostatOptions{Step: step})
		})
	}
}

func TestEnergyLongLabelsFit(t *testing.T) {
	labels := []string{"Monday, October 19, 2026 at 06:00", "Tuesday, October 20, 2026 at 06:00", "Wednesday, October 21, 2026 at 06:00"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Energy(c, []float64{4, 6, 5}, EnergyOptions{Labels: labels}) }), 320, 400)
	settles(t, tt)
	var right float32
	for _, l := range labels {
		r := rectOf(t, tt, l)
		inside(t, l, r, right, 310)
		right = r.X + r.W
	}
}

func TestSceneRunsFull(t *testing.T) {
	runs, active := math.MaxInt-1, ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Scene(c, []SceneItem{{"Evening", "sunset"}, {"Away", "house"}}, &active, SceneOptions{Runs: &runs})
	}), 420, 400)
	settles(t, tt)
	click(t, tt, "Evening")
	settles(t, tt)
	if runs != math.MaxInt || active != "Evening" {
		t.Errorf("last run: %d %q", runs, active)
	}
	before := baseline()
	click(t, tt, "Away")
	cues(t, before)
	if runs != math.MaxInt || active != "Evening" {
		t.Errorf("a full count still ran: %d %q", runs, active)
	}
	panics(t, "run count of 0 or more", func(c *ui.Context) {
		r, a := -1, ""
		Scene(c, []SceneItem{{"Evening", "sunset"}}, &a, SceneOptions{Runs: &r})
	})
}

func TestEnergyLongPeakFits(t *testing.T) {
	labels := []string{"Wednesday, October 21, 2026, 06:00–07:00 America/New_York", "Thursday, October 22, 2026, 06:00–07:00 America/New_York"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Energy(c, []float64{4, 6}, EnergyOptions{Labels: labels}) }), 320, 400)
	settles(t, tt)
	title := rectOf(t, tt, "Energy today")
	chip := rectOf(t, tt, "Peak "+labels[1])
	inside(t, "title", title, 10, chip.X)
	inside(t, "peak", chip, title.X+title.W, 310)
	if title.W < 40 {
		t.Errorf("title squeezed to %.1f", title.W)
	}
}

func TestLightNilNamesCaller(t *testing.T) {
	level, tone := 0.5, 0
	panics(t, `"Ceiling light"`, func(c *ui.Context) { Light(c, "Ceiling light", nil, &level, &tone, LightOptions{}) })
}
