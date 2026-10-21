package typography

import (
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

func TestRelativeTime(t *testing.T) {
	when := time.Now().Add(-3 * time.Minute)
	tt := ui.NewTester(func(c *ui.Context) { RelativeTime(c, when, RelativeTimeOptions{}) }, 300, 200)
	find(t, tt, "3 min ago")
	tt.Advance(4 * time.Minute)
	find(t, tt, "7 min ago")
	custom := ui.NewTester(func(c *ui.Context) {
		RelativeTime(c, when, RelativeTimeOptions{Minute: "分钟", AgoSuffix: " 前"})
	}, 300, 200)
	if !custom.HasText("分钟") {
		t.Errorf("custom units %q", custom.Texts())
	}
}

func TestRelativeTimeBoundaries(t *testing.T) {
	t.Run("future unit", func(t *testing.T) {
		when := time.Now().Add(65 * time.Minute)
		tt := ui.NewTester(func(c *ui.Context) { RelativeTime(c, when, RelativeTimeOptions{}) }, 300, 200)
		find(t, tt, "in 1 hr")
		tt.Advance(10 * time.Minute)
		if !tt.HasText("in 55 min") {
			t.Errorf("texts %q", tt.Texts())
		}
	})
	t.Run("past unit", func(t *testing.T) {
		when := time.Now().Add(-61 * time.Minute)
		tt := ui.NewTester(func(c *ui.Context) { RelativeTime(c, when, RelativeTimeOptions{}) }, 300, 200)
		find(t, tt, "1 hr ago")
		tt.Advance(59 * time.Minute)
		if !tt.HasText("2 hr ago") {
			t.Errorf("texts %q", tt.Texts())
		}
	})
}

func TestTimestamp(t *testing.T) {
	when := time.Date(2026, time.January, 5, 13, 7, 0, 0, time.UTC)
	tt := ui.NewTester(func(c *ui.Context) { Timestamp(c, when) }, 500, 200)
	display := func() string {
		return i18n.FormatDate("en", when, i18n.DateOptions{Date: i18n.MediumStyle, Time: i18n.ShortStyle})
	}
	find(t, tt, display())
	when = when.Add(24 * time.Hour)
	tt.Frame()
	settles(t, tt, "Timestamp")
	r := find(t, tt, display())
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(time.Second)
	full := i18n.FormatDate("en", when, i18n.DateOptions{Date: i18n.FullStyle, Time: i18n.LongStyle})
	if !tt.HasText(full) {
		t.Errorf("tooltip %q, want %q", tt.Texts(), full)
	}
}
