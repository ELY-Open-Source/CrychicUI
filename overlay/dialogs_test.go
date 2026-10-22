package overlay

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestModal(t *testing.T) {
	open, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Modal(c, "Settings", &open, ModalOptions{}, func() { ui.Text(c, "Body") }).Changed() {
			changed++
		}
	}, 600, 400)
	before := baseline()
	open = true
	tt.Frame()
	if !tt.HasText("Body") || !tt.FramePending() {
		t.Fatalf("did not scale in: %q", tt.Texts())
	}
	settles(t, tt, "Modal")
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || changed != 1 || tt.HasText("Body") {
		t.Fatalf("Escape: open %v, changed %d", open, changed)
	}
	open = true
	tt.Frame()
	click(t, tt, "Close")
	if open || changed != 2 {
		t.Errorf("close button: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Open, sound.Close, sound.Open, sound.Close)
	panics(t, "no content", func() { render(func(c *ui.Context) { o := true; Modal(c, "x", &o, ModalOptions{}, nil) }) })
}

func TestConfirmAndDangerDialogs(t *testing.T) {
	open, yes, changed := true, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		root, confirm := ConfirmDialog(c, "Publish?", "Everyone will see it.", &open, ConfirmDialogOptions{})
		if root.Changed() {
			changed++
		}
		if confirm.Clicked() {
			yes++
		}
	}, 600, 400)
	click(t, tt, "Confirm")
	if open || yes != 1 {
		t.Fatalf("confirm: open %v, confirmed %d", open, yes)
	}
	before := baseline()
	open = true
	tt.Frame()
	click(t, tt, "Cancel")
	if open || yes != 1 || changed != 1 {
		t.Errorf("cancel: open %v, confirmed %d, changed %d", open, yes, changed)
	}
	cues(t, before, sound.Open, sound.Tap)

	open, yes = false, 0
	td := ui.NewTester(func(c *ui.Context) {
		if _, confirm := DangerDialog(c, "Delete project?", "This cannot be undone.", &open, DangerDialogOptions{}); confirm.Clicked() {
			yes++
		}
	}, 600, 400)
	before = baseline()
	open = true
	td.Frame()
	click(t, td, "Delete")
	if open || yes != 1 {
		t.Errorf("danger: open %v, confirmed %d", open, yes)
	}
	cues(t, before, sound.Warning, sound.Delete)
}

func TestPromptDialog(t *testing.T) {
	value, open, changed := "old", false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if root, _ := PromptDialog(c, "Rename", &value, &open, PromptDialogOptions{}); root.Changed() {
			changed++
		}
	}, 600, 400)
	open = true
	tt.Frame()
	before := baseline()
	tt.Command("selectAll")
	tt.Type("Ada")
	tt.Frame()
	if value != "old" || changed != 0 {
		t.Fatalf("typing reached the value: %q, changed %d", value, changed)
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if value != "Ada" || open || changed != 1 {
		t.Fatalf("Enter: value %q, open %v, changed %d", value, open, changed)
	}
	cues(t, before, sound.Success)
	open = true
	tt.Frame()
	tt.Command("selectAll")
	tt.Type(" ")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if !open || value != "Ada" {
		t.Errorf("blank submitted: open %v, value %q", open, value)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || value != "Ada" || changed != 2 {
		t.Errorf("Escape: open %v, value %q, changed %d", open, value, changed)
	}
}

func TestFullscreenModal(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		FullscreenModal(c, "Editor", &open, FullscreenModalOptions{}, func() { ui.Text(c, "Canvas") })
	}, 600, 400)
	before := baseline()
	open = true
	tt.Frame()
	if !tt.HasText("Canvas") || !tt.FramePending() {
		t.Fatalf("did not slide up: %q", tt.Texts())
	}
	settles(t, tt, "FullscreenModal")
	if r := find(t, tt, "Editor"); r.W < 500 {
		t.Errorf("does not fill the window: %+v", r)
	}
	click(t, tt, "Close")
	if open {
		t.Error("still open")
	}
	cues(t, before, sound.Open)
}

func TestWhatsNew(t *testing.T) {
	page, open, changed := 0, true, 0
	fs := []Feature{{"sparkles", "Themes", "Pick a look."}, {"zap", "Speed", "Twice as fast."}, {"bell", "Alerts", "Hear about it."}}
	tt := ui.NewTester(func(c *ui.Context) {
		if WhatsNew(c, "What's new", fs, &page, &open, WhatsNewOptions{}).Changed() {
			changed++
		}
	}, 600, 400)
	before := baseline()
	click(t, tt, "Next")
	if page != 1 || !tt.FramePending() {
		t.Fatalf("page %d, sliding %v", page, tt.FramePending())
	}
	click(t, tt, "Back")
	click(t, tt, "Next")
	click(t, tt, "Next")
	if page != 2 || tt.HasText("Next") {
		t.Fatalf("last page %d, %q", page, tt.Texts())
	}
	click(t, tt, "Done")
	if open || changed != 5 {
		t.Errorf("done: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Swipe, sound.Swipe, sound.Swipe, sound.Swipe, sound.Tap)
	panics(t, "page past the end", func() {
		render(func(c *ui.Context) { p, o := 3, true; WhatsNew(c, "x", fs, &p, &o, WhatsNewOptions{}) })
	})
}
