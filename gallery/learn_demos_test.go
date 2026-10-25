package gallery

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestLearnDemosLog(t *testing.T) {
	press := func(names ...string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, name := range names {
				if err := tt.Click(name); err != nil {
					return err
				}
				tt.Frame()
			}
			return nil
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Flashcard":  {press("What temperature does cone 6 reach?"), "answer shown true"},
		"Quiz":       {press("Stoneware", "Check"), "choice 1, checked false\nchoice 1, checked true"},
		"LessonPath": {press("Centering"), "lesson 2"},
		"Vocabulary": {press("Raku", "Listen to Raku"), "open 1, tones 0\nopen 1, tones 1"},
		"FillBlank": {func(tt *ui.Tester) error {
			if err := press("Missing word")(tt); err != nil {
				return err
			}
			tt.Type("glaze")
			tt.Frame()
			return press("Check")(tt)
		}, "checked false, 5 runes\nchecked true, 5 runes"},
		"MatchPairs": {press("Rib", "Cuts the pot free", "Rib", "Smooths the walls"),
			"picks 1, matched 0\npicks 2, matched 0\npicks 3, matched 0\npicks 4, matched 1"},
		"TrueFalse": {press("True"), "said 1"},
		"OrderSteps": {func(tt *ui.Tester) error {
			a, ok := tt.Find("Wedge the clay")
			if !ok {
				return fmt.Errorf("no step")
			}
			b, _ := tt.Find("Open the form")
			tt.Press(a.X+a.W/2, a.Y+a.H/2)
			tt.Move(b.X+b.W/2, b.Y+b.H/2)
			tt.Release(b.X+b.W/2, b.Y+b.H/2)
			tt.Frame()
			return nil
		}, "moved, 2 in place"},
		"HintReveal": {press("Show hint"), "hint shown true"},
		"StudyTimer": {press("Start", "Skip"), "timer toggled\nphase skipped"},
		"Syllabus":   {press("Bowls"), "open week 4"},
		"Instructor": {press("Follow"), "following true"},
		"Recall":     {press("Good"), "grade 2"},
	}
	quiet := map[string]bool{"QuizResult": true, "CourseCard": true, "StudyStreak": true, "Certificate": true, "Grade": true, "Notes": true,
		"Progress": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "learn" {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(2 * time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 20 {
		t.Errorf("checked %d learn demos", seen)
	}
}
