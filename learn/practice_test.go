package learn

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func toolPairs() []Pair {
	return []Pair{{Left: "Rib", Right: "Smooths walls"}, {Left: "Wire", Right: "Cuts it free"}, {Left: "Sponge", Right: "Lifts water"}}
}

func TestMatchPairsLinks(t *testing.T) {
	matched := make([]bool, 3)
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if MatchPairs(c, toolPairs(), []int{2, 0, 1}, &matched, MatchPairsOptions{}).Changed() {
			changed++
		}
	}), 460, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Rib")
	click(t, tt, "Cuts it free")
	cues(t, before, sound.Tap, sound.Error)
	if slices.Contains(matched, true) || changed != 0 {
		t.Fatalf("a wrong pair matched: %v", matched)
	}
	before = baseline()
	click(t, tt, "Smooths walls")
	click(t, tt, "Rib")
	cues(t, before, sound.Tap, sound.Success)
	settles(t, tt)
	if !matched[0] || changed != 1 || !exact(tt, "1 of 3 matched") {
		t.Errorf("right side first: %v %q", matched, tt.Texts())
	}
	click(t, tt, "Wire")
	click(t, tt, "Cuts it free")
	click(t, tt, "Sponge")
	click(t, tt, "Lifts water")
	settles(t, tt)
	if !exact(tt, "All matched") {
		t.Errorf("all matched: %q", tt.Texts())
	}
	panics(t, "permutation", func(c *ui.Context) {
		m := make([]bool, 3)
		MatchPairs(c, toolPairs(), []int{0, 0, 1}, &m, MatchPairsOptions{})
	})
}

func TestMatchPairsHeldPick(t *testing.T) {
	pairs := toolPairs()
	matched := make([]bool, 3)
	picks := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MatchPairs(c, pairs, []int{2, 0, 1}, &matched, MatchPairsOptions{Picks: &picks})
	}), 460, 400)
	settles(t, tt)
	click(t, tt, "Rib")
	pairs[0], pairs[1] = pairs[1], pairs[0]
	tt.Frame()
	before := baseline()
	click(t, tt, "Smooths walls")
	cues(t, before, sound.Success)
	if !matched[1] || picks != 2 {
		t.Errorf("a reorder lost the held pick: matched %v picks %d", matched, picks)
	}
}

func TestMatchPairsLineDrawsAfterReorder(t *testing.T) {
	pairs := toolPairs()
	matched := make([]bool, 3)
	tt := ui.NewTester(padded(func(c *ui.Context) { MatchPairs(c, pairs, []int{2, 0, 1}, &matched, MatchPairsOptions{}) }), 460, 400)
	settles(t, tt)
	click(t, tt, "Rib")
	click(t, tt, "Smooths walls")
	settles(t, tt)
	pairs[0], pairs[1] = pairs[1], pairs[0]
	matched[0], matched[1] = matched[1], matched[0]
	tt.Frame()
	settles(t, tt)
	l, r := rectOf(t, tt, "Wire"), rectOf(t, tt, "Cuts it free")
	// green counts the link pixels between the two columns.
	green := func() int {
		img, n := tt.Image(), 0
		for y := int(l.Y - 60); y < int(l.Y+120); y++ {
			for x := int(l.X+l.W) + 2; x < int(r.X)-2; x++ {
				if p := img.RGBAAt(x, y); int(p.G)-int(p.R) > 40 {
					n++
				}
			}
		}
		return n
	}
	one := green()
	click(t, tt, "Wire")
	click(t, tt, "Cuts it free")
	tt.Advance(60 * time.Millisecond)
	early := green()
	settles(t, tt)
	if full := green(); early-one >= (full-one)*3/4 {
		t.Errorf("the new line was %d of %d pixels drawn at 60 ms", early-one, full-one)
	}
}

func TestReportsAndIdentity(t *testing.T) {
	text, checked, changed := "", false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if FillBlank(c, "Bisque hardens the pot so", "can go on.", "glaze", &text, &checked, FillBlankOptions{}).Changed() {
			changed++
		}
	}), 420, 300)
	settles(t, tt)
	click(t, tt, "Missing word")
	tt.Type("g")
	tt.Frame()
	tt.Type("l")
	tt.Frame()
	if changed != 2 {
		t.Errorf("FillBlank reported %d of 2 edits", changed)
	}
	q := clayQuestion()
	picked, done := -1, false
	tt = ui.NewTester(padded(func(c *ui.Context) { Quiz(c, q, &picked, &done, QuizOptions{}) }), 420, 500)
	settles(t, tt)
	swapUnderPress(t, tt, "Stoneware", func() { q.Choices[0], q.Choices[1] = q.Choices[1], q.Choices[0]; q.Answer = 1 })
	if picked != -1 {
		t.Errorf("a held click picked %q, which slid under it", q.Choices[picked])
	}
	followed := false
	k := Teacher{Name: "Priya Nair", Rating: 4.9}
	tt = ui.NewTester(padded(func(c *ui.Context) { Instructor(c, k, &followed, InstructorOptions{}) }), 400, 500)
	settles(t, tt)
	click(t, tt, "Follow")
	k.Bio = "Priya teaches beginners."
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if followed {
		t.Error("Follow lost focus when the bio appeared")
	}
}

