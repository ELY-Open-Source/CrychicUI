package health

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestIntervalsRounds(t *testing.T) {
	runs, resets := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, run, reset := Intervals(c, 2*time.Second, time.Second, 2, IntervalsOptions{})
		if run.Clicked() {
			runs++
		}
		if reset.Clicked() {
			resets++
		}
	}), 400, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Start")
	tt.Advance(2100 * time.Millisecond)
	if !exact(tt, "Rest") || !exact(tt, "Round 1 of 2") {
		t.Errorf("after work: %q", tt.Texts())
	}
	tt.Advance(time.Second)
	if !exact(tt, "Work") || !exact(tt, "Round 2 of 2") {
		t.Errorf("after rest: %q", tt.Texts())
	}
	tt.Advance(2 * time.Second)
	cues(t, before, sound.Tap, sound.Notify, sound.Notify, sound.Notify)
	if !exact(tt, "Done") || runs != 1 {
		t.Errorf("after the last round: runs %d %q", runs, tt.Texts())
	}
	before = baseline()
	_ = tt.Click("Start")
	tt.Advance(100 * time.Millisecond)
	cues(t, before)
	click(t, tt, "Reset")
	if resets != 1 || !exact(tt, "Round 1 of 2") || !exact(tt, "0:02") {
		t.Errorf("reset: %q", tt.Texts())
	}
	panics(t, "1 to 50 rounds", func(c *ui.Context) { Intervals(c, time.Second, time.Second, 51, IntervalsOptions{}) })
}

func TestBreathePhases(t *testing.T) {
	running := false
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Breathe(c, &running, BreatheOptions{}).Changed() {
			changed++
		}
	}), 400, 500)
	settles(t, tt)
	if tt.FramePending() || !exact(tt, "Ready when you are") {
		t.Errorf("at rest: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Start")
	cues(t, before, sound.Tap)
	tt.Frame()
	if !running || changed != 1 || !exact(tt, "Breathe in") || !exact(tt, "4") {
		t.Errorf("started: %q", tt.Texts())
	}
	tt.Advance(4100 * time.Millisecond)
	if !exact(tt, "Hold") || !tt.FramePending() {
		t.Errorf("after inhale: %q", tt.Texts())
	}
	tt.Advance(4 * time.Second)
	if !exact(tt, "Breathe out") {
		t.Errorf("after hold: %q", tt.Texts())
	}
	click(t, tt, "Stop")
	settles(t, tt)
	if running || !exact(tt, "Ready when you are") || tt.FramePending() {
		t.Errorf("stopped: %q", tt.Texts())
	}
	panics(t, "1 to 60 s", func(c *ui.Context) { r := false; Breathe(c, &r, BreatheOptions{Hold: 500 * time.Millisecond}) })
}

func TestMoodPicks(t *testing.T) {
	mood, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Mood(c, []MoodDay{{Label: "Mon", Mood: 3}, {Label: "Tue", Mood: -1}}, &mood, MoodOptions{}).Changed() {
			changed++
		}
	}), 420, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Good")
	click(t, tt, "Good")
	cues(t, before, sound.Select)
	if mood != 3 || changed != 1 {
		t.Errorf("mood %d changed %d", mood, changed)
	}
	panics(t, "mood of -1 to 4", func(c *ui.Context) { m := 5; Mood(c, nil, &m, MoodOptions{}) })
}

func TestHabitsToggle(t *testing.T) {
	habits := []string{"Stretch", "Read"}
	days := []string{"Fri", "Sat", "Sun"}
	done := [][]bool{{true, true, false}, {false, true, true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Habits(c, habits, days, &done, HabitsOptions{}) }), 460, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Stretch on Sun")
	cues(t, before, sound.Toggle)
	if !done[0][2] || !exact(tt, "3") {
		t.Errorf("toggle: %v %q", done, tt.Texts())
	}
	swapUnderPress(t, tt, "Read on Fri", func() { habits[0], habits[1] = habits[1], habits[0]; done[0], done[1] = done[1], done[0] })
	if done[0][0] || !done[1][0] {
		t.Errorf("a reorder under the press toggled %v", done)
	}
	panics(t, "has 2 marks for 3 days", func(c *ui.Context) {
		d := [][]bool{{true, true}}
		Habits(c, []string{"a"}, []string{"x", "y", "z"}, &d, HabitsOptions{})
	})
}

