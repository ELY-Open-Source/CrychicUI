package form

import (
	"hash/fnv"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestFormCard(t *testing.T) {
	var saved, cancelled int
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		_, submit, cancel := FormCard(c, "Account", FormCardOptions{Description: "Your details"}, func() {
			ui.Text(c, "body")
		})
		if submit.Clicked() {
			saved++
		}
		if cancel.Clicked() {
			cancelled++
		}
	}, 400, 300)
	if tt.FramePending() == false {
		t.Error("card did not rise in")
	}
	settles(t, tt, "FormCard")
	find(t, tt, "Your details")
	find(t, tt, "body")
	click(t, tt, "Save")
	click(t, tt, "Cancel")
	if saved != 1 || cancelled != 1 {
		t.Errorf("saved %d, cancelled %d", saved, cancelled)
	}
	cues(t, before, sound.Tap, sound.Tap)
	panics(t, "nil content", func() { render(func(c *ui.Context) { FormCard(c, "A", FormCardOptions{}, nil) }) })
}

func TestFormSection(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		FormSection(c, "Billing", FormSectionOptions{Description: "Where we send invoices"}, func() { ui.Text(c, "body") })
	}, 400, 200)
	if !tt.FramePending() {
		t.Error("section did not fade in")
	}
	settles(t, tt, "FormSection")
	if find(t, tt, "Where we send invoices").Y >= find(t, tt, "body").Y {
		t.Error("description not above the fields")
	}
	panics(t, "empty title", func() { render(func(c *ui.Context) { FormSection(c, "", FormSectionOptions{}, func() {}) }) })
}

func TestFormRowStacksWhenNarrow(t *testing.T) {
	view := func(c *ui.Context) {
		FormRow(c, []func(){func() { ui.Text(c, "first") }, func() { ui.Text(c, "last") }})
	}
	wide := ui.NewTester(view, 800, 200)
	wide.Advance(time.Second)
	if a, b := find(t, wide, "first"), find(t, wide, "last"); a.Y != b.Y || a.X >= b.X {
		t.Errorf("wide: %v and %v not side by side", a, b)
	}
	narrow := ui.NewTester(view, 300, 200)
	narrow.Advance(time.Second)
	if a, b := find(t, narrow, "first"), find(t, narrow, "last"); a.Y >= b.Y {
		t.Errorf("narrow: %v and %v not stacked", a, b)
	}
	panics(t, "no fields", func() { render(func(c *ui.Context) { FormRow(c, nil) }) })
}

func TestFieldGroup(t *testing.T) {
	var name, problem string
	tt := ui.NewTester(func(c *ui.Context) {
		FieldGroup(c, "Name", FieldGroupOptions{Help: "As on your card", Error: problem}, func() {
			ui.TextInput(c, &name)
		})
	}, 400, 200)
	find(t, tt, "As on your card")
	typeInto(t, tt, "Name", "Ada")
	if name != "Ada" {
		t.Fatalf("label click did not focus the input: %q", name)
	}
	problem = "Too short"
	tt.Frame()
	find(t, tt, "Too short")
	var frames []uint64
	for range 15 {
		tt.Advance(20 * time.Millisecond)
		frames = append(frames, shot(tt))
	}
	settles(t, tt, "FieldGroup")
	if !slices.ContainsFunc(frames, func(h uint64) bool { return h != shot(tt) }) {
		t.Error("error did not shake")
	}
	panics(t, "nil control", func() { render(func(c *ui.Context) { FieldGroup(c, "A", FieldGroupOptions{}, nil) }) })
}

// shot hashes the last frame.
func shot(tt *ui.Tester) uint64 {
	h := fnv.New64a()
	h.Write(tt.Image().Pix)
	return h.Sum64()
}

func TestFieldBoxDisablesItsFields(t *testing.T) {
	off, clicks := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		FieldBox(c, "Shipping", FieldBoxOptions{Disabled: off}, func() {
			if ui.Button(c, "Inside").Clicked() {
				clicks++
			}
		})
	}, 400, 200)
	find(t, tt, "Shipping")
	click(t, tt, "Inside")
	off = false
	tt.Frame()
	click(t, tt, "Inside")
	if clicks != 1 {
		t.Errorf("%d clicks; the disabled one should not count", clicks)
	}
	settles(t, tt, "FieldBox")
}

func TestStepHeader(t *testing.T) {
	at := 0
	tt := ui.NewTester(func(c *ui.Context) { StepHeader(c, []string{"Cart", "Address", "Pay"}, at) }, 500, 100)
	find(t, tt, "1")
	find(t, tt, "Pay")
	settles(t, tt, "StepHeader")
	at = 1
	tt.Frame()
	if !tt.FramePending() {
		t.Error("no pop or fill on a new step")
	}
	settles(t, tt, "StepHeader step")
	if tt.HasText("1") {
		t.Error("finished step still numbered; want a check")
	}
	panics(t, "step past the end", func() { render(func(c *ui.Context) { StepHeader(c, []string{"A"}, 1) }) })
}

func TestFieldGroupShakesOnFirstError(t *testing.T) {
	var name string
	tt := ui.NewTester(func(c *ui.Context) {
		FieldGroup(c, "Name", FieldGroupOptions{Error: "Required"}, func() { ui.TextInput(c, &name) })
	}, 400, 200)
	var frames []uint64
	for range 15 {
		tt.Advance(20 * time.Millisecond)
		frames = append(frames, shot(tt))
	}
	settles(t, tt, "FieldGroup first error")
	if !slices.ContainsFunc(frames, func(h uint64) bool { return h != shot(tt) }) {
		t.Error("first error did not shake")
	}
}
