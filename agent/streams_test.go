package agent

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPromptBox(t *testing.T) {
	draft, changed := "", 0
	var sent []string
	var attached int
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, send, attach := PromptBox(c, &draft, "GPT-6.1", func(s string) { sent = append(sent, s) }, PromptBoxOptions{})
		if root.Changed() {
			changed++
		}
		if send.Clicked() && draft != "" {
			t.Error("send clicked with a draft left")
		}
		if attach.Clicked() {
			attached++
		}
	}), 500, 300)
	settles(t, tt, "prompt box")
	if !tt.HasText("GPT-6.1") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Ask anything")
	tt.Type("review the deck")
	tt.Frame()
	if changed == 0 {
		t.Fatal("typing did not report a change")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(sent) != 1 || sent[0] != "review the deck" || draft != "" {
		t.Fatalf("sent %q draft %q", sent, draft)
	}
	cues(t, before, sound.Send)
	again := baseline()
	click(t, tt, "Attach")
	if attached != 1 {
		t.Errorf("attached %d", attached)
	}
	cues(t, again, sound.Tap)
	panics(t, "empty model", func() {
		render(func(c *ui.Context) { PromptBox(c, &draft, "", func(string) {}, PromptBoxOptions{}) })
	})
	panics(t, "nil send", func() {
		render(func(c *ui.Context) { PromptBox(c, &draft, "m", nil, PromptBoxOptions{}) })
	})
}

func TestStreamText(t *testing.T) {
	text := "one two three four five six seven eight nine ten"
	done := false
	tt := ui.NewTester(padded(func(c *ui.Context) { StreamText(c, text, done) }), 500, 200)
	tt.Frame()
	joined := func() string { return strings.Join(tt.Texts(), " ") }
	first := joined()
	tt.Advance(700 * time.Millisecond)
	second := joined()
	if len(second) <= len(first) {
		t.Errorf("text did not grow: %q then %q", first, second)
	}
	if strings.Contains(second, "nine ten") {
		t.Errorf("revealed everything: %q", second)
	}
	spaced := "one\n  two\tthree four"
	st2 := ui.NewTester(padded(func(c *ui.Context) { StreamText(c, spaced, false) }), 500, 200)
	st2.Frame()
	st2.Advance(170 * time.Millisecond)
	for _, x := range st2.Texts() {
		if x == "StreamText" {
			continue
		}
		if !strings.HasPrefix(spaced, x) {
			t.Errorf("%q is not a prefix of %q", x, spaced)
		}
	}
	tt.Advance(2 * time.Second)
	if !strings.Contains(joined(), "nine ten") || !tt.FramePending() {
		t.Errorf("stream %q pending %v", joined(), tt.FramePending())
	}
	done = true
	tt.Frame()
	settles(t, tt, "stream")
	if !strings.Contains(joined(), "nine ten") {
		t.Errorf("final %q", joined())
	}
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	done = false
	tt.Frame()
	if !strings.Contains(joined(), "nine ten") {
		t.Errorf("reduced motion %q", joined())
	}
}

func TestThinking(t *testing.T) {
	var base time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if base.IsZero() {
			base = c.Now()
		}
		Thinking(c, base, ThinkingOptions{})
	}), 400, 200)
	tt.Frame()
	if !tt.HasText("Thinking") || !tt.HasText("0s") {
		t.Errorf("texts %q", tt.Texts())
	}
	if !moves(tt, 300*time.Millisecond) {
		t.Error("label does not shimmer")
	}
	tt.Advance(2100 * time.Millisecond)
	if !tt.HasText("2s") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "zero time", func() { render(func(c *ui.Context) { Thinking(c, time.Time{}, ThinkingOptions{}) }) })
}

func TestReasoning(t *testing.T) {
	open, changed := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Reasoning(c, "Noon suits all three calendars.", &open, ReasoningOptions{Took: 12 * time.Second}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "reasoning")
	if !tt.HasText("Thought for 12s") || tt.HasText("Noon suits all three calendars.") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Thought for 12s")
	settles(t, tt, "open reasoning")
	if !open || changed != 1 || !tt.HasText("Noon suits all three calendars.") {
		t.Errorf("open %v changed %d texts %q", open, changed, tt.Texts())
	}
	cues(t, before, sound.Toggle)
	plain := ui.NewTester(padded(func(c *ui.Context) { Reasoning(c, "x", &open, ReasoningOptions{}) }), 300, 200)
	plain.Frame()
	if !plain.HasText("Reasoning") {
		t.Errorf("plain %q", plain.Texts())
	}
	panics(t, "negative time", func() {
		render(func(c *ui.Context) { Reasoning(c, "x", &open, ReasoningOptions{Took: -time.Second}) })
	})
}

