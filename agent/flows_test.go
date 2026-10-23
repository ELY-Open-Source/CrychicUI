package agent

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestProposedDiff(t *testing.T) {
	diff := "@@ -1,2 +1,2 @@\n-old line\n+new line"
	applied, rejected, changed := 0, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, apply, reject := ProposedDiff(c, "deck.md", diff, ProposedDiffOptions{})
		if root.Changed() {
			changed++
		}
		if apply.Clicked() {
			applied++
		}
		if reject.Clicked() {
			rejected++
		}
	}), 500, 300)
	settles(t, tt, "diff")
	if !tt.HasText("-old line") || !tt.HasText("+new line") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Apply")
	if applied != 1 || changed != 1 {
		t.Fatalf("applied %d changed %d", applied, changed)
	}
	click(t, tt, "Reject")
	if rejected != 1 || changed != 2 {
		t.Errorf("rejected %d changed %d", rejected, changed)
	}
	cues(t, before, sound.Success, sound.Close)
	panics(t, "empty diff", func() { render(func(c *ui.Context) { ProposedDiff(c, "x", "", ProposedDiffOptions{}) }) })
}

func TestEvalScore(t *testing.T) {
	checks := []EvalCheck{{Name: "Accuracy", Score: 0.82, Pass: 0.7}, {Name: "Tone", Score: 0.55, Pass: 0.7}}
	tt := ui.NewTester(padded(func(c *ui.Context) { EvalScore(c, "Evals", checks) }), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bars do not grow")
	}
	settles(t, tt, "scores")
	if !tt.HasText("0.82") || !tt.HasText("0.55") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "score over one", func() {
		render(func(c *ui.Context) { EvalScore(c, "x", []EvalCheck{{Name: "a", Score: 1.5}}) })
	})
}

func TestRateAnswer(t *testing.T) {
	rating, reason, changed, sends := 0, -1, 0, 0
	reasons := []string{"Wrong facts", "Too long"}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, send := RateAnswer(c, &rating, reasons, &reason, RateAnswerOptions{})
		if root.Changed() {
			changed++
		}
		if send.Clicked() {
			sends++
		}
	}), 500, 300)
	settles(t, tt, "rating")
	before := baseline()
	click(t, tt, "Not helpful")
	if rating != -1 || changed != 1 {
		t.Fatalf("rating %d changed %d", rating, changed)
	}
	click(t, tt, "Send feedback")
	if sends != 0 {
		t.Error("sent without a reason")
	}
	click(t, tt, "Too long")
	if reason != 1 || changed != 2 {
		t.Fatalf("reason %d changed %d", reason, changed)
	}
	click(t, tt, "Send feedback")
	if sends != 1 || changed != 3 {
		t.Errorf("sends %d changed %d", sends, changed)
	}
	cues(t, before, sound.Toggle, sound.Select, sound.Send)
	panics(t, "bad rating", func() {
		v := 2
		render(func(c *ui.Context) { RateAnswer(c, &v, reasons, &reason, RateAnswerOptions{}) })
	})
}

func TestStreamCode(t *testing.T) {
	code := "package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"
	done := false
	tt := ui.NewTester(padded(func(c *ui.Context) { StreamCode(c, "go", code, done) }), 500, 300)
	tt.Frame()
	tt.Advance(130 * time.Millisecond)
	joined := strings.Join(tt.Texts(), " ")
	if strings.Contains(joined, "println") {
		t.Errorf("revealed everything: %q", joined)
	}
	if !strings.HasPrefix(code, "package main\n\nfunc main() {\n") {
		t.Error("bad test source")
	}
	tt.Advance(2 * time.Second)
	if !strings.Contains(strings.Join(tt.Texts(), " "), "println") {
		t.Errorf("stream %q", tt.Texts())
	}
	done = true
	tt.Frame()
	settles(t, tt, "stream")
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	done = false
	tt.Frame()
	if !strings.Contains(strings.Join(tt.Texts(), " "), "println") {
		t.Errorf("reduced motion %q", tt.Texts())
	}
	panics(t, "empty code", func() { render(func(c *ui.Context) { StreamCode(c, "go", "", false) }) })
}

