package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/learn"
	"github.com/egoist/mygo/ui"
)

var learnEntries = []Entry{
	{ID: 1331, Name: "Flashcard", Category: "learn", Summary: "Card flips to the answer; cue Swipe.", Demo: demoLearnFlashcard},
	{ID: 1332, Name: "Quiz", Category: "learn", Summary: "Question with choices and check; cue Success, Error.", Demo: demoLearnQuiz},
	{ID: 1333, Name: "QuizResult", Category: "learn", Summary: "Score ring and review.", Demo: demoLearnQuizResult},
	{ID: 1334, Name: "LessonPath", Category: "learn", Summary: "Lessons on a winding path.", Demo: demoLearnLessonPath},
	{ID: 1335, Name: "CourseCard", Category: "learn", Summary: "Course with progress.", Demo: demoLearnCourseCard},
	{ID: 1336, Name: "StudyStreak", Category: "learn", Summary: "Week of study days.", Demo: demoLearnStudyStreak},
	{ID: 1337, Name: "Vocabulary", Category: "learn", Summary: "Words with meanings and audio tone; cue Tap.", Demo: demoLearnVocabulary},
	{ID: 1338, Name: "FillBlank", Category: "learn", Summary: "Sentence with an input gap; cue Success, Error.", Demo: demoLearnFillBlank},
	{ID: 1339, Name: "MatchPairs", Category: "learn", Summary: "Match left and right; lines draw; cue Success.", Demo: demoLearnMatchPairs},
	{ID: 1340, Name: "TrueFalse", Category: "learn", Summary: "Two choices with feedback; cue Success, Error.", Demo: demoLearnTrueFalse},
	{ID: 1341, Name: "OrderSteps", Category: "learn", Summary: "Drag steps into order; cue Drop.", Demo: demoLearnOrderSteps},
	{ID: 1342, Name: "HintReveal", Category: "learn", Summary: "Hint hidden until asked; cue Open.", Demo: demoLearnHintReveal},
	{ID: 1343, Name: "Certificate", Category: "learn", Summary: "Certificate card.", Demo: demoLearnCertificate},
	{ID: 1344, Name: "Grade", Category: "learn", Summary: "Letter grade badge.", Demo: demoLearnGrade},
	{ID: 1345, Name: "StudyTimer", Category: "learn", Summary: "Focus timer with sessions; cue Notify.", Demo: demoLearnStudyTimer},
	{ID: 1346, Name: "Notes", Category: "learn", Summary: "Highlights with notes.", Demo: demoLearnNotes},
	{ID: 1347, Name: "Syllabus", Category: "learn", Summary: "Weeks and topics.", Demo: demoLearnSyllabus},
	{ID: 1348, Name: "Instructor", Category: "learn", Summary: "Teacher card.", Demo: demoLearnInstructor},
	{ID: 1349, Name: "Recall", Category: "learn", Summary: "Again, hard, good, easy; cue Select.", Demo: demoLearnRecall},
	{ID: 1350, Name: "Progress", Category: "learn", Summary: "Course progress by module.", Demo: demoLearnProgress},
}

func demoLearnQuestions() []learn.Question {
	return []learn.Question{
		{Text: "Which clay body matures at the lowest temperature?", Choices: []string{"Earthenware", "Stoneware", "Porcelain"}, Answer: 0,
			Explain: "Earthenware matures around cone 06, far below stoneware and porcelain."},
		{Text: "What does wedging remove from clay?", Choices: []string{"Water", "Air pockets", "Grog"}, Answer: 1},
		{Text: "Which firing comes first?", Choices: []string{"Glaze firing", "Bisque firing"}, Answer: 1},
		{Text: "What tool cuts a pot from the wheel head?", Choices: []string{"Rib", "Sponge", "Wire"}, Answer: 2},
		{Text: "Which cone is hottest?", Choices: []string{"Cone 06", "Cone 6", "Cone 10"}, Answer: 2},
	}
}

