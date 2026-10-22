package feedback

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestOfflineAndMaintenance(t *testing.T) {
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			if _, r := OfflineState(c, OfflineStateOptions{}); r.Clicked() {
				retries++
			}
			Maintenance(c, MaintenanceOptions{Title: "Back soon"})
		})
	}, 400, 600)
	find(t, tt, "You are offline")
	find(t, tt, "Back soon")
	loops(t, tt, "pulse and rock")
	before := baseline()
	click(t, tt, "Retry")
	if retries != 1 {
		t.Fatalf("retries %d", retries)
	}
	cues(t, before, sound.Tap)
}

func TestLoadingVeilAndBusy(t *testing.T) {
	loading := false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			LoadingVeil(c, loading, LoadingVeilOptions{}, func() { ui.Text(c, "Report").Padding(40) })
			Busy(c, loading, "Saving draft", "Draft saved")
		})
	}, 400, 300)
	find(t, tt, "Draft saved")
	tt.Advance(time.Second)
	if tt.FramePending() || tt.HasText("Loading") {
		t.Fatal("idle veil animates or shows")
	}
	loading = true
	tt.Frame()
	find(t, tt, "Loading")
	find(t, tt, "Saving draft")
	loops(t, tt, "veil spinner")
	loading = false
	tt.Frame()
	settles(t, tt, "veil fade out")
	if tt.HasText("Loading") {
		t.Error("veil stayed")
	}
}

func TestUploadProgress(t *testing.T) {
	ups := []Upload{{"a.png", 0.2}, {"b.pdf", 1}}
	tt := ui.NewTester(func(c *ui.Context) { UploadProgress(c, ups) }, 400, 300)
	find(t, tt, "20%")
	find(t, tt, "100%")
	ups = ups[1:]
	tt.Frame()
	settles(t, tt, "row exit")
	if tt.HasText("a.png") {
		t.Error("removed row stayed")
	}
	panics(t, "repeated name", func() { render(func(c *ui.Context) { UploadProgress(c, []Upload{{"a", 0}, {"a", 0}}) }) })
}

func TestStepCount(t *testing.T) {
	step := 2
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(300).Children(func() { StepCount(c, step, 5, StepCountOptions{}) })
	}, 300, 120)
	find(t, tt, "Step 2 of 5")
	step = 3
	tt.Frame()
	find(t, tt, "Step 3 of 5")
	settles(t, tt, "step fill")
	panics(t, "step past total", func() { render(func(c *ui.Context) { StepCount(c, 6, 5, StepCountOptions{}) }) })
}

func TestTimeLeftBar(t *testing.T) {
	var ends time.Time
	tt := ui.NewTester(func(c *ui.Context) {
		if ends.IsZero() {
			ends = c.Now().Add(2 * time.Second)
		}
		ui.Column(c).Width(400).Children(func() { TimeLeftBar(c, "Session", ends, 2*time.Second) })
	}, 400, 40)
	th := ui.LightTheme()
	r := find(t, tt, "Session")
	at := func(f float32, k ui.Color) bool { return tinted(tt, int(r.X+r.W*f), int(r.Y+r.H/2), k) }
	if !at(0.9, th.Accent) {
		t.Error("fresh bar not full")
	}
	tt.Advance(time.Second)
	if !at(0.25, th.Accent) || at(0.75, th.Accent) {
		t.Errorf("half drained: 25%%=%v 75%%=%v", at(0.25, th.Accent), at(0.75, th.Accent))
	}
	if !tt.FramePending() {
		t.Error("bar stopped draining")
	}
	tt.Advance(2 * time.Second)
	if tt.FramePending() || !at(0.05, th.Border) {
		t.Errorf("drained bar: pending %v, empty %v", tt.FramePending(), at(0.05, th.Border))
	}
	panics(t, "zero total", func() { render(func(c *ui.Context) { TimeLeftBar(c, "x", time.Time{}, 0) }) })
}

func TestRetryIn(t *testing.T) {
	var at time.Time
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if at.IsZero() {
			at = c.Now().Add(2500 * time.Millisecond)
		}
		if RetryIn(c, at, RetryInOptions{}).Clicked() {
			clicks++
		}
	}, 300, 100)
	find(t, tt, "Retry in 3 s")
	tt.Advance(600 * time.Millisecond)
	find(t, tt, "Retry in 2 s")
	if err := tt.Click("Retry in 2 s"); err == nil && clicks > 0 {
		t.Error("disabled button clicked")
	}
	tt.Advance(2 * time.Second)
	before := baseline()
	click(t, tt, "Retry")
	if clicks != 1 {
		t.Fatalf("clicks %d", clicks)
	}
	cues(t, before, sound.Tap)
	far := ui.NewTester(func(c *ui.Context) { RetryIn(c, c.Now().AddDate(300, 0, 0), RetryInOptions{}) }, 300, 100)
	if _, ok := far.Find("Retry in 9223372037 s"); !ok {
		t.Errorf("far deadline shows %q", far.Texts())
	}
}

func TestAttentionKeepsContent(t *testing.T) {
	on, v := false, ""
	tt := ui.NewTester(func(c *ui.Context) {
		Attention(c, on, func() { ui.TextInput(c, &v).Label("Name") })
	}, 300, 100)
	click(t, tt, "Name")
	for _, now := range []bool{false, true, false} {
		on = now
		tt.Frame()
		if !tt.Focused("Name") {
			t.Fatalf("on=%v: content lost focus", now)
		}
	}
}

func TestErrorListAndFieldError(t *testing.T) {
	var name, email ui.Handle
	a, b := "", ""
	msg := ""
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			ui.TextInput(c.Key("name"), &a).Label("Name").Bind(&name)
			ui.TextInput(c.Key("email"), &b).Label("Email").Bind(&email)
			ErrorList(c, []FormError{{"Name is required", &name}, {"Email is invalid", &email}}, ErrorListOptions{})
			FieldError(c, msg)
		})
	}, 400, 400)
	before := baseline()
	click(t, tt, "Email is invalid")
	tt.Frame()
	if !tt.Focused("Email") {
		t.Error("click did not focus the field")
	}
	cues(t, before, sound.Tap)
	msg = "Too short"
	tt.Frame()
	r, moved := find(t, tt, "Too short"), false
	for range 8 {
		tt.Advance(50 * time.Millisecond)
		moved = moved || find(t, tt, "Too short").X != r.X
	}
	if !moved {
		t.Error("field error did not shake")
	}
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Error("shake never settled")
	}
	panics(t, "nil field", func() {
		render(func(c *ui.Context) { ErrorList(c, []FormError{{"x", nil}}, ErrorListOptions{}) })
	})
}

func TestHintAttentionBeacon(t *testing.T) {
	on := true
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(12).Padding(20).Children(func() {
			Hint(c, "Use 8 characters or more.")
			Attention(c, on, func() { ui.Button(c, "New") })
			Beacon(c, "What is new", "Charts now export to PNG.")
		})
	}, 400, 400)
	find(t, tt, "Use 8 characters or more.")
	loops(t, tt, "attention and beacon")
	before := baseline()
	click(t, tt, "What is new")
	find(t, tt, "Charts now export to PNG.")
	cues(t, before, sound.Open)
	if !strings.Contains(strings.Join(tt.Texts(), " "), "New") {
		t.Error("attention lost its content")
	}
}
