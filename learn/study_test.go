package learn

import (
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func clayQuestion() Question {
	return Question{Text: "Which clay body matures lowest?", Choices: []string{"Earthenware", "Stoneware", "Porcelain"}, Answer: 0,
		Explain: "Earthenware matures near cone 06."}
}

func TestFlashcardFlips(t *testing.T) {
	up := false
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Flashcard(c, "Cone 6 reaches?", "About 1,222 °C", &up, FlashcardOptions{}).Changed() {
			changed++
		}
	}), 400, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Cone 6 reaches?")
	cues(t, before, sound.Swipe)
	tt.Advance(100 * time.Millisecond)
	if exact(tt, "Answer") {
		t.Error("the answer shows before the card is edge-on")
	}
	settles(t, tt)
	if !up || changed != 1 || !exact(tt, "About 1,222 °C") || !exact(tt, "Answer") {
		t.Errorf("up %v changed %d: %q", up, changed, tt.Texts())
	}
	click(t, tt, "About 1,222 °C")
	tt.Advance(100 * time.Millisecond)
	click(t, tt, "Cone 6 reaches?")
	for range 6 {
		if exact(tt, "Question") {
			t.Fatal("a flip turned back before edge-on showed the question")
		}
		tt.Advance(60 * time.Millisecond)
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	settles(t, tt)
	click(t, tt, "About 1,222 °C")
	if !exact(tt, "Question") {
		t.Error("with Reduce Motion the question waits")
	}
	panics(t, "needs a front", func(c *ui.Context) { Flashcard(c, " ", "b", &up, FlashcardOptions{}) })
}

func TestQuizChecks(t *testing.T) {
	picked, checked := -1, false
	tt := ui.NewTester(padded(func(c *ui.Context) { Quiz(c, clayQuestion(), &picked, &checked, QuizOptions{}) }), 420, 500)
	settles(t, tt)
	_ = tt.Click("Check")
	tt.Frame()
	if checked {
		t.Fatal("checked with nothing picked")
	}
	before := baseline()
	click(t, tt, "Stoneware")
	click(t, tt, "Check")
	cues(t, before, sound.Select, sound.Error)
	settles(t, tt)
	if !checked || !exact(tt, "Not quite") || !exact(tt, "Earthenware matures near cone 06.") {
		t.Errorf("wrong pick: %q", tt.Texts())
	}
	_ = tt.Click("Porcelain")
	tt.Frame()
	if picked != 1 {
		t.Error("a checked quiz took a new pick")
	}
	picked, checked = -1, false
	tt.Frame()
	before = baseline()
	click(t, tt, "Earthenware")
	click(t, tt, "Check")
	cues(t, before, sound.Select, sound.Success)
	if !exact(tt, "Correct") {
		t.Errorf("right pick: %q", tt.Texts())
	}
	panics(t, "answer among them", func(c *ui.Context) {
		q := clayQuestion()
		q.Answer = 3
		p, k := -1, false
		Quiz(c, q, &p, &k, QuizOptions{})
	})
	panics(t, "one before checking", func(c *ui.Context) { p, k := -1, true; Quiz(c, clayQuestion(), &p, &k, QuizOptions{}) })
}

func TestQuizResultScores(t *testing.T) {
	qs := []Question{clayQuestion(), {Text: "First firing?", Choices: []string{"Glaze", "Bisque"}, Answer: 1},
		{Text: "Hottest cone?", Choices: []string{"06", "6", "10"}, Answer: 2}}
	tt := ui.NewTester(padded(func(c *ui.Context) { QuizResult(c, qs, []int{0, 0, -1}, QuizResultOptions{}) }), 420, 600)
	settles(t, tt)
	for _, s := range []string{"33%", "1 of 3 correct", "Keep practicing", "You chose Glaze", "Answer: Bisque", "Skipped", "Answer: 10"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { QuizResult(c, qs, []int{0, 1, 0}, QuizResultOptions{Pass: 0.6}) }), 420, 600)
	settles(t, tt)
	if !exact(tt, "Passed") || !exact(tt, "67%") {
		t.Errorf("two of three at a 0.6 pass: %q", tt.Texts())
	}
	panics(t, "pass share", func(c *ui.Context) { QuizResult(c, qs, []int{0, 0, 0}, QuizResultOptions{Pass: math.NaN()}) })
	panics(t, "out of range", func(c *ui.Context) { QuizResult(c, qs, []int{0, 0, 3}, QuizResultOptions{}) })
}

