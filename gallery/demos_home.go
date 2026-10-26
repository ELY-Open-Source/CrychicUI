package gallery

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/home"
	"github.com/egoist/mygo/ui"
)

var homeEntries = []Entry{
	{ID: 1371, Name: "Light", Category: "home", Summary: "Light with brightness and tone; cue Toggle.", Demo: demoHomeLight},
	{ID: 1372, Name: "Thermostat", Category: "home", Summary: "Dial for target temperature; cue Tick.", Demo: demoHomeThermostat},
	{ID: 1373, Name: "Lock", Category: "home", Summary: "Hold to lock or unlock; cue Success.", Demo: demoHomeLock},
	{ID: 1374, Name: "Camera", Category: "home", Summary: "Camera tile with live dot.", Demo: demoHomeCamera},
	{ID: 1375, Name: "Scene", Category: "home", Summary: "Scene buttons; cue Tap.", Demo: demoHomeScene},
	{ID: 1376, Name: "Rooms", Category: "home", Summary: "Room tabs; cue Select.", Demo: demoHomeRooms},
	{ID: 1377, Name: "Device", Category: "home", Summary: "Device tile with state; cue Toggle.", Demo: demoHomeDevice},
	{ID: 1378, Name: "Fan", Category: "home", Summary: "Fan speed with spinning blades; cue Tick.", Demo: demoHomeFan},
	{ID: 1379, Name: "Blinds", Category: "home", Summary: "Blinds position slider; cue Tick.", Demo: demoHomeBlinds},
	{ID: 1380, Name: "Energy", Category: "home", Summary: "Energy use chart.", Demo: demoHomeEnergy},
	{ID: 1381, Name: "Plug", Category: "home", Summary: "Plug toggle with power draw; cue Toggle.", Demo: demoHomePlug},
	{ID: 1382, Name: "Speaker", Category: "home", Summary: "Speaker with volume; cue Tick.", Demo: demoHomeSpeaker},
	{ID: 1383, Name: "AirQuality", Category: "home", Summary: "Air quality index ring.", Demo: demoHomeAirQuality},
	{ID: 1384, Name: "Humidity", Category: "home", Summary: "Humidity gauge.", Demo: demoHomeHumidity},
	{ID: 1385, Name: "Sensors", Category: "home", Summary: "Door and window sensors.", Demo: demoHomeSensors},
	{ID: 1386, Name: "Automation", Category: "home", Summary: "When and then rule; cue Toggle.", Demo: demoHomeAutomation},
	{ID: 1387, Name: "Vacuum", Category: "home", Summary: "Vacuum status and dock; cue Tap.", Demo: demoHomeVacuum},
	{ID: 1388, Name: "Garage", Category: "home", Summary: "Garage door with state; cue Toggle.", Demo: demoHomeGarage},
	{ID: 1389, Name: "Sprinklers", Category: "home", Summary: "Zones with run; cue Toggle.", Demo: demoHomeSprinklers},
	{ID: 1390, Name: "HomeSummary", Category: "home", Summary: "Home status overview.", Demo: demoHomeSummary},
}

func demoHomeLight(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	on := boolState(root, "on", true)
	level := numberState(root, "level", 0.7)
	tone := intState(root, "tone", 0)
	root.Children(func() {
		if home.Light(c, "Ceiling light", on, level, tone, home.LightOptions{Room: "Living room"}).Changed() {
			log(fmt.Sprintf("on %v, %d%%, tone %d", *on, int(math.Round(*level*100)), *tone))
		}
	})
}

func demoHomeThermostat(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340)
	target := numberState(root, "target", 21.5)
	root.Children(func() {
		if home.Thermostat(c, "Living room", target, 20, home.ThermostatOptions{}).Changed() {
			log("target " + strconv.FormatFloat(*target, 'f', 1, 64))
		}
	})
}

func demoHomeLock(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(320)
	locked := boolState(root, "locked", true)
	root.Children(func() {
		if home.Lock(c, "Front door", locked, home.LockOptions{}).Changed() {
			log(fmt.Sprintf("locked %v", *locked))
		}
	})
}

func demoHomeCamera(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Gap(12).Children(func() {
		home.Camera(c, "Front porch", photos()[1], time.Date(2026, 10, 25, 19, 42, 8, 0, time.UTC), home.CameraOptions{Room: "Entrance"})
	})
}

func demoHomeScene(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	active := textState(root, "active", "Evening")
	runs := intState(root, "runs", 0)
	root.Children(func() {
		scenes := []home.SceneItem{{Name: "Good morning", Icon: "sunrise"}, {Name: "Evening", Icon: "sunset"},
			{Name: "Movie night", Icon: "clapperboard"}, {Name: "Away", Icon: "house"}, {Name: "Bedtime", Icon: "moon"}, {Name: "Focus", Icon: "lamp-desk"}}
		before := *runs
		changed := home.Scene(c, scenes, active, home.SceneOptions{Runs: runs}).Changed()
		if changed || *runs != before {
			log(fmt.Sprintf("runs %d, scene %s", *runs, *active))
		}
	})
}

func demoHomeRooms(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	picked := textState(root, "room", "Living room")
	root.Children(func() {
		rooms := []home.Room{{Name: "Living room", Icon: "sofa", On: 3}, {Name: "Kitchen", Icon: "utensils-crossed", On: 1},
			{Name: "Bedroom", Icon: "bed-double"}, {Name: "Bathroom", Icon: "bath"}, {Name: "Studio", Icon: "palette", On: 2}}
		if home.Rooms(c, rooms, picked, home.RoomsOptions{}).Changed() {
			log("room " + *picked)
		}
	})
}

