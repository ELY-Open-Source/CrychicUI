package code

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestRegexTesterClearsMatchedText(t *testing.T) {
	pattern, text := "a+", "banana"
	tt := ui.NewTester(padded(func(c *ui.Context) { RegexTester(c, "Rx", &pattern, &text, RegexTesterOptions{}) }), 500, 500)
	settles(t, tt, "regex")
	click(t, tt, "Text")
	tt.Key(ui.Ctrl, ui.KeyA)
	tt.Key(0, ui.KeyBackspace)
	tt.Frame()
	if text != "" || !tt.HasText("0 matches") {
		t.Errorf("text %q texts %q", text, tt.Texts())
	}
}

func TestJSONFormatKeepsActionFocus(t *testing.T) {
	text := "{"
	tt := ui.NewTester(padded(func(c *ui.Context) { JSONFormat(c, "invoice.json", &text, JSONFormatOptions{}) }), 500, 400)
	settles(t, tt, "json")
	click(t, tt, "JSON")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Format") {
		t.Fatal("Tab does not reach Format")
	}
	before := baseline()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if !tt.Focused("Format") {
		t.Error("the error row took Format's focus")
	}
	cues(t, before, sound.Error, sound.Error)
}

func TestBreadcrumbKeepsIdentityAcrossPaths(t *testing.T) {
	parts := []string{"billing", "internal", "invoice"}
	at := 2
	tt := ui.NewTester(padded(func(c *ui.Context) { Breadcrumb(c, "Path", parts, &at) }), 500, 200)
	settles(t, tt, "crumbs")
	swapUnderPress(tt, find(t, tt, "internal"), func() { parts = []string{"billing", "invoice", "close.go"} })
	if at != 2 {
		t.Errorf("a new crumb under the pointer fired: at %d", at)
	}
}

func TestBreadcrumbBoundsLongParts(t *testing.T) {
	long := strings.Repeat("segment", 20)
	at := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(320).Padding(10).Children(func() { Breadcrumb(c, "Path", []string{long, "x.go"}, &at) })
	}, 600, 300)
	settles(t, tt, "long crumb")
	if r := find(t, tt, long); r.X+r.W > 330 {
		t.Errorf("crumb runs to %v past its 330 box", r.X+r.W)
	}
}

func TestSinceText(t *testing.T) {
	now := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	for at, want := range map[time.Time]string{
		time.Date(2500, 1, 1, 0, 0, 0, 0, time.UTC):   "in 473y",
		time.Date(1500, 1, 1, 0, 0, 0, 0, time.UTC):   "526y ago",
		time.Date(2025, 10, 25, 9, 0, 0, 0, time.UTC): "11mo ago",
		time.Date(2025, 10, 24, 9, 0, 0, 0, time.UTC): "1y ago",
		now.Add(30 * time.Second):                     "now",
		now.Add(-3 * time.Hour):                       "3h ago",
	} {
		if got := sinceText(now, at); got != want {
			t.Errorf("sinceText(%v) = %q, want %q", at, got, want)
		}
	}
}

func TestUnixTimeFarAndNow(t *testing.T) {
	epoch := int64(16725225600)
	tt := ui.NewTester(padded(func(c *ui.Context) { UnixTime(c, "Unix", &epoch, UnixTimeOptions{}) }), 500, 300)
	settles(t, tt, "far")
	if !tt.HasText("Friday, 1 January 2500") || tt.HasText("292y") || !tt.HasText("y") {
		t.Errorf("far future texts %q", tt.Texts())
	}
	click(t, tt, "Now")
	settles(t, tt, "now")
	if !tt.HasText("now") || tt.HasText("now ago") {
		t.Errorf("present texts %q", tt.Texts())
	}
}