func demoLearnFlashcard(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	flipped := ui.Local(root, "flipped", func() bool { return false })
	root.Children(func() {
		if learn.Flashcard(c, "What temperature does cone 6 reach?", "About 1,222 °C (2,232 °F)", flipped, learn.FlashcardOptions{}).Changed() {
			log(fmt.Sprintf("answer shown %v", *flipped))
		}
	})
}

func demoLearnQuiz(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	picked := ui.Local(root, "picked", func() int { return -1 })
	checked := ui.Local(root, "checked", func() bool { return false })
	root.Children(func() {
		if learn.Quiz(c, demoLearnQuestions()[0], picked, checked, learn.QuizOptions{}).Changed() {
			log(fmt.Sprintf("choice %d, checked %v", *picked, *checked))
		}
	})
}

func demoLearnQuizResult(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Children(func() {
		learn.QuizResult(c, demoLearnQuestions(), []int{0, 1, 0, 2, 1}, learn.QuizResultOptions{Title: "Clay basics quiz"})
	})
}

func demoLearnLessonPath(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	picked := ui.Local(root, "picked", func() int { return 3 })
	lessons := []learn.Lesson{
		{Title: "Wedging clay", Icon: "hand", Minutes: 12, Done: true},
		{Title: "Centering", Icon: "circle-dot", Minutes: 18, Done: true},
		{Title: "Opening the form", Icon: "circle-dashed", Minutes: 15, Done: true},
		{Title: "Pulling walls", Icon: "move-vertical", Minutes: 22},
		{Title: "Trimming the foot", Icon: "scissors", Minutes: 16},
		{Title: "Glazing", Icon: "paintbrush", Minutes: 25},
	}
	root.Children(func() {
		if learn.LessonPath(c, lessons, picked, learn.LessonPathOptions{}).Changed() {
			log(fmt.Sprintf("lesson %d", *picked+1))
		}
	})
}

func demoLearnCourseCard(c *ui.Context, log func(string)) {
	ui.Row(c).Width(520).Gap(16).Children(func() {
		ui.Column(c).Grow(1).Basis(0).Children(func() {
			learn.CourseCard(c, learn.Course{Title: "Wheel throwing basics", Author: "Priya Nair", Level: "Beginner", Lessons: 12, Done: 5,
				Minutes: 200, Cover: productShots()[1]}, learn.CourseCardOptions{})
		})
		ui.Column(c).Grow(1).Basis(0).Children(func() {
			learn.CourseCard(c, learn.Course{Title: "Glaze chemistry", Author: "Daniel Okafor", Level: "Advanced", Lessons: 8, Done: 8,
				Minutes: 45, Cover: productShots()[5]}, learn.CourseCardOptions{})
		})
	})
}

func demoLearnStudyStreak(c *ui.Context, log func(string)) {
	days := []learn.StudyDay{{Label: "Mon", Minutes: 25}, {Label: "Tue", Minutes: 40}, {Label: "Wed"}, {Label: "Thu", Minutes: 30},
		{Label: "Fri", Minutes: 22}, {Label: "Sat", Minutes: 12, Today: true}, {Label: "Sun"}}
	ui.Column(c).Width(420).Children(func() { learn.StudyStreak(c, days, 3, learn.StudyStreakOptions{}) })
}

