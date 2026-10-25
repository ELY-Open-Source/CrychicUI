package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/health"
	"github.com/egoist/mygo/ui"
)

var healthEntries = []Entry{
	{ID: 1351, Name: "ActivityRings", Category: "health", Summary: "Three rings closing.", Demo: demoHealthActivityRings},
	{ID: 1352, Name: "Steps", Category: "health", Summary: "Step count toward a goal.", Demo: demoHealthSteps},
	{ID: 1353, Name: "HeartRate", Category: "health", Summary: "Beating line and bpm.", Demo: demoHealthHeartRate},
	{ID: 1354, Name: "Sleep", Category: "health", Summary: "Sleep stages bars.", Demo: demoHealthSleep},
	{ID: 1355, Name: "Water", Category: "health", Summary: "Glasses filled; cue Pop.", Demo: demoHealthWater},
	{ID: 1356, Name: "Calories", Category: "health", Summary: "Eaten, burned, left ring.", Demo: demoHealthCalories},
	{ID: 1357, Name: "Macros", Category: "health", Summary: "Protein, carbs, fat bars.", Demo: demoHealthMacros},
	{ID: 1358, Name: "Workout", Category: "health", Summary: "Workout card with stats.", Demo: demoHealthWorkout},
	{ID: 1359, Name: "Intervals", Category: "health", Summary: "Work and rest timer; cue Notify.", Demo: demoHealthIntervals},
	{ID: 1360, Name: "Breathe", Category: "health", Summary: "Circle growing and shrinking.", Demo: demoHealthBreathe},
	{ID: 1361, Name: "Mood", Category: "health", Summary: "Mood faces by day; cue Select.", Demo: demoHealthMood},
	{ID: 1362, Name: "Habits", Category: "health", Summary: "Habit grid by day; cue Toggle.", Demo: demoHealthHabits},
	{ID: 1363, Name: "Weight", Category: "health", Summary: "Weight trend line.", Demo: demoHealthWeight},
	{ID: 1364, Name: "Medication", Category: "health", Summary: "Doses with taken marks; cue Toggle.", Demo: demoHealthMedication},
	{ID: 1365, Name: "Cycle", Category: "health", Summary: "Cycle day ring.", Demo: demoHealthCycle},
	{ID: 1366, Name: "BMI", Category: "health", Summary: "BMI on a band scale.", Demo: demoHealthBMI},
	{ID: 1367, Name: "Meal", Category: "health", Summary: "Meal card with nutrients.", Demo: demoHealthMeal},
	{ID: 1368, Name: "BodyMap", Category: "health", Summary: "Body regions to pick; cue Select.", Demo: demoHealthBodyMap},
	{ID: 1369, Name: "Vitals", Category: "health", Summary: "Vital tiles.", Demo: demoHealthVitals},
	{ID: 1370, Name: "GoalStreak", Category: "health", Summary: "Days met in a row.", Demo: demoHealthGoalStreak},
}

var demoWeekDays = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

func demoHealthActivityRings(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		health.ActivityRings(c, [3]health.Ring{{Label: "Move", Value: 432, Goal: 600, Unit: "kcal"},
			{Label: "Exercise", Value: 38, Goal: 30, Unit: "min"}, {Label: "Stand", Value: 9, Goal: 12, Unit: "h"}}, health.ActivityRingsOptions{})
	})
}

func demoHealthSteps(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		health.Steps(c, 7840, 10000, health.StepsOptions{Week: []int{9120, 11340, 6210, 10480, 8790, 12030, 7840}, Days: demoWeekDays})
	})
}

func demoHealthHeartRate(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		health.HeartRate(c, 72, []int{64, 66, 70, 75, 82, 96, 118, 124, 109, 92, 80, 76, 72}, health.HeartRateOptions{})
	})
}

func demoHealthSleep(c *ui.Context, log func(string)) {
	spans := []health.SleepSpan{{Stage: health.Awake, Minutes: 12}, {Stage: health.Light, Minutes: 48}, {Stage: health.Deep, Minutes: 62},
		{Stage: health.Light, Minutes: 34}, {Stage: health.REM, Minutes: 28}, {Stage: health.Light, Minutes: 52}, {Stage: health.Deep, Minutes: 40},
		{Stage: health.REM, Minutes: 36}, {Stage: health.Awake, Minutes: 6}, {Stage: health.Light, Minutes: 58}, {Stage: health.REM, Minutes: 42},
		{Stage: health.Awake, Minutes: 8}}
	ui.Column(c).Width(440).Children(func() {
		health.Sleep(c, spans, health.SleepOptions{Start: time.Date(2026, 10, 24, 23, 10, 0, 0, time.UTC)})
	})
}

func demoHealthWater(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	glasses := ui.Local(root, "glasses", func() int { return 5 })
	root.Children(func() {
		if health.Water(c, glasses, 8, health.WaterOptions{}).Changed() {
			log(fmt.Sprintf("glasses %d", *glasses))
		}
	})
}

func demoHealthCalories(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() { health.Calories(c, 1640, 430, 2100, health.CaloriesOptions{}) })
}

func demoHealthMacros(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		health.Macros(c, [3]health.Macro{{Grams: 92, Goal: 120}, {Grams: 186, Goal: 230}, {Grams: 71, Goal: 65}}, health.MacrosOptions{})
	})
}

