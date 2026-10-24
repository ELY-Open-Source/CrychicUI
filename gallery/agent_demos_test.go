package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestAgentDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"PromptBox": {func(tt *ui.Tester) error {
			if err := tt.Click("Ask about the launch plan"); err != nil {
				return err
			}
			tt.Type("hi")
			tt.Key(0, ui.KeyEnter)
			return nil
		}, "sent 2 runes"},
		"StreamText":    {press("Restart"), "restarted"},
		"Reasoning":     {press("Thought for 14s"), "open true"},
		"ToolCall":      {press("Finish write"), "done true"},
		"ToolResult":    {press("list_files"), "open false"},
		"Plan":          {press("Next step"), "step 2"},
		"TaskProgress":  {press("Advance"), "step 3 of 5"},
		"Citation":      {press("1 Layout: rows, columns and grids"), "opened 1"},
		"SourceCard":    {press("Go 1.27 release notes"), "opened 2"},
		"ModelPicker":   {press("Local 8B"), "selected Local 8B"},
		"TokenCount":    {press("Add 512"), "3360 tokens"},
		"ContextMeter":  {press("Add 20k tokens"), "116000 tokens"},
		"CostMeter":     {press("Run a batch"), "USD 36.40"},
		"Suggestions":   {press("Draft a polite reply"), "fill: Draft a polite reply"},
		"Branches":      {press("Next"), "at 3 of 3"},
		"AnswerActions": {press("Like"), "vote 1"},
		"Artifact":      {press("Open"), "opened launch-deck.html"},
		"Approval":      {press("Allow"), "allowed"},
		"ProposedDiff":  {press("Apply"), "applied"},
		"VoiceOrb":      {press("Louder"), "level 0.8"},
		"Memories":      {press("Forget"), "2 left"},
		"AgentStatus":   {press("Finish"), "state 3"},
		"AgentBoard":    {press("Unblock Gamma"), "gamma working"},
		"ContextFiles":  {press("Remove"), "2 files left"},
		"SlashCommands": {press("/deploy"), "selected /deploy"},
		"SystemPrompt": {func(tt *ui.Tester) error {
			if err := tt.Click("System prompt"); err != nil {
				return err
			}
			tt.Type("x")
			return nil
		}, "edited 34 runes"},
		"Temperature": {func(tt *ui.Tester) error {
			tt.Key(0, ui.KeyTab)
			tt.Key(0, ui.KeyRight)
			return nil
		}, "value 0.51"},
		"RateAnswer":      {press("Not helpful"), "rating -1, reason -1"},
		"StreamCode":      {press("Restart"), "restarted"},
		"ImageGenerating": {press("10% more"), "50%"},
		"StopButton":      {press("Stop"), "stopped"},
		"PromptHistory":   {press("Older"), "at 2"},
		"AgentWelcome":    {press("Draft the deck"), "picked Draft the deck"},
		"PromptLibrary": {func(tt *ui.Tester) error {
			tt.Key(0, ui.KeyTab)
			tt.Type("release")
			tt.Frame()
			return tt.Click("Release note")
		}, "selected Release note"},
	}
	// Quiet demos only show; they have no events.
	quiet := map[string]bool{"Thinking": true, "AgentAvatar": true, "RunTimeline": true, "EvalScore": true, "SearchSteps": true, "AgentHandoff": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "agent" {
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
		tt.Advance(time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		found := false
		for _, l := range logs {
			found = found || strings.HasPrefix(l, act.want)
		}
		if !found {
			t.Errorf("%s logged %q, want a line starting %q", e.Name, logs, act.want)
		}
	}
	if seen != len(acts)+len(quiet) {
		t.Errorf("checked %d agent demos, want %d", seen, len(acts)+len(quiet))
	}
}

func TestSuggestionsDemoLogsTheBox(t *testing.T) {
	var logs []string
	tt := ui.NewTester(func(c *ui.Context) { demoSuggestions(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
	tt.Advance(time.Second)
	if err := tt.Click("Ask anything"); err != nil {
		t.Fatal(err)
	}
	tt.Type("hi")
	tt.Frame()
	if err := tt.Click("Attach"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	joined := strings.Join(logs, "\n")
	if !strings.Contains(joined, "draft 2 runes") || !strings.Contains(joined, "attach") {
		t.Errorf("logs %q", logs)
	}
}

func TestStreamTextDemoRestarts(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { demoStreamText(c, func(string) {}) }, 900, 700)
	tt.Advance(4 * time.Second)
	if !tt.HasText("milestones.") {
		t.Fatalf("stream did not finish: %q", tt.Texts())
	}
	if err := tt.Click("Restart"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if tt.HasText("milestones.") {
		t.Error("full text stays right after a restart")
	}
	tt.Advance(4 * time.Second)
	if !tt.HasText("milestones.") {
		t.Error("stream did not finish after the restart")
	}
}
