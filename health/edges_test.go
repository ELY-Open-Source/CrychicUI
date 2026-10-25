package health

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestHabitsDayReorderUnderPress(t *testing.T) {
	habits := []string{"Stretch"}
	days := []string{"Mon", "Tue"}
	done := [][]bool{{false, false}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Habits(c, habits, days, &done, HabitsOptions{}) }), 460, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Stretch on Mon", func() { days[0], days[1] = days[1], days[0] })
	if done[0][0] || done[0][1] {
		t.Errorf("a day reorder under the press toggled %v", done)
	}
}

func TestBreatheCopyBesideButton(t *testing.T) {
	running := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Breathe(c, &running, BreatheOptions{Ready: "toggle"}) }), 400, 500)
	settles(t, tt)
	click(t, tt, "Start")
	if !running {
		t.Error("Start did not run")
	}
}

func TestWorkoutHugeValueKeepsUnit(t *testing.T) {
	s := Session{Kind: "Run", Icon: "footprints", When: time.Date(2026, 10, 25, 7, 0, 0, 0, time.UTC), Minutes: 30, Calories: math.MaxInt}
	tt := ui.NewTester(padded(func(c *ui.Context) { Workout(c, s, WorkoutOptions{}) }), 320, 400)
	settles(t, tt)
	value, unit := rectOf(t, tt, "9,223,372,036,854,775,807"), rectOf(t, tt, "kcal")
	inside(t, "kcal", unit, value.X+value.W, 294)
}

func TestVitalsLongValueKeepsUnit(t *testing.T) {
	vitals := []Vital{
		{Name: "Blood pressure", Value: "120/80–130/85", Unit: "mmHg", Icon: "heart"},
		{Name: "Resting pulse", Value: "58", Unit: "bpm", Icon: "activity"},
	}
	wide := ui.NewTester(padded(func(c *ui.Context) { Vitals(c, vitals, VitalsOptions{}) }), 900, 400)
	settles(t, wide)
	full := rectOf(t, wide, "mmHg").W
	tt := ui.NewTester(padded(func(c *ui.Context) { Vitals(c, vitals, VitalsOptions{}) }), 320, 400)
	settles(t, tt)
	value, unit := rectOf(t, tt, "120/80–130/85"), rectOf(t, tt, "mmHg")
	inside(t, "mmHg", unit, value.X+value.W, 310)
	if unit.W < full-0.5 {
		t.Errorf("mmHg cut to %.1f of %.1f", unit.W, full)
	}
}
