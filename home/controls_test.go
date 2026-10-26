package home

import (
	"image"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestLightSwitchDimTone(t *testing.T) {
	on, level, tone := true, 0.7, 0
	tt := ui.NewTester(padded(func(c *ui.Context) { Light(c, "Ceiling light", &on, &level, &tone, LightOptions{}) }), 420, 400)
	settles(t, tt)
	if !exact(tt, "70%") {
		t.Errorf("level: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Switch Ceiling light")
	cues(t, before, sound.Toggle)
	if on || !exact(tt, "Off") {
		t.Errorf("switch off: %v %q", on, tt.Texts())
	}
	click(t, tt, "Cool")
	if tone != 0 {
		t.Errorf("a dark light took a tone: %d", tone)
	}
	click(t, tt, "Switch Ceiling light")
	before = baseline()
	click(t, tt, "Cool")
	cues(t, before, sound.Select)
	before = baseline()
	click(t, tt, "Cool")
	cues(t, before)
	if tone != 2 {
		t.Errorf("tone %d", tone)
	}
	panics(t, "level of 0.01 to 1", func(c *ui.Context) { l := 0.0; Light(c, "Lamp", &on, &l, &tone, LightOptions{}) })
	panics(t, "tone of 0 to 2", func(c *ui.Context) { k := 3; Light(c, "Lamp", &on, &level, &k, LightOptions{}) })
}

func TestThermostatSteps(t *testing.T) {
	target, current := 21.5, 20.0
	tt := ui.NewTester(padded(func(c *ui.Context) { Thermostat(c, "Living room", &target, current, ThermostatOptions{}) }), 420, 500)
	settles(t, tt)
	if !exact(tt, "Heating") || !exact(tt, "Now 20.0°") {
		t.Errorf("mode: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Warmer")
	cues(t, before, sound.Tick)
	if target != 22 || !exact(tt, "22.0") {
		t.Errorf("warmer: %v %q", target, tt.Texts())
	}
	target = 30
	tt.Frame()
	before = baseline()
	click(t, tt, "Warmer")
	cues(t, before)
	click(t, tt, "Cooler")
	if target != 29.5 {
		t.Errorf("cooler from the top: %v", target)
	}
	target, current = 18, 21
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "Cooling") {
		t.Errorf("cooling: %q", tt.Texts())
	}
	panics(t, "must sit in 10..30", func(c *ui.Context) { v := 35.0; Thermostat(c, "Hall", &v, 20, ThermostatOptions{}) })
	panics(t, "finite range", func(c *ui.Context) { v := 20.0; Thermostat(c, "Hall", &v, 20, ThermostatOptions{Min: 30, Max: 10}) })
}

func hold(tt *ui.Tester, r ui.Rect, d time.Duration) {
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	tt.Advance(d)
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
}

func TestLockHold(t *testing.T) {
	locked := true
	tt := ui.NewTester(padded(func(c *ui.Context) { Lock(c, "Front door", &locked, LockOptions{}) }), 400, 400)
	settles(t, tt)
	before := baseline()
	hold(tt, rectOf(t, tt, "Unlock Front door"), 300*time.Millisecond)
	cues(t, before)
	if !locked {
		t.Error("a short press unlocked")
	}
	before = baseline()
	hold(tt, rectOf(t, tt, "Unlock Front door"), 2*time.Second)
	cues(t, before, sound.Success)
	if locked || !exact(tt, "Unlocked") || !exact(tt, "Hold to lock") {
		t.Errorf("hold: %v %q", locked, tt.Texts())
	}
	panics(t, "hold of 0.1 to 5 s", func(c *ui.Context) { Lock(c, "Door", &locked, LockOptions{Hold: 50 * time.Millisecond}) })
}

func TestCameraLiveAndOffline(t *testing.T) {
	var frame *ui.Bitmap
	at := time.Date(2026, 10, 25, 19, 42, 8, 0, time.UTC)
	tt := ui.NewTester(padded(func(c *ui.Context) { Camera(c, "Porch", frame, at, CameraOptions{}) }), 420, 400)
	settles(t, tt)
	if !exact(tt, "Offline") || exact(tt, "LIVE") || tt.FramePending() {
		t.Errorf("offline: %q", tt.Texts())
	}
	frame = ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 4, 4)))
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "LIVE") || !exact(tt, "19:42:08") || !tt.FramePending() {
		t.Errorf("live: %q", tt.Texts())
	}
	panics(t, "needs a name and a frame time", func(c *ui.Context) {
		Camera(c, "Porch", ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2))), time.Time{}, CameraOptions{})
	})
}