func TestPercentEdges(t *testing.T) {
	for _, c := range [][3]int{{199, 200, 99}, {1, 1000, 1}, {0, 5, 0}, {5, 5, 100}, {2, 3, 67}} {
		if got := percent(c[0], c[1]); got != c[2] {
			t.Errorf("percent(%d, %d) = %d, want %d", c[0], c[1], got, c[2])
		}
	}
}

func clayLessons() []Lesson {
	return []Lesson{{Title: "Wedging", Icon: "hand", Minutes: 12, Done: true}, {Title: "Centering", Icon: "circle-dot", Minutes: 18, Done: true},
		{Title: "Pulling walls", Icon: "move-vertical", Minutes: 22}, {Title: "Glazing", Icon: "paintbrush", Minutes: 25}}
}

func TestLessonPathPicks(t *testing.T) {
	lessons := clayLessons()
	picked := 2
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(300).Children(func() {
			if LessonPath(c, lessons, &picked, LessonPathOptions{}).Changed() {
				changed++
			}
		})
	}), 320, 600)
	settles(t, tt)
	if !tt.FramePending() {
		t.Error("the current lesson does not pulse")
	}
	before := baseline()
	_ = tt.Click("Glazing, locked")
	tt.Frame()
	click(t, tt, "Centering")
	cues(t, before, sound.Select)
	if picked != 1 || changed != 1 {
		t.Errorf("picked %d changed %d", picked, changed)
	}
	for _, s := range []string{"Wedging", "Centering", "Pulling walls", "Glazing", "Lesson 4 · 25 min"} {
		inside(t, s, rectOf(t, tt, s), 10, 310)
	}
	panics(t, "done after an open one", func(c *ui.Context) {
		l := clayLessons()
		l[3].Done = true
		p := -1
		LessonPath(c, l, &p, LessonPathOptions{})
	})
	panics(t, "repeated", func(c *ui.Context) {
		l := clayLessons()
		l[1].Title = "Wedging"
		p := -1
		LessonPath(c, l, &p, LessonPathOptions{})
	})
}

func TestCourseCardText(t *testing.T) {
	k := Course{Title: "Wheel throwing basics", Author: "Priya Nair", Lessons: 12, Done: 5, Minutes: 200}
	tt := ui.NewTester(padded(func(c *ui.Context) { CourseCard(c, k, CourseCardOptions{}) }), 360, 500)
	settles(t, tt)
	for _, s := range []string{"42% complete", "12 lessons", "3h 20m"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	k = Course{Title: "Glaze tests", Author: "Daniel Okafor", Lessons: 1, Done: 1, Minutes: 45}
	tt = ui.NewTester(padded(func(c *ui.Context) { CourseCard(c, k, CourseCardOptions{}) }), 360, 500)
	settles(t, tt)
	for _, s := range []string{"Completed", "1 lesson", "45 min"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	panics(t, "done <= lessons", func(c *ui.Context) {
		CourseCard(c, Course{Title: "x", Author: "y", Lessons: 2, Done: 3}, CourseCardOptions{})
	})
}

func TestStudyStreakWeek(t *testing.T) {
	days := []StudyDay{{Label: "Mon", Minutes: 25}, {Label: "Tue", Minutes: 40}, {Label: "Wed"}, {Label: "Thu", Minutes: 30},
		{Label: "Fri", Minutes: 22}, {Label: "Sat", Minutes: 12, Today: true}, {Label: "Sun"}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(300).Children(func() { StudyStreak(c, days, 3, StudyStreakOptions{}) })
	}), 320, 400)
	settles(t, tt)
	if !exact(tt, "129 min this week") || !exact(tt, "3-day streak") {
		t.Errorf("header: %q", tt.Texts())
	}
	for _, s := range []string{"Mon", "Sun", "40 min"} {
		inside(t, s, rectOf(t, tt, s), 10, 310)
	}
	panics(t, "marked today", func(c *ui.Context) {
		StudyStreak(c, []StudyDay{{Label: "Mon", Today: true}, {Label: "Tue", Today: true}}, 0, StudyStreakOptions{})
	})
	panics(t, "minutes", func(c *ui.Context) {
		StudyStreak(c, []StudyDay{{Label: "Mon", Minutes: math.MaxInt}, {Label: "Tue", Minutes: 1}}, 0, StudyStreakOptions{})
	})
}

