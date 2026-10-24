package agent

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestProposedDiff(t *testing.T) {
	diff := "@@ -1,3 +1,3 @@\n # Launch\n-Ship in June\n+Ship in July\n+Add pricing slide\n"
	for _, tc := range []struct {
		button, shown string
		want          Decision
		cue           sound.Cue
	}{{"Apply", "Applied", Allowed, sound.Success}, {"Reject", "Rejected", Denied, sound.Close}} {
		decision, changed := Undecided, 0
		tt := ui.NewTester(padded(func(c *ui.Context) {
			if root, _, _ := ProposedDiff(c, "docs/deck.md", diff, &decision, ProposedDiffOptions{}); root.Changed() {
				changed++
			}
		}), 600, 400)
		settles(t, tt, "diff")
		if !tt.HasText("+2") || !tt.HasText("−1") || !tt.HasText("Ship in July") {
			t.Fatalf("texts %q", tt.Texts())
		}
		before := baseline()
		click(t, tt, tc.button)
		settles(t, tt, "decided diff")
		if decision != tc.want || changed != 1 || !tt.HasText(tc.shown) || tt.Click("Apply") == nil {
			t.Errorf("%s: decision %d changed %d texts %q", tc.button, decision, changed, tt.Texts())
		}
		cues(t, before, tc.cue)
	}
	panics(t, "nil decision", func() { render(func(c *ui.Context) { ProposedDiff(c, "a", "+b", nil, ProposedDiffOptions{}) }) })
}

func TestEvalScore(t *testing.T) {
	checks := []EvalCheck{{Name: "Accuracy", Score: 0.92, Pass: 0.8}, {Name: "Tone", Score: 0.55, Pass: 0.7}}
	tt := ui.NewTester(padded(func(c *ui.Context) { EvalScore(c, "Evals", checks, EvalScoreOptions{}) }), 500, 300)
	settles(t, tt, "evals")
	for _, s := range []string{"1 of 2 passed", "0.92", "0.55"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	panics(t, "over one", func() {
		render(func(c *ui.Context) { EvalScore(c, "E", []EvalCheck{{Name: "x", Score: 2}}, EvalScoreOptions{}) })
	})
}

func TestRateAnswer(t *testing.T) {
	rating, reason, sent, sends := 0, -1, false, 0
	reasons := []string{"Inaccurate", "Too long"}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, send := RateAnswer(c, &rating, reasons, &reason, &sent, RateAnswerOptions{}); send.Clicked() {
			sends++
		}
	}), 500, 300)
	settles(t, tt, "rate")
	if tt.HasText("Too long") {
		t.Fatal("reasons shown before a rating")
	}
	before := baseline()
	click(t, tt, "Not helpful")
	settles(t, tt, "reasons")
	click(t, tt, "Send feedback")
	if sent {
		t.Fatal("sent without a reason")
	}
	click(t, tt, "Too long")
	click(t, tt, "Send feedback")
	settles(t, tt, "thanks")
	if !sent || rating != -1 || reason != 1 || !tt.HasText("Thanks for the feedback") {
		t.Errorf("rating %d reason %d sent %v texts %q", rating, reason, sent, tt.Texts())
	}
	cues(t, before, sound.Toggle, sound.Select, sound.Send)
}

func TestStreamCode(t *testing.T) {
	code := "package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"
	done, copies := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, cp := StreamCode(c, "go", code, done, StreamCodeOptions{}); cp.Clicked() {
			copies++
		}
	}), 500, 400)
	tt.Frame()
	tt.Advance(130 * time.Millisecond)
	if strings.Contains(strings.Join(tt.Texts(), " "), "println") {
		t.Errorf("revealed everything: %q", tt.Texts())
	}
	tt.Advance(2 * time.Second)
	if !tt.HasText("    println(\"hi\")") || !tt.HasText("5") {
		t.Errorf("stream %q", tt.Texts())
	}
	done = true
	tt.Frame()
	settles(t, tt, "stream")
	before := baseline()
	click(t, tt, "Copy")
	if tt.Clipboard() != code || copies != 1 {
		t.Errorf("clipboard %q copies %d", tt.Clipboard(), copies)
	}
	cues(t, before, sound.Pop)
	panics(t, "empty code", func() { render(func(c *ui.Context) { StreamCode(c, "go", "", false, StreamCodeOptions{}) }) })
}

