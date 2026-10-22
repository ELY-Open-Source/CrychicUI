package form

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestActionBar(t *testing.T) {
	dirty := false
	saved, cancelled := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Size(400, 300).Children(func() {
			_, save, cancel := ActionBar(c, dirty, ActionBarOptions{}, func() { ui.Text(c, "content") })
			if save.Clicked() {
				saved++
			}
			if cancel.Clicked() {
				cancelled++
			}
		})
	}, 400, 300)
	settles(t, tt, "ActionBar hidden")
	find(t, tt, "content")
	click(t, tt, "Save")
	if saved != 0 {
		t.Fatal("hidden bar saved")
	}
	dirty = true
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bar did not slide")
	}
	settles(t, tt, "ActionBar shown")
	before := baseline()
	click(t, tt, "Save")
	click(t, tt, "Cancel")
	if saved != 1 || cancelled != 1 {
		t.Errorf("saved %d, cancelled %d", saved, cancelled)
	}
	cues(t, before, sound.Success, sound.Tap)
}

func TestDirtyBar(t *testing.T) {
	dirty := true
	discarded, saved := 0, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		_, discard, save := DirtyBar(c, dirty, DirtyBarOptions{})
		if discard.Clicked() {
			discarded++
		}
		if save.Clicked() {
			saved++
		}
	}, 400, 120)
	settles(t, tt, "DirtyBar")
	find(t, tt, "Unsaved changes")
	click(t, tt, "Discard")
	click(t, tt, "Save")
	dirty = false
	tt.Frame()
	settles(t, tt, "DirtyBar hidden")
	click(t, tt, "Save")
	if discarded != 1 || saved != 1 {
		t.Errorf("discarded %d, saved %d", discarded, saved)
	}
	cues(t, before, sound.Close, sound.Success)
}

func TestErrorSummary(t *testing.T) {
	var name, email ui.Handle
	a, b := "", ""
	var errs []feedback.FormError
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			ErrorSummary(c, errs, ErrorSummaryOptions{})
			ui.TextInput(c.Key("name"), &a).Label("Name").Bind(&name)
			ui.TextInput(c.Key("email"), &b).Label("Email").Bind(&email)
		})
	}, 400, 400)
	if tt.HasText("Fix 0 problems") {
		t.Error("summary shown without errors")
	}
	cues(t, before)
	errs = []feedback.FormError{{Message: "Name is required", Field: &name}, {Message: "Email is invalid", Field: &email}}
	tt.Frame()
	settles(t, tt, "ErrorSummary")
	find(t, tt, "Fix 2 problems")
	click(t, tt, "Email is invalid")
	if !tt.Focused("Email") {
		t.Error("click did not focus the field")
	}
	cues(t, before, sound.Error, sound.Tap)
}

func TestHelpPopover(t *testing.T) {
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() { HelpPopover(c, "About the slug", "The slug is the last part of the link.") })
	}, 400, 300)
	click(t, tt, "About the slug")
	settles(t, tt, "HelpPopover")
	find(t, tt, "The slug is the last part of the link.")
	tt.Key(0, ui.KeyEscape)
	settles(t, tt, "HelpPopover closed")
	if tt.HasText("The slug is the last part of the link.") {
		t.Error("Escape left it open")
	}
	click(t, tt, "About the slug")
	settles(t, tt, "HelpPopover again")
	find(t, tt, "The slug is the last part of the link.")
	cues(t, before, sound.Open, sound.Open)
	panics(t, "no help", func() { render(func(c *ui.Context) { HelpPopover(c, "A", "") }) })
}

func TestRequiredLegend(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { RequiredLegend(c, RequiredLegendOptions{Text: "Needed"}) }, 300, 60)
	if !tt.FramePending() {
		t.Error("legend did not fade in")
	}
	settles(t, tt, "RequiredLegend")
	find(t, tt, "Needed")
	find(t, tt, "*")
}

func TestAutosaveForm(t *testing.T) {
	var text string
	var edits uint64
	saves := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if AutosaveForm(c, &edits, AutosaveFormOptions{}, func() {
			if ui.TextInput(c, &text).Label("Note").Changed() {
				edits++
			}
		}).Changed() {
			saves++
		}
	}, 400, 200)
	if tt.HasText("Saved") || tt.HasText("Saving…") {
		t.Error("status shown before any edit")
	}
	typeInto(t, tt, "Note", "hi")
	tt.Advance(500 * time.Millisecond)
	typeInto(t, tt, "Note", "!")
	tt.Advance(700 * time.Millisecond)
	if saves != 0 {
		t.Fatal("saved before the pause ended")
	}
	tt.Advance(time.Second)
	settles(t, tt, "AutosaveForm")
	if saves != 1 {
		t.Errorf("saved %d times, want once", saves)
	}
	panics(t, "negative delay", func() {
		render(func(c *ui.Context) { AutosaveForm(c, &edits, AutosaveFormOptions{Delay: -1}, func() {}) })
	})
}