func TestSearchSteps(t *testing.T) {
	steps := []SearchStep{{Label: "Searching", Done: true}, {Label: "Reading", Count: 3}, {Label: "Writing"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { SearchSteps(c, "Search", steps) }), 400, 300)
	tt.Frame()
	if !moves(tt, 150*time.Millisecond) {
		t.Error("active step does not spin")
	}
	if !tt.HasText("Reading 3") {
		t.Errorf("texts %q", tt.Texts())
	}
	steps[1].Done = true
	steps[2] = SearchStep{Label: "Writing"}
	tt.Frame()
	if !tt.HasText("Reading 3") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "empty label", func() {
		render(func(c *ui.Context) { SearchSteps(c, "x", []SearchStep{{Label: ""}}) })
	})
}

func TestImageGenerating(t *testing.T) {
	progress := 0.4
	tt := ui.NewTester(padded(func(c *ui.Context) { ImageGenerating(c, "Image", progress) }), 500, 300)
	tt.Frame()
	if !tt.HasText("40%") {
		t.Errorf("texts %q", tt.Texts())
	}
	if !moves(tt, 200*time.Millisecond) {
		t.Error("sheen does not move")
	}
	progress = 1
	tt.Frame()
	settles(t, tt, "placeholder")
	if !tt.HasText("100%") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "over one", func() { render(func(c *ui.Context) { ImageGenerating(c, "x", 1.2) }) })
}

func TestAgentHandoff(t *testing.T) {
	at := time.Date(2026, 10, 22, 9, 30, 0, 0, time.UTC)
	tt := ui.NewTester(padded(func(c *ui.Context) { AgentHandoff(c, "Planner", "Writer", "Took the deck", at) }), 500, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("event does not fade in")
	}
	settles(t, tt, "handoff")
	if !tt.HasText("Took the deck") || !tt.HasText("9:30") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "zero time", func() { render(func(c *ui.Context) { AgentHandoff(c, "a", "b", "", time.Time{}) }) })
}

func TestStopButton(t *testing.T) {
	running, changed := true, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if StopButton(c, running, StopButtonOptions{}).Changed() {
			changed++
		}
	}), 400, 200)
	tt.Frame()
	if !moves(tt, 200*time.Millisecond) {
		t.Error("stop does not pulse")
	}
	before := baseline()
	click(t, tt, "Stop")
	if changed != 1 {
		t.Fatalf("changed %d", changed)
	}
	cues(t, before, sound.Close)
	running = false
	tt.Frame()
	settles(t, tt, "stop")
	again := baseline()
	click(t, tt, "Stop")
	cues(t, again)
	if changed != 1 {
		t.Error("disabled stop reported a change")
	}
}

func TestPromptHistory(t *testing.T) {
	entries := []string{"Draft the deck", "Review the brief", "Plan the sprint"}
	at, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if PromptHistory(c, entries, &at, PromptHistoryOptions{}).Changed() {
			changed++
		}
	}), 400, 200)
	settles(t, tt, "history")
	if !tt.HasText("1 of 3") || !tt.HasText("Draft the deck") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Older")
	if at != 1 || changed != 1 || !tt.HasText("Review the brief") {
		t.Fatalf("at %d changed %d texts %q", at, changed, tt.Texts())
	}
	click(t, tt, "Older")
	click(t, tt, "Older")
	if at != 2 {
		t.Error("moved past the oldest entry")
	}
	click(t, tt, "Newer")
	if at != 1 {
		t.Errorf("at %d", at)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	panics(t, "bad position", func() {
		n := 3
		render(func(c *ui.Context) { PromptHistory(c, entries, &n, PromptHistoryOptions{}) })
	})
	panics(t, "empty entry", func() {
		n := 0
		render(func(c *ui.Context) { PromptHistory(c, []string{""}, &n, PromptHistoryOptions{}) })
	})
}

func TestAgentWelcome(t *testing.T) {
	chosen, changed := -1, 0
	prompts := []string{"Plan the launch", "Draft the deck"}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if AgentWelcome(c, "What are we shipping?", prompts, []string{"Search", "Code", "Files"}, &chosen).Changed() {
			changed++
		}
	}), 500, 400)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("welcome does not rise in")
	}
	settles(t, tt, "welcome")
	if !tt.HasText("What are we shipping?") || !tt.HasText("Search") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Draft the deck")
	if chosen != 1 || changed != 1 {
		t.Fatalf("chosen %d changed %d", chosen, changed)
	}
	cues(t, before, sound.Tap)
	panics(t, "empty prompt", func() {
		render(func(c *ui.Context) { AgentWelcome(c, "g", []string{""}, nil, &chosen) })
	})
}