func demoLearnVocabulary(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	open := ui.Local(root, "open", func() int { return 0 })
	played := ui.Local(root, "played", func() int { return 0 })
	words := []learn.Word{
		{Term: "Kintsugi", Reading: "kin-tsu-gi", Meaning: "Mending breaks with gold-dusted lacquer.", Example: "The cracked tea bowl came back from kintsugi with a seam of gold.", Tone: 392},
		{Term: "Raku", Reading: "ra-ku", Meaning: "Low-fired ware pulled glowing from the kiln.", Example: "We dropped the raku pots into sawdust to smoke the crackle.", Tone: 440},
		{Term: "Tenmoku", Reading: "ten-mo-ku", Meaning: "An iron glaze that breaks from rust to black.", Example: "Tenmoku pools dark in the throwing rings.", Tone: 494},
		{Term: "Wabi-sabi", Reading: "wa-bi sa-bi", Meaning: "Beauty in the worn and the imperfect.", Tone: 330},
	}
	root.Children(func() {
		if learn.Vocabulary(c, words, open, learn.VocabularyOptions{Title: "Pottery terms", Played: played}).Changed() {
			log(fmt.Sprintf("open %d, tones %d", *open, *played))
		}
	})
}

func demoLearnFillBlank(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	text := ui.Local(root, "text", func() string { return "" })
	checked := ui.Local(root, "checked", func() bool { return false })
	root.Children(func() {
		if learn.FillBlank(c, "Bisque firing hardens the pot so", "can go on without the clay slumping.", "glaze", text, checked,
			learn.FillBlankOptions{}).Changed() {
			log(fmt.Sprintf("checked %v, %d runes", *checked, len([]rune(*text))))
		}
	})
}

func demoLearnMatchPairs(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(460)
	pairs := []learn.Pair{{Left: "Rib", Right: "Smooths the walls"}, {Left: "Wire", Right: "Cuts the pot free"},
		{Left: "Needle tool", Right: "Levels the rim"}, {Left: "Sponge", Right: "Lifts extra water"}}
	matched := ui.Local(root, "matched", func() []bool { return make([]bool, len(pairs)) })
	picks := ui.Local(root, "picks", func() int { return 0 })
	root.Children(func() {
		if learn.MatchPairs(c, pairs, []int{2, 0, 3, 1}, matched, learn.MatchPairsOptions{Title: "Match each tool", Picks: picks}).Changed() {
			n := 0
			for _, m := range *matched {
				if m {
					n++
				}
			}
			log(fmt.Sprintf("picks %d, matched %d", *picks, n))
		}
	})
}

func demoLearnTrueFalse(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	said := ui.Local(root, "said", func() learn.Said { return learn.Unsaid })
	root.Children(func() {
		if learn.TrueFalse(c, "Porcelain is fired hotter than earthenware.", true, said,
			learn.TrueFalseOptions{Explain: "Porcelain matures near cone 10; earthenware near cone 06."}).Changed() {
			log(fmt.Sprintf("said %d", *said))
		}
	})
}

func demoLearnOrderSteps(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	correct := []string{"Wedge the clay", "Center it on the wheel", "Open the form", "Pull up the walls", "Trim the foot"}
	steps := ui.Local(root, "steps", func() []string {
		return []string{"Open the form", "Wedge the clay", "Pull up the walls", "Center it on the wheel", "Trim the foot"}
	})
	root.Children(func() {
		if learn.OrderSteps(c, steps, correct, learn.OrderStepsOptions{}).Changed() {
			n := 0
			for i, s := range *steps {
				if s == correct[i] {
					n++
				}
			}
			log(fmt.Sprintf("moved, %d in place", n))
		}
	})
}

func demoLearnHintReveal(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	shown := ui.Local(root, "shown", func() bool { return false })
	root.Children(func() {
		if learn.HintReveal(c, "Brace your elbows on your thighs and let the wheel's speed, not your arms, do the centering.", shown,
			learn.HintRevealOptions{}).Changed() {
			log(fmt.Sprintf("hint shown %v", *shown))
		}
	})
}

func demoLearnCertificate(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Children(func() {
		learn.Certificate(c, learn.Award{Name: "Maya Chen", Course: "Wheel Throwing Basics", Issuer: "Kiln & Co. Clay School",
			Signer: "Priya Nair", ID: "KC-2026-0418", Date: time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)}, learn.CertificateOptions{})
	})
}