func TestTrueFalseOnce(t *testing.T) {
	said := Unsaid
	tt := ui.NewTester(padded(func(c *ui.Context) {
		TrueFalse(c, "Porcelain fires hotter than earthenware.", true, &said, TrueFalseOptions{Explain: "Cone 10 against cone 06."})
	}), 420, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "False")
	_ = tt.Click("True")
	tt.Frame()
	cues(t, before, sound.Error)
	if said != SaidFalse || !exact(tt, "Not quite") || !exact(tt, "Cone 10 against cone 06.") {
		t.Errorf("said %d: %q", said, tt.Texts())
	}
	panics(t, "answer state", func(c *ui.Context) { s := Said(7); TrueFalse(c, "x", true, &s, TrueFalseOptions{}) })
}

func TestOrderStepsDrop(t *testing.T) {
	correct := []string{"Wedge", "Center", "Open", "Pull"}
	steps := []string{"Center", "Wedge", "Open", "Pull"}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if OrderSteps(c, &steps, correct, OrderStepsOptions{}).Changed() {
			changed++
		}
	}), 420, 500)
	settles(t, tt)
	if !exact(tt, "2 of 4 in place") {
		t.Errorf("start: %q", tt.Texts())
	}
	before := baseline()
	drag(tt, rectOf(t, tt, "Wedge"), rectOf(t, tt, "Center"))
	cues(t, before, sound.Drop)
	settles(t, tt)
	if !slices.Equal(steps, correct) || changed != 1 || !exact(tt, "In order") {
		t.Errorf("steps %v: %q", steps, tt.Texts())
	}
	panics(t, "not the answer's steps", func(c *ui.Context) {
		s := []string{"Wedge", "Glaze", "Open", "Pull"}
		OrderSteps(c, &s, correct, OrderStepsOptions{})
	})
}

func TestHintRevealOpens(t *testing.T) {
	shown := false
	tt := ui.NewTester(padded(func(c *ui.Context) { HintReveal(c, "Brace your elbows.", &shown, HintRevealOptions{}) }), 400, 400)
	settles(t, tt)
	if exact(tt, "Brace your elbows.") {
		t.Fatal("the hint shows before it is asked for")
	}
	before := baseline()
	click(t, tt, "Show hint")
	cues(t, before, sound.Open)
	if !shown || !exact(tt, "Brace your elbows.") {
		t.Errorf("shown %v: %q", shown, tt.Texts())
	}
}

func TestCertificateDate(t *testing.T) {
	a := Award{Name: "Maya Chen", Course: "Wheel Throwing", Issuer: "Clay School", Date: time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)}
	tt := ui.NewTester(padded(func(c *ui.Context) { Certificate(c, a, CertificateOptions{}) }), 460, 600)
	settles(t, tt)
	if !exact(tt, "October 18, 2026") || !exact(tt, "Maya Chen") {
		t.Errorf("certificate: %q", tt.Texts())
	}
	panics(t, "date", func(c *ui.Context) { Certificate(c, Award{Name: "a", Course: "b", Issuer: "c"}, CertificateOptions{}) })
}

func TestGradeBands(t *testing.T) {
	for _, c := range []struct {
		score         float64
		letter, shown string
	}{{89.99, "B", "89.9%"}, {90, "A", "90%"}, {59.5, "F", "59.5%"}, {100, "A", "100%"}} {
		tt := ui.NewTester(padded(func(ctx *ui.Context) { Grade(ctx, c.score, GradeOptions{}) }), 300, 200)
		settles(t, tt)
		if !exact(tt, c.letter) || !exact(tt, c.shown) {
			t.Errorf("%v: %q", c.score, tt.Texts())
		}
	}
	panics(t, "0 to 100", func(c *ui.Context) { Grade(c, math.NaN(), GradeOptions{}) })
	panics(t, "0 to 100", func(c *ui.Context) { Grade(c, 100.5, GradeOptions{}) })
}

