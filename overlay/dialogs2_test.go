package overlay

import (
	"fmt"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestShortcutSheet(t *testing.T) {
	groups := []ShortcutGroup{{"Edit", []Shortcut{{"Copy", []string{"Ctrl", "C"}}}}, {"View", []Shortcut{{"Zoom in", []string{"Ctrl", "+"}}}}}
	open, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ShortcutSheet(c, "Shortcuts", groups, &open, ShortcutSheetOptions{}).Changed() {
			changed++
		}
	}, 700, 500)
	before := baseline()
	open = true
	tt.Frame()
	settles(t, tt, "sheet")
	if !tt.HasText("Zoom in") || !tt.HasText("Ctrl") {
		t.Fatalf("sheet %q", tt.Texts())
	}
	click(t, tt, "Close")
	if open || changed != 1 {
		t.Errorf("close: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Open, sound.Close)
	panics(t, "no keys", func() {
		render(func(c *ui.Context) {
			o := true
			ShortcutSheet(c, "x", []ShortcutGroup{{"A", []Shortcut{{"Copy", nil}}}}, &o, ShortcutSheetOptions{})
		})
	})
}

func TestStepDialog(t *testing.T) {
	step, open, changed := 0, true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if StepDialog(c, "Checkout", []string{"Name", "Plan", "Pay"}, &step, &open, func(i int) { ui.Text(c, fmt.Sprint("Page ", i)) },
			StepDialogOptions{Count: "Paso %d de %d"}).Changed() {
			changed++
		}
	}, 600, 500)
	settles(t, tt, "dialog")
	if !tt.HasText("Paso 1 de 3") || !tt.HasText("Page 0") {
		t.Fatalf("first step %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Back")
	click(t, tt, "Next")
	if step != 1 || !tt.FramePending() {
		t.Fatalf("Next: step %d, sliding %v", step, tt.FramePending())
	}
	settles(t, tt, "slide")
	click(t, tt, "Next")
	settles(t, tt, "last")
	click(t, tt, "Done")
	if open || changed != 3 {
		t.Errorf("Done: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
}

func TestRenameDialog(t *testing.T) {
	name, open, changed := "report.pdf", false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if root, _ := RenameDialog(c, "Rename file", &name, &open, RenameDialogOptions{}); root.Changed() {
			changed++
		}
	}, 600, 400)
	open = true
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if !open || name != "report.pdf" {
		t.Fatalf("unchanged Enter: open %v, name %q", open, name)
	}
	before := baseline()
	tt.Type("notes")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if name != "notes.pdf" || open || changed != 1 {
		t.Errorf("rename: name %q, open %v, changed %d", name, open, changed)
	}
	cues(t, before, sound.Success)
}

func TestSessionTimeout(t *testing.T) {
	deadline, open, stays, changed := time.Now().Add(5*time.Second), true, 0, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		root, stay, _ := SessionTimeout(c, deadline, &open, SessionTimeoutOptions{})
		if stay.Clicked() {
			stays++
		}
		if root.Changed() {
			changed++
		}
	}, 600, 400)
	if !tt.HasText("You will be signed out in 5 s.") {
		t.Fatalf("start %q", tt.Texts())
	}
	tt.Advance(1100 * time.Millisecond)
	if !tt.HasText("You will be signed out in 4 s.") {
		t.Fatalf("a second on %q", tt.Texts())
	}
	click(t, tt, "Stay signed in")
	if open || stays != 1 {
		t.Fatalf("stay: open %v, stays %d", open, stays)
	}
	cues(t, before, sound.Warning, sound.Tap)
	deadline, open = tt.Now().Add(2*time.Second), true
	tt.Frame()
	tt.Advance(3 * time.Second)
	if open || changed != 2 {
		t.Errorf("expiry: open %v, changed %d", open, changed)
	}
	panics(t, "no deadline", func() {
		render(func(c *ui.Context) { o := true; SessionTimeout(c, time.Time{}, &o, SessionTimeoutOptions{}) })
	})
}

func TestFeedbackWidget(t *testing.T) {
	text, open, sent := "", false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).FillWidth().FillHeight().Justify(ui.End).Padding(16).Children(func() {
			if _, send := FeedbackWidget(c, &text, &open, FeedbackWidgetOptions{}); send.Clicked() {
				sent++
			}
		})
	}, 600, 500)
	before := baseline()
	click(t, tt, "Feedback")
	settles(t, tt, "form")
	tt.Type("Great")
	tt.Frame()
	click(t, tt, "Send")
	if text != "Great" || sent != 1 || open {
		t.Fatalf("send: text %q, sent %d, open %v", text, sent, open)
	}
	cues(t, before, sound.Open, sound.Send)
	text, open = "  ", true
	tt.Frame()
	settles(t, tt, "reopened")
	click(t, tt, "Send")
	if sent != 1 || !open {
		t.Errorf("blank send: sent %d, open %v", sent, open)
	}
}

func TestAskReportsCloses(t *testing.T) {
	for _, rename := range []bool{false, true} {
		value, open, changed := "report.pdf", false, 0
		tt := ui.NewTester(func(c *ui.Context) {
			var root ui.Element
			if rename {
				root, _ = RenameDialog(c, "Rename file", &value, &open, RenameDialogOptions{})
			} else {
				root, _ = PromptDialog(c, "Rename file", &value, &open, PromptDialogOptions{})
			}
			if root.Changed() {
				changed++
			}
		}, 600, 400)
		closes := []func(){
			func() { click(t, tt, "Cancel") },
			func() { tt.Key(0, ui.KeyEscape); tt.Frame() },
			func() { tt.ClickAt(5, 5); tt.Frame() },
		}
		for i, shut := range closes {
			open = true
			tt.Frame()
			settles(t, tt, "dialog")
			shut()
			if open || changed != i+1 || value != "report.pdf" {
				t.Errorf("rename %v, close %d: open %v, changed %d, value %q", rename, i, open, changed, value)
			}
		}
	}
}
