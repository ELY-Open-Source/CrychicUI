package display

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func TestLocale(t *testing.T) {
	label := "EN"
	tt := ui.NewTester(func(c *ui.Context) { Locale(c, label) }, 300, 200)
	label = "中文"
	tt.Frame()
	settles(t, tt, "Locale")
	find(t, tt, "中文")
}

func TestPresence(t *testing.T) {
	status, label := PresenceOffline, "Offline"
	tt := ui.NewTester(func(c *ui.Context) { Presence(c, status, label) }, 300, 200)
	status, label = PresenceOnline, "Online"
	tt.Frame()
	settles(t, tt, "Presence")
	if !tt.HasText("Online") || tt.HasText("Offline") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestOrSeparator(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { OrSeparator(c, "OR") }, 300, 200)
	settles(t, tt, "OrSeparator")
	find(t, tt, "OR")
}

func TestWatermark(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Watermark(c, "DRAFT", func() { ui.Box(c).Size(180, 90) })
	}, 300, 200)
	settles(t, tt, "Watermark")
	if !tt.HasText("DRAFT") {
		t.Error("watermark text missing")
	}
}

func TestAccentCard(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { AccentCard(c, func() { ui.Text(c, "Body") }) }, 300, 200)
	settles(t, tt, "AccentCard")
	find(t, tt, "Body")
}

func TestComparison(t *testing.T) {
	after := 12.0
	tt := ui.NewTester(func(c *ui.Context) { Comparison(c, 10, after) }, 300, 200)
	after = 14
	tt.Frame()
	settles(t, tt, "Comparison")
	find(t, tt, "14")
}

func TestScore(t *testing.T) {
	value := 42.0
	tt := ui.NewTester(func(c *ui.Context) { Score(c, value, ScoreOptions{}) }, 300, 200)
	value = 84
	tt.Frame()
	settles(t, tt, "Score ring")
	find(t, tt, "84")
	panics(t, "Score value", func() { ui.Render(func(c *ui.Context) { Score(c, 101, ScoreOptions{}) }, 10, 10, 1) })
}

func TestLevel(t *testing.T) {
	progress := 0.25
	tt := ui.NewTester(func(c *ui.Context) { Level(c, 2, progress, LevelOptions{Label: "Stage"}) }, 300, 200)
	progress = 0.75
	tt.Frame()
	settles(t, tt, "Level fill")
	find(t, tt, "Stage 2")
	panics(t, "Level progress", func() { ui.Render(func(c *ui.Context) { Level(c, 2, 2, LevelOptions{}) }, 10, 10, 1) })
}

func TestStreak(t *testing.T) {
	days := 4
	tt := ui.NewTester(func(c *ui.Context) { Streak(c, days) }, 300, 200)
	if !tt.FramePending() {
		t.Error("flame did not flicker")
	}
	days = 5
	tt.Frame()
	tt.Advance(500 * time.Millisecond)
	find(t, tt, "5")
	if !tt.FramePending() {
		t.Error("flame stopped")
	}
	panics(t, "Streak days", func() { ui.Render(func(c *ui.Context) { Streak(c, -1) }, 10, 10, 1) })
}

func TestAchievement(t *testing.T) {
	unlocked := false
	tt := ui.NewTester(func(c *ui.Context) { Achievement(c, "First ship", "trophy", unlocked) }, 300, 200)
	unlocked = true
	tt.Frame()
	settles(t, tt, "Achievement shine")
	find(t, tt, "First ship")
	panics(t, "Achievement name", func() { ui.Render(func(c *ui.Context) { Achievement(c, "", "trophy", false) }, 10, 10, 1) })
}

func TestCounter(t *testing.T) {
	target := 100
	tt := ui.NewTester(func(c *ui.Context) { Counter(c, target) }, 300, 200)
	tt.Advance(300 * time.Millisecond)
	if tt.HasText("0") || tt.HasText("100") {
		t.Fatalf("counter did not count: %q", tt.Texts())
	}
	tt.Advance(2 * time.Second)
	find(t, tt, "100")
	target = 200
	tt.Frame()
	settles(t, tt, "Counter target")
	find(t, tt, "200")
}

func TestLiveClock(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { LiveClock(c) }, 300, 200)
	expected := func(at time.Time) string {
		return i18n.FormatDate("en", at, i18n.DateOptions{Time: i18n.MediumStyle})
	}
	find(t, tt, expected(tt.Now()))
	tt.Advance(time.Second)
	find(t, tt, expected(tt.Now()))
}

func TestCountdownTo(t *testing.T) {
	target := time.Now().Add(65 * time.Second)
	tt := ui.NewTester(func(c *ui.Context) { CountdownTo(c, target) }, 400, 200)
	tt.Frame()
	if !tt.HasText("0d 00h 01m 05s") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Advance(2 * time.Second)
	if !tt.HasText("0d 00h 01m 03s") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestWeather(t *testing.T) {
	condition := WeatherSunny
	temperature := 22.5
	tt := ui.NewTester(func(c *ui.Context) { Weather(c, condition, temperature, "°C") }, 300, 200)
	condition = WeatherRainy
	temperature = 18
	tt.Frame()
	settles(t, tt, "Weather")
	find(t, tt, "18°C")
	panics(t, "Weather condition", func() {
		ui.Render(func(c *ui.Context) { Weather(c, WeatherCondition(9), 20, "°C") }, 10, 10, 1)
	})
}

func TestIDCard(t *testing.T) {
	fields := []KeyValuePair{{Key: "Role", Value: "Engineer"}, {Key: "Team", Value: "Interface"}}
	tt := ui.NewTester(func(c *ui.Context) { IDCard(c, "Ada Lovelace", nil, fields) }, 500, 300)
	settles(t, tt, "IDCard")
	find(t, tt, "Ada Lovelace")
	find(t, tt, "Engineer")
	panics(t, "IDCard fields", func() { ui.Render(func(c *ui.Context) { IDCard(c, "Ada", nil, nil) }, 10, 10, 1) })
}

func TestCountdownFormatsDays(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { CountdownTo(c, c.Now().Add(49*time.Hour)) }, 400, 200)
	tt.Frame()
	if !tt.HasText("2d 01h 00m 00s") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMetricRejectsNaN(t *testing.T) {
	panics(t, "Metric NaN", func() {
		ui.Render(func(c *ui.Context) { Metric(c, math.NaN(), "ms") }, 10, 10, 1)
	})
}
