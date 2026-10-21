package display

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func TestWatermarkDoesNotBlockContent(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Watermark(c, "DRAFT", func() {
			if ui.Button(c, "Inside").Size(160, 80).Clicked() {
				clicks++
			}
		})
	}, 300, 200)
	tt.ClickAt(25, 10)
	if clicks != 1 {
		t.Error("watermark blocked the content click")
	}
}

func TestLevelLeavesRoomDuringFill(t *testing.T) {
	progress := 0.75
	var accent ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		accent = c.Theme().Accent
		Level(c, 2, progress, LevelOptions{})
	}, 300, 200)
	tt.Frame()
	progress = 1
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	if near(tt.Image(), 295, 25, accent) {
		t.Error("remaining space disappeared before the fill arrived")
	}
	tt.Advance(2 * time.Second)
	if !near(tt.Image(), 295, 25, accent) {
		t.Error("fill never completed")
	}
}

func TestCounterKeepsLargeTargets(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Counter(c, 16777217) }, 400, 200)
	tt.Advance(2 * time.Second)
	if !tt.HasText("16777217") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCounterKeepsHugeTargets(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Counter(c, 9007199254740993) }, 500, 200)
	tt.Advance(2 * time.Second)
	if !tt.HasText("9007199254740993") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestPriceRejectsNaNOldPrice(t *testing.T) {
	panics(t, "Price old NaN", func() {
		ui.Render(func(c *ui.Context) { Price(c, 10, math.NaN(), "USD") }, 10, 10, 1)
	})
}

func TestLiveClockSecondBoundary(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { LiveClock(c) }, 300, 200)
	at := tt.Now()
	shown := i18n.FormatDate("en", at, i18n.DateOptions{Time: i18n.MediumStyle})
	wait := time.Second - time.Duration(at.Nanosecond())
	tt.Advance(wait)
	next := i18n.FormatDate("en", tt.Now(), i18n.DateOptions{Time: i18n.MediumStyle})
	if next == shown || !tt.HasText(next) {
		t.Errorf("clock stayed %q, want %q", shown, next)
	}
}

func TestCountdownTargetFraction(t *testing.T) {
	target := time.Now().Add(60*time.Second + 250*time.Millisecond)
	builds := 0
	tt := ui.NewTester(func(c *ui.Context) {
		builds++
		CountdownTo(c, target)
	}, 400, 200)
	tt.Frame()
	built := builds
	wait := target.Sub(tt.Now())%time.Second + time.Millisecond
	tt.Advance(wait)
	if !tt.HasText("0d 00h 00m 59s") {
		t.Errorf("wait %v, remaining %v, builds %d, texts %q", wait, target.Sub(tt.Now()), builds-built, tt.Texts())
	}
}

func TestCountdownWholeSecondBoundary(t *testing.T) {
	var target time.Time
	builds := 0
	tt := ui.NewTester(func(c *ui.Context) {
		builds++
		if target.IsZero() {
			target = c.Now().Add(65 * time.Second)
		}
		CountdownTo(c, target)
	}, 400, 200)
	built := builds
	tt.Advance(time.Second)
	if !tt.HasText("0d 00h 01m 04s") {
		t.Errorf("builds %d, remaining %v, texts %q", builds-built, target.Sub(tt.Now()), tt.Texts())
	}
}
