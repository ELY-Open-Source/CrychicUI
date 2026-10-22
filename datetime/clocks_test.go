package datetime

import (
	"hash/fnv"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// shot hashes the last frame.
func shot(tt *ui.Tester) uint64 {
	h := fnv.New64a()
	h.Write(tt.Image().Pix)
	return h.Sum64()
}

func TestAnalogClockSweeps(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { AnalogClock(c, AnalogClockOptions{}) }, 200, 200)
	a := shot(tt)
	tt.Advance(250 * time.Millisecond)
	if shot(tt) == a || !tt.FramePending() {
		t.Error("second hand did not sweep")
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	b := shot(tt)
	tt.Advance(time.Duration(time.Second-time.Duration(tt.Now().Nanosecond())).Truncate(time.Millisecond) - 50*time.Millisecond)
	if shot(tt) != b {
		t.Error("hand moved within a second under reduced motion")
	}
	tt.Advance(time.Second)
	if shot(tt) == b {
		t.Error("hand did not tick under reduced motion")
	}
	panics(t, "negative size", func() { render(func(c *ui.Context) { AnalogClock(c, AnalogClockOptions{Size: -1}) }) })
}

func TestDigitalClock(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { DigitalClock(c, DigitalClockOptions{Seconds: true}) }, 400, 100)
	first := tt.Texts()
	var frames []uint64
	for range 4 {
		tt.Advance(500 * time.Millisecond)
		frames = append(frames, shot(tt))
	}
	if frames[0] == frames[2] && frames[1] == frames[3] && frames[0] == frames[1] {
		t.Error("colon did not blink")
	}
	tt.Advance(time.Minute)
	if slices.Equal(tt.Texts(), first) {
		t.Errorf("digits stayed %q", first)
	}
}

func TestWorldClock(t *testing.T) {
	ahead := time.FixedZone("ahead", 0)
	tt := ui.NewTester(func(c *ui.Context) {
		_, here := c.Now().Zone()
		ahead = time.FixedZone("ahead", here+3*3600)
		WorldClock(c, []Zone{{"Home", c.Now().Location()}, {"East", ahead}}, WorldClockOptions{})
	}, 400, 200)
	onMinute(tt)
	settles(t, tt, "WorldClock")
	find(t, tt, "Same time")
	find(t, tt, "3 h ahead")
	panics(t, "nil location", func() { render(func(c *ui.Context) { WorldClock(c, []Zone{{"X", nil}}, WorldClockOptions{}) }) })
}

func TestStopwatch(t *testing.T) {
	before := baseline()
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Stopwatch(c, StopwatchOptions{}).Changed() {
			changes++
		}
	}, 400, 400)
	find(t, tt, "00:00.00")
	click(t, tt, "Start")
	tt.Advance(1500 * time.Millisecond)
	click(t, tt, "Lap 1")
	tt.Advance(500 * time.Millisecond)
	click(t, tt, "Stop")
	settles(t, tt, "Stopwatch")
	find(t, tt, "00:02.00")
	find(t, tt, "00:01.50")
	click(t, tt, "Reset")
	settles(t, tt, "Stopwatch reset")
	find(t, tt, "00:00.00")
	if tt.HasText("00:01.50") {
		t.Error("laps stayed after reset")
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	if changes != 4 {
		t.Errorf("%d changes, want 4", changes)
	}
}

func TestCountdown(t *testing.T) {
	before := baseline()
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Countdown(c, 3*time.Second, CountdownOptions{}).Changed() {
			changes++
		}
	}, 400, 400)
	find(t, tt, "0:03")
	click(t, tt, "Start")
	tt.Advance(1200 * time.Millisecond)
	find(t, tt, "0:02")
	tt.Advance(2 * time.Second)
	settles(t, tt, "Countdown done")
	find(t, tt, "0:00")
	click(t, tt, "Reset")
	settles(t, tt, "Countdown reset")
	find(t, tt, "0:03")
	cues(t, before, sound.Tap, sound.Warning, sound.Tap)
	if changes != 3 {
		t.Errorf("%d changes, want 3", changes)
	}
	panics(t, "zero total", func() { render(func(c *ui.Context) { Countdown(c, 0, CountdownOptions{}) }) })
}

func TestPomodoro(t *testing.T) {
	before := baseline()
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Pomodoro(c, PomodoroOptions{Work: 2 * time.Second, Break: time.Second}).Changed() {
			changes++
		}
	}, 400, 400)
	find(t, tt, "Focus")
	click(t, tt, "Start")
	tt.Advance(2100 * time.Millisecond)
	tt.Advance(100 * time.Millisecond)
	find(t, tt, "Break")
	tt.Advance(time.Second)
	tt.Advance(100 * time.Millisecond)
	find(t, tt, "Focus")
	cues(t, before, sound.Tap, sound.Notify, sound.Notify)
	if changes != 3 {
		t.Errorf("%d changes, want 3", changes)
	}
}

func TestDigitalClockFollowsTheHourCycle(t *testing.T) {
	var am, pm, written string
	tt := ui.NewTester(func(c *ui.Context) {
		am, pm = i18n.DayPeriods(locale(c), i18n.Abbreviated)
		written = i18n.FormatDate(locale(c), c.Now(), i18n.DateOptions{Skeleton: "jm"})
		DigitalClock(c, DigitalClockOptions{})
	}, 400, 100)
	tt.SetLocale("en-US")
	tt.Frame()
	formats := strings.Contains(written, am) || strings.Contains(written, pm)
	if !formats {
		t.Fatalf("en-US wrote %q without a day period", written)
	}
	_, a := tt.Find(am)
	_, p := tt.Find(pm)
	shows := a || p
	if formats != shows {
		t.Errorf("MyGo writes %q; the clock shows a day period: %v", written, shows)
	}
}
