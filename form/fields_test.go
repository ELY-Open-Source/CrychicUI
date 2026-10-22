package form

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestInviteForm(t *testing.T) {
	var emails []string
	role, sent := "Member", 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, send := InviteForm(c, &emails, &role, []string{"Member", "Admin"}, InviteFormOptions{}); send.Clicked() {
			sent++
		}
	}, 500, 400)
	click(t, tt, "Send invites")
	tt.Key(0, ui.KeyTab)
	tt.Type("bad,")
	tt.Frame()
	find(t, tt, "Remove the addresses that are not emails.")
	click(t, tt, "Send invites")
	emails = []string{"ada@example.com"}
	tt.Frame()
	click(t, tt, "Send invites")
	if sent != 1 {
		t.Errorf("sent %d", sent)
	}
	cues(t, before, sound.Tick, sound.Pop, sound.Send)
}

func TestTypeToConfirm(t *testing.T) {
	var typed string
	confirmed := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, ok := TypeToConfirm(c, "prod-db", &typed, TypeToConfirmOptions{}); ok.Clicked() {
			confirmed++
		}
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	tt.Type("prod-d")
	tt.Frame()
	click(t, tt, "Delete")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Type("b")
	tt.Frame()
	click(t, tt, "Delete")
	if confirmed != 1 {
		t.Errorf("confirmed %d, typed %q", confirmed, typed)
	}
	cues(t, before, sound.Delete)
}

func TestDynamicFields(t *testing.T) {
	var fields []DynamicField
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		DynamicFields(c, &fields, []string{"Phone", "Website"}, DynamicFieldsOptions{})
	}, 500, 300)
	click(t, tt, "Add field")
	click(t, tt, "Add field")
	settles(t, tt, "DynamicFields add")
	if len(fields) != 2 || fields[0].Kind != "Phone" {
		t.Fatalf("fields %+v", fields)
	}
	click(t, tt, "Remove field")
	settles(t, tt, "DynamicFields remove")
	if len(fields) != 1 {
		t.Errorf("fields %+v", fields)
	}
	cues(t, before, sound.Pop, sound.Pop, sound.Delete)
	bad := []DynamicField{{Kind: "Fax"}}
	panics(t, "unknown kind", func() {
		render(func(c *ui.Context) { DynamicFields(c, &bad, []string{"Phone"}, DynamicFieldsOptions{}) })
	})
}

func TestConditionalField(t *testing.T) {
	on := false
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		ConditionalField(c, "Ship elsewhere", &on, func() { ui.Text(c, "Other address") })
	}, 400, 200)
	if tt.HasText("Other address") {
		t.Fatal("content shown while off")
	}
	click(t, tt, "Ship elsewhere")
	if !on || !tt.FramePending() {
		t.Fatalf("on %v, springing %v", on, tt.FramePending())
	}
	settles(t, tt, "ConditionalField")
	find(t, tt, "Other address")
	cues(t, before, sound.Toggle)
}

func TestReviewStep(t *testing.T) {
	edit := -1
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		ReviewStep(c, []ReviewSection{
			{"Contact", []display.KeyValuePair{{Key: "Name", Value: "Ada"}}},
			{"Shipping", []display.KeyValuePair{{Key: "City", Value: "London"}}},
		}, &edit, ReviewStepOptions{})
	}, 400, 300)
	find(t, tt, "London")
	click(t, tt, "Edit Shipping")
	if edit != 1 {
		t.Errorf("edit %d", edit)
	}
	cues(t, before, sound.Tap)
}

func TestConsent(t *testing.T) {
	agreed, went := false, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, submit := Consent(c, "I accept the terms", &agreed, ConsentOptions{}); submit.Clicked() {
			went++
		}
	}, 400, 200)
	click(t, tt, "Continue")
	click(t, tt, "I accept the terms")
	click(t, tt, "Continue")
	if !agreed || went != 1 {
		t.Errorf("agreed %v, went %d", agreed, went)
	}
	cues(t, before, sound.Toggle, sound.Tap)
}