func demoLearnGrade(c *ui.Context, log func(string)) {
	ui.Column(c).Width(340).Gap(16).Children(func() {
		learn.Grade(c, 94.5, learn.GradeOptions{Label: "Glaze chemistry quiz"})
		learn.Grade(c, 81, learn.GradeOptions{Label: "Throwing practical"})
		learn.Grade(c, 67, learn.GradeOptions{Label: "Kiln log"})
	})
}

func demoLearnStudyTimer(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(360)
	sessions := ui.Local(root, "sessions", func() int { return 1 })
	root.Children(func() {
		timer, run, skip := learn.StudyTimer(c, sessions, learn.StudyTimerOptions{})
		if run.Clicked() {
			log("timer toggled")
		}
		if skip.Clicked() {
			log("phase skipped")
		}
		if timer.Changed() {
			log(fmt.Sprintf("sessions %d", *sessions))
		}
	})
}

func demoLearnNotes(c *ui.Context, log func(string)) {
	items := []learn.Highlight{
		{Quote: "A pot is made as much by the space inside as by the clay around it.", Note: "Use this for the opening-the-form lesson.",
			Where: "The Craft of Clay · p. 42", Tint: 0},
		{Quote: "Centering is not force; it is patience applied evenly.", Note: "Priya says the same thing every class.",
			Where: "The Craft of Clay · p. 57", Tint: 2},
		{Quote: "Trim when the clay is leather-hard: firm enough to hold, soft enough to carve.", Where: "Studio handbook · p. 9", Tint: 1},
	}
	ui.Column(c).Width(440).Children(func() { learn.Notes(c, items, learn.NotesOptions{}) })
}

func demoLearnSyllabus(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	open := ui.Local(root, "open", func() int { return 2 })
	weeks := []learn.Week{
		{Title: "Clay and wedging", Topics: []string{"Clay bodies", "Ram's head wedging", "Studio safety"}, Done: true},
		{Title: "Centering", Topics: []string{"Wheel speed", "Coning up and down"}, Done: true},
		{Title: "Cylinders", Topics: []string{"Opening", "Pulling walls", "Compressing the rim"}},
		{Title: "Bowls", Topics: []string{"Wide openings", "Shaping curves"}},
		{Title: "Trimming and glazing", Topics: []string{"Leather-hard timing", "Trimming a foot", "Dipping glaze"}},
	}
	root.Children(func() {
		if learn.Syllabus(c, weeks, open, learn.SyllabusOptions{}).Changed() {
			log(fmt.Sprintf("open week %d", *open+1))
		}
	})
}

func demoLearnInstructor(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(380)
	followed := ui.Local(root, "followed", func() bool { return false })
	root.Children(func() {
		if learn.Instructor(c, learn.Teacher{Name: "Priya Nair", Role: "Ceramicist · Kiln & Co. Clay School", Rating: 4.9, Students: 12480,
			Courses: 6, Bio: "Priya has thrown on the wheel for 18 years and teaches beginners to trust slow, even pressure."},
			followed, learn.InstructorOptions{}).Changed() {
			log(fmt.Sprintf("following %v", *followed))
		}
	})
}

func demoLearnRecall(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(420)
	grade := ui.Local(root, "grade", func() int { return -1 })
	root.Children(func() {
		if learn.Recall(c, [4]string{"<1m", "6m", "1d", "4d"}, grade, learn.RecallOptions{}).Changed() {
			log(fmt.Sprintf("grade %d", *grade))
		}
	})
}

func demoLearnProgress(c *ui.Context, log func(string)) {
	modules := []learn.Module{{Title: "Clay basics", Done: 4, Total: 4}, {Title: "Centering", Done: 3, Total: 3},
		{Title: "Cylinders and bowls", Done: 2, Total: 5}, {Title: "Trimming", Total: 3}, {Title: "Glazing", Total: 4}}
	ui.Column(c).Width(420).Children(func() { learn.Progress(c, modules, learn.ProgressOptions{Title: "Wheel throwing basics"}) })
}