func TestWeightChange(t *testing.T) {
	points := []WeightPoint{{Label: "Aug 2", Value: 68.4}, {Label: "Sep 1", Value: 67.5}, {Label: "Oct 25", Value: 66.6}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Weight(c, points, WeightOptions{Goal: 65}) }), 420, 400)
	settles(t, tt)
	for _, s := range []string{"66.6", "-1.8 kg", "Goal 65.0 kg", "Aug 2", "Oct 25"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	points[2].Value = 68.8
	tt.Frame()
	if !exact(tt, "+0.4 kg") {
		t.Errorf("a gain lacks its sign: %q", tt.Texts())
	}
	panics(t, "goal of 0 to 1000", func(c *ui.Context) { Weight(c, points, WeightOptions{Goal: 2000}) })
	panics(t, "value of 0 to 1000", func(c *ui.Context) { Weight(c, []WeightPoint{{"a", 1}, {"b", -1}}, WeightOptions{}) })
}

func TestMedicationMarks(t *testing.T) {
	doses := []Dose{{Name: "Vitamin D", At: "08:00", Taken: true}, {Name: "Iron", At: "08:00"}, {Name: "Iron", At: "20:00"}}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Medication(c, &doses, MedicationOptions{}).Changed() {
			changed++
		}
	}), 420, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Iron at 08:00")
	cues(t, before, sound.Toggle)
	if !doses[1].Taken || changed != 1 || !exact(tt, "2 of 3 taken") {
		t.Errorf("mark: %v %q", doses, tt.Texts())
	}
	swapUnderPress(t, tt, "Iron at 20:00", func() { doses[1], doses[2] = doses[2], doses[1] })
	if doses[1].Taken || !doses[2].Taken {
		t.Errorf("a reorder under the press marked %v", doses)
	}
	click(t, tt, "Iron at 20:00")
	settles(t, tt)
	if !exact(tt, "All taken") {
		t.Errorf("all: %q", tt.Texts())
	}
	panics(t, "repeated", func(c *ui.Context) {
		d := []Dose{{Name: "a", At: "1"}, {Name: "a", At: "1"}}
		Medication(c, &d, MedicationOptions{})
	})
}

func TestMealAndVitals(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Meal(c, Dish{Name: "Salmon bowl", Calories: 1640, Protein: 38}, MealOptions{})
	}), 400, 500)
	settles(t, tt)
	if !exact(tt, "1,640 kcal") || !exact(tt, "38g") {
		t.Errorf("meal: %q", tt.Texts())
	}
	panics(t, "nutrients of 0", func(c *ui.Context) { Meal(c, Dish{Name: "a", Fat: -1}, MealOptions{}) })
	tt = ui.NewTester(padded(func(c *ui.Context) {
		Vitals(c, []Vital{{Name: "Blood pressure", Value: "128/84", Icon: "heart", Range: High}, {Name: "Oxygen", Value: "98", Icon: "wind"}}, VitalsOptions{})
	}), 420, 400)
	settles(t, tt)
	if !exact(tt, "High") || !exact(tt, "Normal") || !exact(tt, "128/84") {
		t.Errorf("vitals: %q", tt.Texts())
	}
	panics(t, "repeated", func(c *ui.Context) {
		Vitals(c, []Vital{{"a", "1", "", "heart", InRange}, {"a", "2", "", "heart", Low}}, VitalsOptions{})
	})
}

func TestBodyMapPicks(t *testing.T) {
	region, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if BodyMap(c, &region, BodyMapOptions{}).Changed() {
			changed++
		}
	}), 400, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Left arm")
	click(t, tt, "Left arm")
	cues(t, before, sound.Select)
	if region != 4 || changed != 1 {
		t.Errorf("region %d changed %d", region, changed)
	}
	r, l := rectOf(t, tt, "Right arm"), rectOf(t, tt, "Left arm")
	if r.X >= l.X {
		t.Error("the right arm should sit on the viewer's left")
	}
	panics(t, "region of -1 to 6", func(c *ui.Context) { r := 7; BodyMap(c, &r, BodyMapOptions{}) })
}

func TestGoalStreakRuns(t *testing.T) {
	days := []bool{true, true, true, false, true}
	tt := ui.NewTester(padded(func(c *ui.Context) { GoalStreak(c, days, GoalStreakOptions{}) }), 400, 300)
	settles(t, tt)
	if !exact(tt, "1 day in a row") || !exact(tt, "Best 3") {
		t.Errorf("one: %q", tt.Texts())
	}
	days = []bool{true, false, true, true, true, true}
	tt.Frame()
	if !exact(tt, "4 days in a row") || !exact(tt, "Best 4") {
		t.Errorf("four: %q", tt.Texts())
	}
	panics(t, "1 to 31 days", func(c *ui.Context) { GoalStreak(c, nil, GoalStreakOptions{}) })
}