func TestToolCall(t *testing.T) {
	open, state, changed := false, StepRunning, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ToolCall(c, "web_search", `{"query":"launch deck"}`, state, &open, ToolCallOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	tt.Frame()
	if !tt.HasText("Running") || !tt.FramePending() {
		t.Fatalf("running %q pending %v", tt.Texts(), tt.FramePending())
	}
	before := baseline()
	click(t, tt, "web_search")
	tt.Advance(time.Second)
	if !open || changed != 1 || !tt.HasText("Input") {
		t.Fatalf("open %v changed %d texts %q", open, changed, tt.Texts())
	}
	cues(t, before, sound.Toggle)
	state = StepDone
	tt.Frame()
	settles(t, tt, "done call")
	if !tt.HasText("Done") || tt.HasText("Running") {
		t.Errorf("done %q", tt.Texts())
	}
	state = StepFailed
	tt.Frame()
	if !tt.HasText("Failed") {
		t.Errorf("failed %q", tt.Texts())
	}
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), `ToolCall "web_search"`) {
			t.Errorf("bad state panic %v", r)
		}
	}()
	render(func(c *ui.Context) { ToolCall(c, "web_search", "", StepFailed+1, &open, ToolCallOptions{}) })
}

func TestToolResult(t *testing.T) {
	open, changed, copies := false, 0, 0
	body := "deck-v3.html\ndeck-notes.md\nlaunch-plan.pdf\n"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, cp := ToolResult(c, "list_files", body, &open, ToolResultOptions{})
		if root.Changed() {
			changed++
		}
		if cp.Clicked() {
			copies++
		}
	}), 500, 400)
	settles(t, tt, "result")
	if !tt.HasText("3 lines") || tt.HasText("deck-v3.html\ndeck-notes.md\nlaunch-plan.pdf\n") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "list_files")
	settles(t, tt, "open result")
	if !open || changed != 1 || !tt.HasText(body) {
		t.Errorf("open %v changed %d texts %q", open, changed, tt.Texts())
	}
	click(t, tt, "Copy")
	if tt.Clipboard() != body || copies != 1 {
		t.Errorf("clipboard %q copies %d", tt.Clipboard(), copies)
	}
	cues(t, before, sound.Toggle, sound.Pop)
}

func TestPlan(t *testing.T) {
	steps := []PlanStep{
		{Title: "Read the brief", State: StepDone},
		{Title: "Draft the deck", Detail: "12 slides", State: StepRunning},
		{Title: "Review together", State: StepPending},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { Plan(c, "Launch plan", steps, PlanOptions{}) }), 500, 400)
	tt.Frame()
	if !tt.HasText("1 of 3 done") || !tt.HasText("12 slides") || !tt.FramePending() {
		t.Fatalf("texts %q pending %v", tt.Texts(), tt.FramePending())
	}
	steps[1].State, steps[2].State = StepDone, StepDone
	tt.Frame()
	settles(t, tt, "finished plan")
	if !tt.HasText("3 of 3 done") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no title", func() {
		render(func(c *ui.Context) { Plan(c, "P", []PlanStep{{}}, PlanOptions{}) })
	})
}

func TestTaskProgress(t *testing.T) {
	steps := []string{"Gather notes", "Outline", "Draft slides", "Polish"}
	current := 2
	tt := ui.NewTester(padded(func(c *ui.Context) { TaskProgress(c, "Launch deck", steps, current, TaskProgressOptions{}) }), 500, 250)
	tt.Frame()
	for _, s := range []string{"50%", "Draft slides", "Step 3 of 4"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	current = 4
	tt.Frame()
	settles(t, tt, "finished task")
	if !tt.HasText("100%") || !tt.HasText("Finished") || !tt.HasText("Step 4 of 4") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "past the end", func() {
		render(func(c *ui.Context) { TaskProgress(c, "T", steps, 5, TaskProgressOptions{}) })
	})
}
