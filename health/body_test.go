package health

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestActivityRingsLegend(t *testing.T) {
	rings := [3]Ring{{Label: "Move", Value: 1300, Goal: 600, Unit: "kcal"}, {Label: "Exercise", Value: 38, Goal: 30, Unit: "min"},
		{Label: "Stand", Value: 9, Goal: 12, Unit: "h"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ActivityRings(c, rings, ActivityRingsOptions{}) }), 420, 300)
	settles(t, tt)
	if !exact(tt, "1,300 / 600 kcal") || !exact(tt, "9 / 12 h") {
		t.Errorf("legend: %q", tt.Texts())
	}
	panics(t, "goal of 1", func(c *ui.Context) {
		ActivityRings(c, [3]Ring{{Label: "a"}, {Label: "b", Goal: 1}, {Label: "c", Goal: 1}}, ActivityRingsOptions{})
	})
	panics(t, "repeated", func(c *ui.Context) {
		ActivityRings(c, [3]Ring{{Label: "a", Goal: 1}, {Label: "a", Goal: 1}, {Label: "c", Goal: 1}}, ActivityRingsOptions{})
	})
}

func TestStepsGoal(t *testing.T) {
	steps := 7840
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Steps(c, steps, 10000, StepsOptions{Week: []int{9120, 12030}, Days: []string{"Fri", "Sat"}})
	}), 400, 400)
	settles(t, tt)
	if !exact(tt, "7,840") || !exact(tt, "2,160 to go") || !exact(tt, "Sat") {
		t.Errorf("under goal: %q", tt.Texts())
	}
	steps = 10000
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "Goal reached") {
		t.Errorf("at goal: %q", tt.Texts())
	}
	panics(t, "label for each", func(c *ui.Context) { Steps(c, 1, 10, StepsOptions{Week: []int{1, 2}, Days: []string{"Mon"}}) })
	panics(t, "day of -1", func(c *ui.Context) { Steps(c, 1, 10, StepsOptions{Week: []int{-1}, Days: []string{"Mon"}}) })
}

func TestHeartRateZones(t *testing.T) {
	bpm := 72
	tt := ui.NewTester(padded(func(c *ui.Context) { HeartRate(c, bpm, []int{64, 124, 72}, HeartRateOptions{}) }), 400, 400)
	settles(t, tt)
	if !tt.FramePending() || !exact(tt, "Resting") || !exact(tt, "64–124 bpm") {
		t.Errorf("resting: pending %v %q", tt.FramePending(), tt.Texts())
	}
	bpm = 150
	tt.Frame()
	if !exact(tt, "Cardio") {
		t.Errorf("150 bpm: %q", tt.Texts())
	}
	panics(t, "sample 1 is 300", func(c *ui.Context) { HeartRate(c, 70, []int{70, 300}, HeartRateOptions{}) })
}

func TestSleepTotals(t *testing.T) {
	spans := []SleepSpan{{Stage: Awake, Minutes: 26}, {Stage: Light, Minutes: 192}, {Stage: Deep, Minutes: 102}, {Stage: REM, Minutes: 106}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Sleep(c, spans, SleepOptions{Start: time.Date(2026, 10, 24, 23, 10, 0, 0, time.UTC)})
	}), 440, 500)
	settles(t, tt)
	for _, s := range []string{"6h 40m asleep", "26m", "3h 12m", "23:10", "06:16"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	panics(t, "within a day", func(c *ui.Context) { Sleep(c, []SleepSpan{{Light, 1000}, {Deep, 500}}, SleepOptions{}) })
	panics(t, "needs a stage", func(c *ui.Context) { Sleep(c, []SleepSpan{{Stage(9), 10}}, SleepOptions{}) })
}

func TestWaterGlasses(t *testing.T) {
	glasses := 5
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Water(c, &glasses, 8, WaterOptions{}).Changed() {
			changed++
		}
	}), 400, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Glass 7")
	click(t, tt, "Glass 7")
	click(t, tt, "Add a glass")
	cues(t, before, sound.Pop, sound.Pop, sound.Pop)
	if glasses != 7 || changed != 3 || !exact(tt, "1,750 of 2,000 ml") {
		t.Errorf("glasses %d changed %d: %q", glasses, changed, tt.Texts())
	}
	glasses = 12
	tt.Frame()
	_ = tt.Click("Add a glass")
	tt.Frame()
	if glasses != 12 {
		t.Errorf("added past 12: %d", glasses)
	}
	panics(t, "goal and glasses of 0 to 12", func(c *ui.Context) { g := 0; Water(c, &g, 13, WaterOptions{}) })
}

