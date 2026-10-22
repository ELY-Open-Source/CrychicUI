package form

import (
	"fmt"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestInlineForm(t *testing.T) {
	var value string
	sent, changed := 0, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		root, send := InlineForm(c, "Email", &value, InlineFormOptions{})
		if root.Changed() {
			changed++
		}
		if send.Clicked() {
			sent++
		}
	}, 400, 100)
	click(t, tt, "Submit")
	typeInto(t, tt, "Email", "  ")
	click(t, tt, "Submit")
	if sent != 0 {
		t.Fatal("blank value was sent")
	}
	cues(t, before)
	typeInto(t, tt, "Email", "ada@example.com")
	if changed == 0 {
		t.Error("typing did not report a change")
	}
	click(t, tt, "Submit")
	if sent != 1 {
		t.Fatalf("sent %d times", sent)
	}
	cues(t, before, sound.Send)
}

func TestWizard(t *testing.T) {
	at, finished := 0, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		_, finish := Wizard(c, []string{"One", "Two"}, &at, func(i int) { ui.Text(c, fmt.Sprintf("page %d", i)) },
			WizardOptions{})
		if finish.Clicked() {
			finished++
		}
	}, 500, 300)
	find(t, tt, "page 0")
	click(t, tt, "Next")
	settles(t, tt, "Wizard")
	if at != 1 || !tt.HasText("page 1") || tt.HasText("page 0") {
		t.Fatalf("at %d, texts %q", at, tt.Texts())
	}
	click(t, tt, "Finish")
	if finished != 1 {
		t.Errorf("finished %d", finished)
	}
	click(t, tt, "Back")
	settles(t, tt, "Wizard back")
	if at != 0 || !tt.HasText("page 0") {
		t.Errorf("back went to %d", at)
	}
	cues(t, before, sound.Tap, sound.Success, sound.Tap)
	panics(t, "nil pages", func() {
		render(func(c *ui.Context) { Wizard(c, []string{"A"}, &at, nil, WizardOptions{}) })
	})
}

func TestSurvey(t *testing.T) {
	answers := []int{-1, -1}
	qs := []Question{{"Tabs or spaces?", []string{"Tabs", "Spaces"}}, {"Light or dark?", []string{"Light", "Dark"}}}
	before := baseline()
	backs := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if _, back := Survey(c, qs, &answers, SurveyOptions{}); back.Clicked() {
			backs++
		}
	}, 400, 300)
	find(t, tt, "1 of 2")
	click(t, tt, "Spaces")
	settles(t, tt, "Survey")
	if answers[0] != 1 || !tt.HasText("Light or dark?") || tt.HasText("Tabs or spaces?") {
		t.Fatalf("answers %v, texts %q", answers, tt.Texts())
	}
	click(t, tt, "Back")
	settles(t, tt, "Survey back")
	find(t, tt, "Tabs or spaces?")
	if backs != 1 {
		t.Errorf("back reported %d times", backs)
	}
	click(t, tt, "Spaces")
	click(t, tt, "Dark")
	settles(t, tt, "Survey end")
	find(t, tt, "Thanks for your answers.")
	if answers[1] != 1 {
		t.Errorf("answers %v", answers)
	}
	cues(t, before, sound.Select, sound.Tap, sound.Select, sound.Select)
	short := []int{-1}
	panics(t, "answer slots", func() { render(func(c *ui.Context) { Survey(c, qs, &short, SurveyOptions{}) }) })
	bad := []int{5, -1}
	panics(t, "answer out of range", func() { render(func(c *ui.Context) { Survey(c, qs, &bad, SurveyOptions{}) }) })
}