func demoHomeDevice(c *ui.Context, log func(string)) {
	root := ui.Grid(c).Columns(2).Width(400).Gap(10)
	coffee := boolState(root, "coffee", false)
	tv := boolState(root, "tv", true)
	root.Children(func() {
		if home.Device(c, "Coffee maker", "coffee", coffee, home.DeviceOptions{Room: "Kitchen"}).Changed() {
			log(fmt.Sprintf("coffee maker on %v", *coffee))
		}
		if home.Device(c, "Television", "tv", tv, home.DeviceOptions{Room: "Living room"}).Changed() {
			log(fmt.Sprintf("television on %v", *tv))
		}
	})
}

func demoHomeFan(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	speed := intState(root, "speed", 2)
	root.Children(func() {
		if home.Fan(c, "Ceiling fan", speed, home.FanOptions{Room: "Bedroom"}).Changed() {
			log(fmt.Sprintf("speed %d", *speed))
		}
	})
}

func demoHomeBlinds(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	open := numberState(root, "open", 0.35)
	root.Children(func() {
		if home.Blinds(c, "Bay window", open, home.BlindsOptions{Room: "Living room"}).Changed() {
			log(fmt.Sprintf("open %d%%", int(math.Round(*open*100))))
		}
	})
}

func demoHomeEnergy(c *ui.Context, log func(string)) {
	use := []float64{0.3, 0.2, 0.2, 0.2, 0.3, 0.4, 0.9, 1.1, 0.8, 0.5, 0.4, 0.5, 0.7, 0.6, 0.5, 0.5, 0.7, 1.0, 1.6, 1.4, 1.2, 0.9, 0.6, 0.4}
	labels := make([]string, len(use))
	for i := range labels {
		labels[i] = fmt.Sprintf("%02d:00", i)
	}
	ui.Column(c).Width(440).Children(func() {
		home.Energy(c, use, home.EnergyOptions{Labels: labels})
	})
}

func demoHomePlug(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	on := boolState(root, "on", true)
	root.Children(func() {
		if home.Plug(c, "Desk lamp plug", on, 42.6, home.PlugOptions{Room: "Studio", Today: 0.38}).Changed() {
			log(fmt.Sprintf("plug %v", *on))
		}
	})
}

func demoHomeSpeaker(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	playing := boolState(root, "playing", true)
	volume := numberState(root, "volume", 0.45)
	root.Children(func() {
		tr := home.Track{Title: "Holocene", Artist: "Bon Iver", Art: photos()[2]}
		if home.Speaker(c, "Kitchen speaker", tr, playing, volume, home.SpeakerOptions{Room: "Kitchen"}).Changed() {
			log(fmt.Sprintf("playing %v, volume %d%%", *playing, int(math.Round(*volume*100))))
		}
	})
}

func demoHomeAirQuality(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		home.AirQuality(c, 42, home.AirQualityOptions{Pollutants: []home.Pollutant{
			{Name: "PM2.5", Value: "8", Unit: "µg/m³"}, {Name: "PM10", Value: "14", Unit: "µg/m³"},
			{Name: "CO₂", Value: "640", Unit: "ppm"}, {Name: "VOC", Value: "0.3", Unit: "mg/m³"}}})
	})
}

func demoHomeHumidity(c *ui.Context, log func(string)) {
	ui.Column(c).Width(340).Children(func() {
		home.Humidity(c, 48, home.HumidityOptions{Room: "Bedroom"})
	})
}

func demoHomeSensors(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		home.Sensors(c, []home.Sensor{
			{Name: "Front door", Room: "Entrance", Since: "Closed 2h ago"},
			{Name: "Balcony door", Room: "Living room", Since: "Open 12 min", Open: true},
			{Name: "Kitchen window", Room: "Kitchen", Kind: home.WindowSensor, Since: "Open 40 min", Open: true},
			{Name: "Bedroom window", Room: "Bedroom", Kind: home.WindowSensor, Since: "Closed since 07:10"},
		}, home.SensorsOptions{})
	})
}

func demoHomeAutomation(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	on := boolState(root, "on", true)
	root.Children(func() {
		r := home.Rule{Name: "Sunset lights", When: "The sun sets", Then: "Turn on the living room lights at 40%",
			WhenIcon: "sunset", ThenIcon: "lightbulb", Last: "Ran yesterday at 18:14"}
		if home.Automation(c, r, on, home.AutomationOptions{}).Changed() {
			log(fmt.Sprintf("rule %v", *on))
		}
	})
}

func demoHomeVacuum(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	state := ui.Local(root, "state", func() home.VacuumState { return home.Docked })
	root.Children(func() {
		if home.Vacuum(c, "Robo vacuum", state, 0.86, home.VacuumOptions{Room: "Living room"}).Changed() {
			log(fmt.Sprintf("vacuum %d", *state))
		}
	})
}

func demoHomeGarage(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	open := boolState(root, "open", false)
	root.Children(func() {
		if home.Garage(c, "Garage", open, home.GarageOptions{}).Changed() {
			log(fmt.Sprintf("door %v", *open))
		}
	})
}

func demoHomeSprinklers(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	running := textState(root, "running", "")
	root.Children(func() {
		zones := []home.Zone{{Name: "Front lawn", Minutes: 15}, {Name: "Back garden", Minutes: 20},
			{Name: "Vegetable beds", Minutes: 10}, {Name: "Hedges", Minutes: 8}}
		if home.Sprinklers(c, zones, running, home.SprinklersOptions{}).Changed() {
			log("running " + *running)
		}
	})
}

func demoHomeSummary(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		home.HomeSummary(c, home.HomeStatus{Name: "14 Alder Street", Inside: 21.5, Energy: 8.2, Lights: 4, On: 7, Locked: true},
			home.HomeSummaryOptions{})
	})
}
