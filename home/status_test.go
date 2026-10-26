package home

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPlugDraw(t *testing.T) {
	on := true
	tt := ui.NewTester(padded(func(c *ui.Context) { Plug(c, "Desk lamp plug", &on, 42.6, PlugOptions{Today: 0.38}) }), 420, 300)
	settles(t, tt)
	if !exact(tt, "42.6 W") || !exact(tt, "0.38 kWh today") {
		t.Errorf("draw: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Switch Desk lamp plug")
	cues(t, before, sound.Toggle)
	if on || !exact(tt, "0.0 W") {
		t.Errorf("off: %v %q", on, tt.Texts())
	}
	panics(t, "0 to 1e5 watts", func(c *ui.Context) { Plug(c, "Plug", &on, -1, PlugOptions{}) })
}

func TestSpeakerPlayAndVolume(t *testing.T) {
	playing, volume := true, 0.3
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Speaker(c, "Kitchen speaker", Track{Title: "Holocene", Artist: "Bon Iver"}, &playing, &volume, SpeakerOptions{})
	}), 420, 300)
	settles(t, tt)
	if !tt.FramePending() {
		t.Error("a playing speaker shows no beat")
	}
	before := baseline()
	click(t, tt, "Pause")
	cues(t, before, sound.Tap)
	settles(t, tt)
	if playing || tt.FramePending() {
		t.Errorf("pause: %v pending %v", playing, tt.FramePending())
	}
	if _, ok := tt.Find("Play"); !ok {
		t.Errorf("no Play after pause: %q", tt.Texts())
	}
	before = baseline()
	click(t, tt, "Volume")
	cues(t, before, sound.Tick)
	if volume != 0.5 {
		t.Errorf("volume %v", volume)
	}
	panics(t, "volume of 0 to 1", func(c *ui.Context) { v := 2.0; Speaker(c, "S", Track{Title: "T"}, &playing, &v, SpeakerOptions{}) })
}

func TestAirQualityBands(t *testing.T) {
	for aqi, want := range map[int]string{0: "Good", 50: "Good", 51: "Moderate", 100: "Moderate", 101: "Poor", 200: "Poor", 201: "Hazardous", 500: "Hazardous"} {
		tt := ui.NewTester(padded(func(c *ui.Context) { AirQuality(c, aqi, AirQualityOptions{}) }), 420, 400)
		settles(t, tt)
		if !exact(tt, want) {
			t.Errorf("%d: %q", aqi, tt.Texts())
		}
	}
	panics(t, "index of 0 to 500", func(c *ui.Context) { AirQuality(c, 501, AirQualityOptions{}) })
	panics(t, "blank or repeated", func(c *ui.Context) {
		AirQuality(c, 10, AirQualityOptions{Pollutants: []Pollutant{{"PM10", "1", "x"}, {"PM10", "2", "x"}}})
	})
}

func TestHumidityBands(t *testing.T) {
	for v, want := range map[float64]string{39.9: "Dry", 40: "Comfortable", 60: "Comfortable", 60.1: "Humid"} {
		tt := ui.NewTester(padded(func(c *ui.Context) { Humidity(c, v, HumidityOptions{}) }), 420, 400)
		settles(t, tt)
		if !exact(tt, want) {
			t.Errorf("%v: %q", v, tt.Texts())
		}
	}
	panics(t, "band inside it", func(c *ui.Context) { Humidity(c, 50, HumidityOptions{Low: 60, High: 40}) })
}

func TestSensorsSummary(t *testing.T) {
	sensors := []Sensor{{Name: "Front door", Open: true}, {Name: "Kitchen window", Kind: WindowSensor, Open: true}, {Name: "Back door"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Sensors(c, sensors, SensorsOptions{}) }), 420, 400)
	settles(t, tt)
	if !exact(tt, "2 open") {
		t.Errorf("open: %q", tt.Texts())
	}
	sensors = []Sensor{{Name: "Back door"}}
	tt.Frame()
	if !exact(tt, "All closed") {
		t.Errorf("closed: %q", tt.Texts())
	}
	sensors = nil
	tt.Frame()
	if !exact(tt, "No sensors yet") || exact(tt, "All closed") {
		t.Errorf("empty: %q", tt.Texts())
	}
	panics(t, "unknown kind", func(c *ui.Context) { Sensors(c, []Sensor{{Name: "Gate", Kind: 7}}, SensorsOptions{}) })
	panics(t, "blank or repeated", func(c *ui.Context) { Sensors(c, []Sensor{{Name: "A"}, {Name: "A"}}, SensorsOptions{}) })
}