func TestScheduleField(t *testing.T) {
	at := time.Date(2026, 3, 4, 9, 30, 0, 0, time.UTC)
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { ScheduleField(c, "Starts", &at, ScheduleFieldOptions{}) }, 500, 200)
	if !tt.FramePending() {
		t.Error("ScheduleField did not rise in")
	}
	settles(t, tt, "ScheduleField")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if at.Hour() != 10 || at.Day() != 4 {
		t.Errorf("at %v", at)
	}
	cues(t, before, sound.Select)
}

func TestNumberRange(t *testing.T) {
	low, high := 2.0, 4.0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) { NumberRange(c, "Price", &low, &high, 0, 10, 1, NumberRangeOptions{}) }, 500, 200)
	if !tt.FramePending() {
		t.Error("NumberRange did not rise in")
	}
	settles(t, tt, "NumberRange")
	tt.Key(0, ui.KeyTab)
	for range 4 {
		tt.Key(0, ui.KeyUp)
	}
	tt.Frame()
	if low != 6 || high != 6 {
		t.Errorf("low %v, high %v", low, high)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	nan := math.NaN()
	panics(t, "NaN step", func() {
		render(func(c *ui.Context) { NumberRange(c, "P", &low, &high, 0, 10, nan, NumberRangeOptions{}) })
	})
	low, high = 8, 2
	panics(t, "unordered values", func() {
		render(func(c *ui.Context) { NumberRange(c, "P", &low, &high, 0, 10, 1, NumberRangeOptions{}) })
	})
	low, high = 2, 4
	panics(t, "empty range", func() {
		render(func(c *ui.Context) { NumberRange(c, "P", &low, &high, 5, 5, 1, NumberRangeOptions{}) })
	})
}

func TestRepeater(t *testing.T) {
	type guest struct{ Name string }
	var guests []guest
	changed := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if Repeater(c, &guests, func(i int, g *guest) {
			ui.TextInput(c, &g.Name).Label(fmt.Sprintf("Guest %d", i+1))
		}, RepeaterOptions{}).Changed() {
			changed++
		}
	}, 500, 400)
	click(t, tt, "Add another")
	added := changed
	typeInto(t, tt, "Guest 1", "Ada")
	if changed <= added {
		t.Error("typing in a group reported no change")
	}
	click(t, tt, "Add another")
	typeInto(t, tt, "Guest 2", "Bob")
	click(t, tt, "Remove group 1")
	settles(t, tt, "Repeater")
	if len(guests) != 1 || guests[0].Name != "Bob" || !tt.HasText("Guest 1") || tt.HasText("Guest 2") {
		t.Errorf("guests %+v, texts %q", guests, tt.Texts())
	}
	cues(t, before, sound.Pop, sound.Pop, sound.Delete)
}

func TestQuestionBlock(t *testing.T) {
	var answer string
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if QuestionBlock(c, "What did you build?", &answer, QuestionBlockOptions{Hint: "One line is fine."}).Changed() {
			changed++
		}
	}, 400, 200)
	settles(t, tt, "QuestionBlock")
	find(t, tt, "One line is fine.")
	tt.Key(0, ui.KeyTab)
	tt.Type("A form kit")
	tt.Frame()
	if answer != "A form kit" || changed == 0 {
		t.Errorf("answer %q, changed %d", answer, changed)
	}
	panics(t, "negative max", func() {
		render(func(c *ui.Context) { QuestionBlock(c, "Q", &answer, QuestionBlockOptions{Max: -1}) })
	})
}

func TestFormResult(t *testing.T) {
	outcome, acted := Pending, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, action := FormResult(c, outcome, FormResultOptions{}, func() { ui.Text(c, "the form") }); action.Clicked() {
			acted++
		}
	}, 400, 400)
	find(t, tt, "the form")
	outcome = Failed
	tt.Frame()
	settles(t, tt, "FormResult failed")
	find(t, tt, "We could not send it. Try again.")
	click(t, tt, "Retry")
	outcome = Succeeded
	tt.Frame()
	settles(t, tt, "FormResult sent")
	if tt.HasText("the form") {
		t.Error("form still shown")
	}
	click(t, tt, "Done")
	if acted != 2 {
		t.Errorf("acted %d", acted)
	}
	cues(t, before, sound.Error, sound.Tap, sound.Success, sound.Tap)
	panics(t, "unknown outcome", func() { render(func(c *ui.Context) { FormResult(c, 9, FormResultOptions{}, func() {}) }) })
}