func TestStudyTimerSessions(t *testing.T) {
	sessions := 0
	changed, runs, skips := 0, 0, 0
	o := StudyTimerOptions{Focus: time.Minute, Rest: 30 * time.Second}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, run, skip := StudyTimer(c, &sessions, o)
		if root.Changed() {
			changed++
		}
		if run.Clicked() {
			runs++
		}
		if skip.Clicked() {
			skips++
		}
	}), 400, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Start")
	tt.Advance(61 * time.Second)
	cues(t, before, sound.Tap, sound.Notify)
	if sessions != 1 || changed != 1 || runs != 1 || !exact(tt, "Break") {
		t.Errorf("after a focus: sessions %d changed %d runs %d %q", sessions, changed, runs, tt.Texts())
	}
	tt.Advance(31 * time.Second)
	if !exact(tt, "Focus") || !exact(tt, "1:00") || !exact(tt, "Start") {
		t.Errorf("after the break: %q", tt.Texts())
	}
	click(t, tt, "Start")
	tt.Advance(20 * time.Second)
	click(t, tt, "Pause")
	tt.Advance(time.Minute)
	if !exact(tt, "0:40") {
		t.Errorf("paused at 0:40: %q", tt.Texts())
	}
	click(t, tt, "Skip")
	if sessions != 1 || skips != 1 || !exact(tt, "Break") {
		t.Errorf("a skipped focus counted: sessions %d skips %d", sessions, skips)
	}
	panics(t, "24 h", func(c *ui.Context) { s := 0; StudyTimer(c, &s, StudyTimerOptions{Focus: 25 * time.Hour}) })
}

func TestNotesCounts(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Notes(c, nil, NotesOptions{}) }), 420, 300)
	settles(t, tt)
	if !exact(tt, "No highlights yet") {
		t.Errorf("empty: %q", tt.Texts())
	}
	tt = ui.NewTester(padded(func(c *ui.Context) { Notes(c, []Highlight{{Quote: "Patience applied evenly."}}, NotesOptions{}) }), 420, 300)
	settles(t, tt)
	if !exact(tt, "1 highlight") {
		t.Errorf("one: %q", tt.Texts())
	}
	panics(t, "tint", func(c *ui.Context) { Notes(c, []Highlight{{Quote: "a", Tint: 4}}, NotesOptions{}) })
}

func TestSyllabusFolds(t *testing.T) {
	weeks := []Week{{Title: "Wedging", Topics: []string{"Clay bodies"}, Done: true}, {Title: "Centering", Topics: []string{"Speed", "Coning"}},
		{Title: "Cylinders", Topics: []string{"Opening"}}}
	open := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Syllabus(c, weeks, &open, SyllabusOptions{}) }), 420, 500)
	settles(t, tt)
	if !exact(tt, "This week") || !exact(tt, "1 topic") || !exact(tt, "2 topics") {
		t.Errorf("weeks: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Centering")
	settles(t, tt)
	click(t, tt, "Centering")
	cues(t, before, sound.Open, sound.Close)
	if open != -1 {
		t.Errorf("open %d after two clicks", open)
	}
	click(t, tt, "Centering")
	settles(t, tt)
	if !exact(tt, "Coning") {
		t.Errorf("open week shows no topics: %q", tt.Texts())
	}
	panics(t, "done after an open one", func(c *ui.Context) {
		o := -1
		Syllabus(c, []Week{{Title: "a", Topics: []string{"x"}}, {Title: "b", Topics: []string{"y"}, Done: true}}, &o, SyllabusOptions{})
	})
}

func TestInstructorFollows(t *testing.T) {
	followed := false
	k := Teacher{Name: "Priya Nair", Role: "Ceramicist", Rating: 4.9, Students: 12480, Courses: 6}
	tt := ui.NewTester(padded(func(c *ui.Context) { Instructor(c, k, &followed, InstructorOptions{}) }), 400, 500)
	settles(t, tt)
	if !exact(tt, "12,480") || !exact(tt, "4.9") {
		t.Errorf("stats: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Follow")
	cues(t, before, sound.Toggle)
	if !followed || !exact(tt, "Following") {
		t.Errorf("followed %v: %q", followed, tt.Texts())
	}
	panics(t, "rating of 0 to 5", func(c *ui.Context) {
		f := false
		Instructor(c, Teacher{Name: "a", Rating: 5.1}, &f, InstructorOptions{})
	})
}

func TestRecallGrades(t *testing.T) {
	grade, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Recall(c, [4]string{"<1m", "6m", "1d", "4d"}, &grade, RecallOptions{}).Changed() {
			changed++
		}
	}), 420, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Good")
	click(t, tt, "Good")
	cues(t, before, sound.Select)
	if grade != 2 || changed != 1 {
		t.Errorf("grade %d changed %d", grade, changed)
	}
	panics(t, "next-review time", func(c *ui.Context) { g := -1; Recall(c, [4]string{"a", "", "c", "d"}, &g, RecallOptions{}) })
}

func TestProgressTotals(t *testing.T) {
	modules := []Module{{Title: "Basics", Done: 4, Total: 4}, {Title: "Forms", Done: 2, Total: 5}, {Title: "Glazing", Total: 4}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Progress(c, modules, ProgressOptions{}) }), 420, 400)
	settles(t, tt)
	for _, s := range []string{"46%", "6 of 13 lessons done", "2 of 5"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	panics(t, "total that fits", func(c *ui.Context) {
		Progress(c, []Module{{Title: "a", Total: math.MaxInt}, {Title: "b", Total: 1}}, ProgressOptions{})
	})
}
