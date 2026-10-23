package agent

import (
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
		if Reasoning(c, "The user wants lunch; noon works.", &open, ReasoningOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "panel")
	if tt.HasText("noon works") {
		t.Error("body shown while closed")
	}
	before := baseline()
	click(t, tt, "Reasoning")
	if !open || changed != 1 {
		t.Fatalf("open %v changed %d", open, changed)
	}
	if !tt.FramePending() {
		t.Error("height does not spring")
	}
	settles(t, tt, "panel")
	if !tt.HasText("noon works") {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Toggle)
	click(t, tt, "Reasoning")
	settles(t, tt, "panel")
	if tt.HasText("noon works") {
		t.Error("body stays after closing")
	}
	panics(t, "empty body", func() { render(func(c *ui.Context) { Reasoning(c, "", &open, ReasoningOptions{}) }) })
}

func TestToolCall(t *testing.T) {
	running := true
	tt := ui.NewTester(padded(func(c *ui.Context) { ToolCall(c, "search", `query "deck"`, running) }), 500, 200)
	tt.Frame()
	if !moves(tt, 150*time.Millisecond) {
		t.Error("status icon does not spin")
	}
	running = false
	tt.Frame()
	if !tt.FramePending() {
		t.Error("check does not pop in")
	}
	settles(t, tt, "tool call")
	if !tt.HasText("search") || !tt.HasText(`query "deck"`) {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "empty name", func() { render(func(c *ui.Context) { ToolCall(c, "", "", false) }) })
}

func TestToolResult(t *testing.T) {
	open, changed := true, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ToolResult(c, "3 files matched", &open, ToolResultOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "result")
	if !tt.HasText("3 files matched") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Result")
	if open || changed != 1 {
		t.Fatalf("open %v changed %d", open, changed)
	}
	settles(t, tt, "result")
	if tt.HasText("3 files matched") {
		t.Error("body stays after closing")
	}
	cues(t, before, sound.Toggle)
	panics(t, "empty body", func() { render(func(c *ui.Context) { ToolResult(c, "", &open, ToolResultOptions{}) }) })
}

func TestPlan(t *testing.T) {
	current := 1
	steps := []string{"Read the brief", "Draft the deck", "Review together"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Plan(c, "Plan", steps, current) }), 500, 300)
	tt.Frame()
	if !moves(tt, 150*time.Millisecond) {
		t.Error("current mark does not spin")
	}
	current = 2
	tt.Frame()
	if !tt.FramePending() {
		t.Error("done mark does not pop")
	}
	tt.Advance(2 * time.Second)
	if !tt.FramePending() {
		t.Error("new current does not spin")
	}
	panics(t, "bad current", func() { render(func(c *ui.Context) { Plan(c, "x", steps, 3) }) })
	panics(t, "empty step", func() { render(func(c *ui.Context) { Plan(c, "x", []string{"a", ""}, 0) }) })
}

func TestTaskProgress(t *testing.T) {
	step := 2
	tt := ui.NewTester(padded(func(c *ui.Context) { TaskProgress(c, "Deck", step, 5, TaskProgressOptions{}) }), 500, 200)
	settles(t, tt, "progress")
	if !tt.HasText("Step 2 of 5") {
		t.Errorf("texts %q", tt.Texts())
	}
	step = 3
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bar does not fill")
	}
	settles(t, tt, "progress")
	if !tt.HasText("Step 3 of 5") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "step past total", func() { render(func(c *ui.Context) { TaskProgress(c, "x", 6, 5, TaskProgressOptions{}) }) })
	panics(t, "zero total", func() { render(func(c *ui.Context) { TaskProgress(c, "x", 0, 0, TaskProgressOptions{}) }) })
}