func TestSceneRunsAndHolds(t *testing.T) {
	scenes := []SceneItem{{"Morning", "sunrise"}, {"Evening", "sunset"}, {"Movie", "clapperboard"}, {"Away", "house"}}
	active := "Evening"
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Scene(c, scenes, &active, SceneOptions{}).Changed() {
			changes++
		}
	}), 420, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Movie")
	cues(t, before, sound.Tap)
	before = baseline()
	click(t, tt, "Movie")
	cues(t, before, sound.Tap)
	if active != "Movie" || changes != 1 {
		t.Errorf("run: %q %d", active, changes)
	}
	swapUnderPress(t, tt, "Morning", func() { scenes[0], scenes[1] = scenes[1], scenes[0] })
	if active == "Evening" {
		t.Error("a reorder under the press ran the scene that slid in")
	}
	panics(t, "blank or repeated", func(c *ui.Context) {
		a := ""
		Scene(c, []SceneItem{{"A", "sun"}, {"A", "moon"}}, &a, SceneOptions{})
	})
	panics(t, "active name among them", func(c *ui.Context) { a := "Nope"; Scene(c, scenes, &a, SceneOptions{}) })
}

func TestRoomsPick(t *testing.T) {
	rooms := []Room{{"Living room", "sofa", 3}, {"Kitchen", "utensils-crossed", 1}, {"Bedroom", "bed-double", 0}}
	picked := "Living room"
	tt := ui.NewTester(padded(func(c *ui.Context) { Rooms(c, rooms, &picked, RoomsOptions{}) }), 600, 200)
	settles(t, tt)
	if !exact(tt, "3 on") || exact(tt, "0 on") {
		t.Errorf("counts: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Kitchen")
	cues(t, before, sound.Select)
	if picked != "Kitchen" {
		t.Errorf("picked %q", picked)
	}
	panics(t, "on count of 0 or more", func(c *ui.Context) { p := "A"; Rooms(c, []Room{{"A", "sofa", -1}}, &p, RoomsOptions{}) })
}

func TestDeviceToggles(t *testing.T) {
	on := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Device(c, "Coffee maker", "coffee", &on, DeviceOptions{Room: "Kitchen"}) }), 300, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Switch Coffee maker")
	cues(t, before, sound.Toggle)
	settles(t, tt)
	if !on || !exact(tt, "On") {
		t.Errorf("device: %v %q", on, tt.Texts())
	}
	panics(t, "needs a name, an icon", func(c *ui.Context) { Device(c, "Fridge", "", &on, DeviceOptions{}) })
}

func TestFanSpeeds(t *testing.T) {
	speed := 2
	tt := ui.NewTester(padded(func(c *ui.Context) { Fan(c, "Ceiling fan", &speed, FanOptions{}) }), 420, 300)
	settles(t, tt)
	if !tt.FramePending() {
		t.Error("a running fan stands still")
	}
	before := baseline()
	click(t, tt, "Off")
	cues(t, before, sound.Tick)
	settles(t, tt)
	if speed != 0 || tt.FramePending() {
		t.Errorf("off: %d pending %v", speed, tt.FramePending())
	}
	// The state text reads Off too; press the button.
	low := rectOf(t, tt, "Low")
	x, y := low.X-low.W/2-4, low.Y+low.H/2
	before = baseline()
	tt.Press(x, y)
	tt.Release(x, y)
	tt.Frame()
	cues(t, before)
	if speed != 0 {
		t.Errorf("speed %d", speed)
	}
	panics(t, "speed of 0 to 3", func(c *ui.Context) { s := 4; Fan(c, "Fan", &s, FanOptions{}) })
}

func TestBlindsSlider(t *testing.T) {
	open := 0.35
	tt := ui.NewTester(padded(func(c *ui.Context) { Blinds(c, "Bay window", &open, BlindsOptions{}) }), 420, 300)
	settles(t, tt)
	if !exact(tt, "35% open") {
		t.Errorf("level: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Position")
	cues(t, before, sound.Tick)
	if open != 0.5 || !exact(tt, "50% open") {
		t.Errorf("slider: %v %q", open, tt.Texts())
	}
	for v, want := range map[float64]string{0: "Closed", 1: "Open"} {
		open = v
		tt.Frame()
		if !exact(tt, want) {
			t.Errorf("%v: %q", v, tt.Texts())
		}
	}
	panics(t, "open share of 0 to 1", func(c *ui.Context) { v := 1.5; Blinds(c, "Window", &v, BlindsOptions{}) })
}

func TestEnergyHover(t *testing.T) {
	use := []float64{1, 2, 3, 0.5}
	labels := []string{"06:00", "07:00", "08:00", "09:00"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Energy(c, use, EnergyOptions{Labels: labels}) }), 420, 400)
	settles(t, tt)
	if !exact(tt, "6.5") || !exact(tt, "Peak 08:00") {
		t.Errorf("total: %q", tt.Texts())
	}
	r := rectOf(t, tt, "06:00")
	tt.Move(r.X+2, r.Y-12)
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "1.0") {
		t.Errorf("hover: %q", tt.Texts())
	}
	panics(t, "each with a label", func(c *ui.Context) { Energy(c, use, EnergyOptions{Labels: labels[:2]}) })
	panics(t, "must be 0 to 1e6", func(c *ui.Context) { Energy(c, []float64{1, -1}, EnergyOptions{Labels: labels[:2]}) })
}