func clayWords() []Word {
	return []Word{{Term: "Raku", Reading: "ra-ku", Meaning: "Low-fired ware.", Example: "Pulled glowing from the kiln.", Tone: 440},
		{Term: "Tenmoku", Meaning: "An iron glaze.", Example: "Breaks rust to black.", Tone: 494}, {Term: "Wabi-sabi", Meaning: "Worn beauty."}}
}

func TestVocabularyOpensAndPlays(t *testing.T) {
	words := clayWords()
	open, played, changed := -1, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Vocabulary(c, words, &open, VocabularyOptions{Played: &played}).Changed() {
			changed++
		}
	}), 420, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Tenmoku")
	settles(t, tt)
	if open != 1 || !exact(tt, "Breaks rust to black.") || exact(tt, "Pulled glowing from the kiln.") {
		t.Errorf("open %d: %q", open, tt.Texts())
	}
	click(t, tt, "Listen to Raku")
	cues(t, before, sound.Tap, sound.Note)
	if played != 1 || changed != 2 {
		t.Errorf("played %d changed %d", played, changed)
	}
	if _, ok := tt.Find("Listen to Wabi-sabi"); ok {
		t.Error("a word without a tone has a listen button")
	}
	swapUnderPress(t, tt, "Raku", func() { words[0], words[1] = words[1], words[0] })
	if open != 1 {
		t.Errorf("a reorder under the press opened row %d", open)
	}
	panics(t, "tone of 0 or 20", func(c *ui.Context) {
		o := -1
		Vocabulary(c, []Word{{Term: "a", Meaning: "b", Tone: 10}}, &o, VocabularyOptions{})
	})
	panics(t, "repeated", func(c *ui.Context) {
		o := -1
		Vocabulary(c, []Word{{Term: "a", Meaning: "b"}, {Term: "a", Meaning: "c"}}, &o, VocabularyOptions{})
	})
}

func TestFillBlankChecks(t *testing.T) {
	text, checked := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Width(400).Children(func() {
			FillBlank(c, "Bisque firing hardens the pot so", "can go on.", "glaze", &text, &checked, FillBlankOptions{})
		})
	}), 420, 400)
	settles(t, tt)
	_ = tt.Click("Check")
	tt.Frame()
	if checked {
		t.Fatal("checked a blank gap")
	}
	click(t, tt, "Missing word")
	tt.Type("wax")
	tt.Frame()
	before := baseline()
	click(t, tt, "Check")
	cues(t, before, sound.Error)
	if !checked || !exact(tt, "The answer is glaze") {
		t.Errorf("wrong word: %q", tt.Texts())
	}
	click(t, tt, "Missing word")
	tt.Key(0, ui.KeyEnd)
	for range 3 {
		tt.Key(0, ui.KeyBackspace)
	}
	tt.Type(" Glaze ")
	tt.Frame()
	if checked {
		t.Error("an edit kept the old verdict")
	}
	before = baseline()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	cues(t, before, sound.Success)
	if !checked || !exact(tt, "Correct") {
		t.Errorf("Enter on the right word: %q", tt.Texts())
	}
}