func TestSearchSteps(t *testing.T) {
	steps := []SearchStep{{Label: "Searching", Count: 12, State: StepDone}, {Label: "Reading", Count: 3, State: StepRunning}, {Label: "Writing"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { SearchSteps(c, "Research", steps) }), 600, 200)
	tt.Frame()
	if !tt.HasText("12") || !tt.HasText("Reading") || !moves(tt, 200*time.Millisecond) {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "negative count", func() {
		render(func(c *ui.Context) { SearchSteps(c, "S", []SearchStep{{Label: "x", Count: -1}}) })
	})
}

func TestImageGenerating(t *testing.T) {
	progress := 0.4
	tt := ui.NewTester(padded(func(c *ui.Context) { ImageGenerating(c, "A lighthouse at dusk", progress) }), 500, 450)
	tt.Frame()
	if !tt.HasText("40%") || !tt.HasText("A lighthouse at dusk") || !moves(tt, 200*time.Millisecond) {
		t.Errorf("texts %q", tt.Texts())
	}
	progress = 1
	tt.Frame()
	settles(t, tt, "finished image")
	panics(t, "over one", func() { render(func(c *ui.Context) { ImageGenerating(c, "x", 2) }) })
}

func TestAgentHandoff(t *testing.T) {
	at := time.Date(2026, 10, 24, 9, 30, 0, 0, time.Local)
	tt := ui.NewTester(padded(func(c *ui.Context) { AgentHandoff(c, "Researcher", "Writer", "Sources gathered", at) }), 600, 300)
	tt.Frame()
	if !tt.HasText("Researcher") || !tt.HasText("Sources gathered") || !moves(tt, 200*time.Millisecond) {
		t.Errorf("texts %q", tt.Texts())
	}
	settles(t, tt, "handoff")
	panics(t, "zero time", func() { render(func(c *ui.Context) { AgentHandoff(c, "a", "b", "", time.Time{}) }) })
}

func TestStopButton(t *testing.T) {
	running, changed := true, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if StopButton(c, running, StopButtonOptions{}).Changed() {
			changed++
		}
	}), 300, 200)
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
	settles(t, tt, "stopped")
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
	}), 500, 300)
	settles(t, tt, "history")
	if !tt.HasText("1 of 3") || !tt.HasText("Plan the sprint") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Older")
	if at != 1 || changed != 1 || !tt.HasText("2 of 3") {
		t.Fatalf("at %d changed %d texts %q", at, changed, tt.Texts())
	}
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if at != 2 {
		t.Errorf("up keys left at %d", at)
	}
	tt.Key(0, ui.KeyDown)
	tt.Frame()
	if at != 1 {
		t.Errorf("down key left at %d", at)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	panics(t, "bad position", func() {
		n := 3
		render(func(c *ui.Context) { PromptHistory(c, entries, &n, PromptHistoryOptions{}) })
	})
}

func TestAgentWelcome(t *testing.T) {
	chosen, changed := -1, 0
	prompts := []Suggestion{{Text: "Draft the deck", Icon: "presentation"}, {Text: "Plan a trip", Icon: "map"}}
	caps := []Capability{{Icon: "search", Title: "Research", Text: "Reads the web"}, {Icon: "code", Title: "Code"}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if AgentWelcome(c, "Good morning, Ana", prompts, caps, &chosen, AgentWelcomeOptions{Subtitle: "How can I help?"}).Changed() {
			changed++
		}
	}), 700, 500)
	tt.Frame()
	if !moves(tt, 100*time.Millisecond) {
		t.Error("welcome does not rise in")
	}
	settles(t, tt, "welcome")
	if !tt.HasText("Reads the web") || !tt.HasText("How can I help?") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Plan a trip")
	if chosen != 1 || changed != 1 {
		t.Errorf("chosen %d changed %d", chosen, changed)
	}
	cues(t, before, sound.Tap)
}