func TestAutomationSwitch(t *testing.T) {
	on := true
	r := Rule{Name: "Sunset lights", When: "The sun sets", Then: "Lights on", WhenIcon: "sunset", ThenIcon: "lightbulb"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Automation(c, r, &on, AutomationOptions{}) }), 420, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Switch Sunset lights")
	cues(t, before, sound.Toggle)
	if on {
		t.Error("rule still on")
	}
	panics(t, "their icons", func(c *ui.Context) { Automation(c, Rule{Name: "R", When: "w", Then: "t"}, &on, AutomationOptions{}) })
}

func TestVacuumRun(t *testing.T) {
	state := Docked
	tt := ui.NewTester(padded(func(c *ui.Context) { Vacuum(c, "Robo vacuum", &state, 0.86, VacuumOptions{}) }), 420, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Dock")
	cues(t, before)
	before = baseline()
	click(t, tt, "Start")
	cues(t, before, sound.Tap)
	settles(t, tt)
	if state != Cleaning || !exact(tt, "Cleaning") || !tt.FramePending() {
		t.Errorf("start: %v %q", state, tt.Texts())
	}
	click(t, tt, "Dock")
	if state != Returning {
		t.Errorf("dock: %v", state)
	}
	click(t, tt, "Start")
	if state != Cleaning {
		t.Errorf("start from returning: %v", state)
	}
	panics(t, "battery of 0 to 1", func(c *ui.Context) { s := Docked; Vacuum(c, "V", &s, 2, VacuumOptions{}) })
}

func TestGarageMoving(t *testing.T) {
	open := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Garage(c, "Garage", &open, GarageOptions{}) }), 420, 500)
	settles(t, tt)
	if !exact(tt, "Closed") {
		t.Errorf("closed: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Open Garage")
	cues(t, before, sound.Toggle)
	tt.Advance(400 * time.Millisecond)
	if !open || !exact(tt, "Opening…") {
		t.Errorf("opening: %v %q", open, tt.Texts())
	}
	tt.Advance(3 * time.Second)
	if !exact(tt, "Open") {
		t.Errorf("open: %q", tt.Texts())
	}
	click(t, tt, "Close Garage")
	tt.Advance(400 * time.Millisecond)
	if !exact(tt, "Closing…") {
		t.Errorf("closing: %q", tt.Texts())
	}
}

func TestSprinklersRun(t *testing.T) {
	zones := []Zone{{"Front lawn", 15}, {"Back garden", 20}, {"Hedges", 8}}
	running := ""
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Sprinklers(c, zones, &running, SprinklersOptions{}).Changed() {
			changes++
		}
	}), 420, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Run Front lawn")
	cues(t, before, sound.Toggle)
	if running != "Front lawn" || !exact(tt, "15:00 left") {
		t.Errorf("run: %q %q", running, tt.Texts())
	}
	tt.Advance(61500 * time.Millisecond)
	if !exact(tt, "13:59 left") {
		t.Errorf("countdown: %q", tt.Texts())
	}
	click(t, tt, "Run Hedges")
	if running != "Hedges" || !exact(tt, "8:00 left") || exact(tt, "13:59 left") {
		t.Errorf("switch: %q %q", running, tt.Texts())
	}
	n := changes
	tt.Advance(8*time.Minute + time.Second)
	if running != "" || changes != n+1 || !exact(tt, "8 min") {
		t.Errorf("finish: %q %d %q", running, changes-n, tt.Texts())
	}
	click(t, tt, "Run Back garden")
	click(t, tt, "Stop Back garden")
	if running != "" {
		t.Errorf("stop: %q", running)
	}
	panics(t, "1 to 120 minutes", func(c *ui.Context) { r := ""; Sprinklers(c, []Zone{{"A", 0}}, &r, SprinklersOptions{}) })
	panics(t, "running name among them", func(c *ui.Context) { r := "B"; Sprinklers(c, []Zone{{"A", 5}}, &r, SprinklersOptions{}) })
}

func TestHomeSummaryStatus(t *testing.T) {
	s := HomeStatus{Name: "14 Alder Street", Inside: 21.46, Energy: 8.2, Lights: 4, On: 7, Locked: true}
	tt := ui.NewTester(padded(func(c *ui.Context) { HomeSummary(c, s, HomeSummaryOptions{}) }), 420, 400)
	settles(t, tt)
	for _, want := range []string{"All doors locked", "21.5°", "8.2 kWh", "4", "7"} {
		if !exact(tt, want) {
			t.Errorf("missing %q: %q", want, tt.Texts())
		}
	}
	s.Locked = false
	tt.Frame()
	if !exact(tt, "A door is unlocked") {
		t.Errorf("unlocked: %q", tt.Texts())
	}
	s.Alerts = 2
	tt.Frame()
	if !exact(tt, "2 alerts") {
		t.Errorf("alerts: %q", tt.Texts())
	}
	panics(t, "-60 to 60 degrees", func(c *ui.Context) { HomeSummary(c, HomeStatus{Name: "Home", Inside: 99}, HomeSummaryOptions{}) })
}