func demoHealthWorkout(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		health.Workout(c, health.Session{Kind: "Morning run", Icon: "footprints", When: time.Date(2026, 10, 25, 7, 12, 0, 0, time.UTC),
			Minutes: 42, Calories: 468, Kilometers: 7.6, HeartRate: 152}, health.WorkoutOptions{})
	})
}

func demoHealthIntervals(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		_, run, reset := health.Intervals(c, 40*time.Second, 20*time.Second, 8, health.IntervalsOptions{})
		if run.Clicked() {
			log("timer toggled")
		}
		if reset.Clicked() {
			log("timer reset")
		}
	})
}

func demoHealthBreathe(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	running := ui.Local(root, "running", func() bool { return false })
	root.Children(func() {
		if health.Breathe(c, running, health.BreatheOptions{}).Changed() {
			log(fmt.Sprintf("breathing %v", *running))
		}
	})
}

func demoHealthMood(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	mood := ui.Local(root, "mood", func() int { return -1 })
	week := []health.MoodDay{{Label: "Mon", Mood: 3}, {Label: "Tue", Mood: 4}, {Label: "Wed", Mood: 2}, {Label: "Thu", Mood: 1},
		{Label: "Fri", Mood: 3}, {Label: "Sat", Mood: -1}}
	root.Children(func() {
		if health.Mood(c, week, mood, health.MoodOptions{}).Changed() {
			log(fmt.Sprintf("mood %d", *mood))
		}
	})
}

func demoHealthHabits(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	habits := []string{"Stretch", "Read 20 pages", "No phone in bed", "Walk outside"}
	done := ui.Local(root, "done", func() [][]bool {
		return [][]bool{{true, true, false, true, true, true, false}, {true, false, true, true, false, true, true},
			{false, false, true, true, true, true, true}, {true, true, true, false, true, false, false}}
	})
	root.Children(func() {
		if health.Habits(c, habits, demoWeekDays, done, health.HabitsOptions{}).Changed() {
			n := 0
			for _, row := range *done {
				for _, m := range row {
					if m {
						n++
					}
				}
			}
			log(fmt.Sprintf("marks %d", n))
		}
	})
}

func demoHealthWeight(c *ui.Context, log func(string)) {
	values := []float64{68.4, 68.1, 68.3, 67.8, 67.5, 67.6, 67.2, 66.9, 67.0, 66.6}
	points := make([]health.WeightPoint, len(values))
	for i, v := range values {
		points[i] = health.WeightPoint{Label: fmt.Sprintf("Aug %d", 2+7*i), Value: v}
	}
	points[len(points)-1].Label = "Oct 25"
	ui.Column(c).Width(420).Children(func() { health.Weight(c, points, health.WeightOptions{Goal: 65}) })
}

func demoHealthMedication(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	doses := ui.Local(root, "doses", func() []health.Dose {
		return []health.Dose{{Name: "Vitamin D", Amount: "1,000 IU", At: "08:00", Taken: true}, {Name: "Iron", Amount: "65 mg", At: "08:00"},
			{Name: "Ibuprofen", Amount: "200 mg", At: "13:00"}, {Name: "Magnesium", Amount: "250 mg", At: "21:30"}}
	})
	root.Children(func() {
		if health.Medication(c, doses, health.MedicationOptions{}).Changed() {
			n := 0
			for _, d := range *doses {
				if d.Taken {
					n++
				}
			}
			log(fmt.Sprintf("taken %d", n))
		}
	})
}

func demoHealthCycle(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() { health.Cycle(c, 12, 28, health.CycleOptions{}) })
}

func demoHealthBMI(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() { health.BMI(c, 66.6, 1.68, health.BMIOptions{}) })
}

func demoHealthMeal(c *ui.Context, log func(string)) {
	ui.Column(c).Width(360).Children(func() {
		health.Meal(c, health.Dish{Name: "Salmon grain bowl", Kind: "Lunch", At: "12:40", Calories: 640, Protein: 38, Carbs: 62, Fat: 24,
			Items: []string{"Roasted salmon, 120 g", "Brown rice and quinoa", "Edamame, cucumber, avocado", "Miso dressing"}}, health.MealOptions{})
	})
}

func demoHealthBodyMap(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(340)
	region := ui.Local(root, "region", func() int { return -1 })
	root.Children(func() {
		if health.BodyMap(c, region, health.BodyMapOptions{}).Changed() {
			log(fmt.Sprintf("region %d", *region))
		}
	})
}

func demoHealthVitals(c *ui.Context, log func(string)) {
	ui.Column(c).Width(420).Children(func() {
		health.Vitals(c, []health.Vital{{Name: "Blood pressure", Value: "128/84", Unit: "mmHg", Icon: "heart-pulse", Range: health.High},
			{Name: "Resting heart", Value: "62", Unit: "bpm", Icon: "heart"}, {Name: "Blood oxygen", Value: "98", Unit: "%", Icon: "wind"},
			{Name: "Temperature", Value: "36.4", Unit: "°C", Icon: "thermometer"}}, health.VitalsOptions{})
	})
}

func demoHealthGoalStreak(c *ui.Context, log func(string)) {
	days := []bool{true, true, false, true, true, true, true, false, true, true, true, true, true, true}
	ui.Column(c).Width(400).Children(func() { health.GoalStreak(c, days, health.GoalStreakOptions{}) })
}