func TestCaloriesLeftAndOver(t *testing.T) {
	eaten := 1640
	tt := ui.NewTester(padded(func(c *ui.Context) { Calories(c, eaten, 430, 2100, CaloriesOptions{}) }), 420, 300)
	settles(t, tt)
	if !exact(tt, "890 left") {
		t.Errorf("left: %q", tt.Texts())
	}
	eaten = 3430
	tt.Frame()
	if !exact(tt, "900 over") {
		t.Errorf("over: %q", tt.Texts())
	}
	panics(t, "goal of 1", func(c *ui.Context) { Calories(c, 1, 1, 0, CaloriesOptions{}) })
}

func TestMacrosShares(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Macros(c, [3]Macro{{Grams: 92, Goal: 120}, {Grams: 186, Goal: 230}, {Grams: 71, Goal: 65}}, MacrosOptions{})
	}), 400, 400)
	settles(t, tt)
	for _, s := range []string{"21%", "42%", "36%", "71g of 65g"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	panics(t, "and a goal", func(c *ui.Context) { Macros(c, [3]Macro{{1, 0}, {1, 1}, {1, 1}}, MacrosOptions{}) })
}

func TestWorkoutStats(t *testing.T) {
	s := Session{Kind: "Morning run", Icon: "footprints", When: time.Date(2026, 10, 25, 7, 12, 0, 0, time.UTC), Minutes: 42, Calories: 468,
		Kilometers: 7.6, HeartRate: 152}
	tt := ui.NewTester(padded(func(c *ui.Context) { Workout(c, s, WorkoutOptions{}) }), 420, 500)
	settles(t, tt)
	for _, w := range []string{"42", "min", "5'32\"", "7.60", "468", "152", "Sun, Oct 25 · 07:12"} {
		if !exact(tt, w) {
			t.Errorf("missing %q: %q", w, tt.Texts())
		}
	}
	s = Session{Kind: "Yoga", Icon: "flower", When: s.When, Minutes: 75}
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "1:15") || exact(tt, "Pace") || exact(tt, "Avg heart rate") {
		t.Errorf("yoga: %q", tt.Texts())
	}
	panics(t, "1 to 1440 minutes", func(c *ui.Context) {
		Workout(c, Session{Kind: "a", Icon: "b", When: s.When, Minutes: 30, Kilometers: 0.005}, WorkoutOptions{})
	})
}

func TestCycleDays(t *testing.T) {
	day := 12
	tt := ui.NewTester(padded(func(c *ui.Context) { Cycle(c, day, 28, CycleOptions{}) }), 400, 500)
	settles(t, tt)
	if !exact(tt, "Day 12") || !exact(tt, "Fertile") || !exact(tt, "Period in 17 days") {
		t.Errorf("day 12: %q", tt.Texts())
	}
	day = 28
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "Period tomorrow") || !exact(tt, "Luteal") {
		t.Errorf("day 28: %q", tt.Texts())
	}
	day = 2
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "Period") {
		t.Errorf("day 2: %q", tt.Texts())
	}
	panics(t, "21 to 45 days", func(c *ui.Context) { Cycle(c, 1, 50, CycleOptions{}) })
}

func TestBMIBands(t *testing.T) {
	for _, c := range []struct {
		kg, m       float64
		shown, band string
	}{{66.6, 1.68, "23.5", "Healthy"}, {70.55, 1.68, "24.9", "Healthy"}, {84.7, 1.68, "30.0", "Obese"}, {84.67, 1.68, "29.9", "Overweight"}, {50, 1.7, "17.3", "Underweight"}} {
		tt := ui.NewTester(padded(func(ctx *ui.Context) { BMI(ctx, c.kg, c.m, BMIOptions{}) }), 400, 300)
		settles(t, tt)
		if !exact(tt, c.shown) || !exact(tt, c.band) {
			t.Errorf("%v kg %v m: %q", c.kg, c.m, tt.Texts())
		}
	}
	panics(t, "0.4 to 2.6 m", func(c *ui.Context) { BMI(c, 70, 3, BMIOptions{}) })
}
