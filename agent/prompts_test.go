package agent

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestSuggestions(t *testing.T) {
	prompts := []string{"Summarize the thread", "Draft a reply", "List next steps"}
	chosen, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Suggestions(c, "Suggestions", prompts, &chosen).Changed() {
			changed++
		}
	}), 500, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("chips do not arrive")
	}
	settles(t, tt, "chips")
	before := baseline()
	click(t, tt, "Draft a reply")
	if chosen != 1 || changed != 1 {
		t.Fatalf("chosen %d changed %d", chosen, changed)
	}
	cues(t, before, sound.Tap)
	panics(t, "empty prompt", func() {
		render(func(c *ui.Context) { Suggestions(c, "x", []string{""}, &chosen) })
	})
}

func TestPromptLibrary(t *testing.T) {
	prompts := []string{"Deploy to production", "Draft a release note", "Review a pull request"}
	query := ""
	chosen, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _ := PromptLibrary(c, "Library", prompts, &query, &chosen, PromptLibraryOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "library")
	tt.Key(0, ui.KeyTab)
	tt.Type("draft")
	tt.Frame()
	settles(t, tt, "hits")
	if tt.HasText("Deploy to production") || !tt.HasText("Draft a release note") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Draft a release note")
	if chosen != 1 || changed != 1 {
		t.Errorf("chosen %d changed %d", chosen, changed)
	}
	cues(t, before, sound.Select)
	query = "zzz"
	tt.Frame()
	settles(t, tt, "empty")
	if !tt.HasText("No prompts match") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "nil query", func() {
		render(func(c *ui.Context) { PromptLibrary(c, "x", prompts, nil, &chosen, PromptLibraryOptions{}) })
	})
}

func TestBranches(t *testing.T) {
	current, total, backs, forths := 2, 3, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, back, forth := Branches(c, &current, total, BranchesOptions{})
		if back.Clicked() {
			backs++
		}
		if forth.Clicked() {
			forths++
		}
	}), 400, 200)
	settles(t, tt, "branches")
	if !tt.HasText("2 / 3") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Next")
	if current != 3 || forths != 1 {
		t.Fatalf("current %d forths %d", current, forths)
	}
	click(t, tt, "Next")
	if current != 3 {
		t.Error("moved past the last branch")
	}
	click(t, tt, "Previous")
	click(t, tt, "Previous")
	if current != 1 || backs != 2 {
		t.Errorf("current %d backs %d", current, backs)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	panics(t, "zero current", func() {
		n := 0
		render(func(c *ui.Context) { Branches(c, &n, 3, BranchesOptions{}) })
	})
	panics(t, "zero total", func() {
		n := 1
		render(func(c *ui.Context) { Branches(c, &n, 0, BranchesOptions{}) })
	})
}

func TestAnswerActions(t *testing.T) {
	vote, changed, retries := 0, 0, 0
	copies := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, cp, retry := AnswerActions(c, "The deck is ready.", &vote, AnswerActionsOptions{})
		if root.Changed() {
			changed++
		}
		if cp.Clicked() {
			copies++
		}
		if retry.Clicked() {
			retries++
		}
	}), 500, 200)
	settles(t, tt, "actions")
	before := baseline()
	click(t, tt, "Like")
	if vote != 1 || changed != 1 {
		t.Fatalf("vote %d changed %d", vote, changed)
	}
	click(t, tt, "Like")
	if vote != 0 {
		t.Errorf("vote %d after un-like", vote)
	}
	click(t, tt, "Dislike")
	if vote != -1 {
		t.Errorf("vote %d", vote)
	}
	click(t, tt, "Retry")
	if retries != 1 {
		t.Errorf("retries %d", retries)
	}
	click(t, tt, "Copy")
	if tt.Clipboard() != "The deck is ready." || copies != 1 {
		t.Errorf("clipboard %q copies %d", tt.Clipboard(), copies)
	}
	cues(t, before, sound.Pop, sound.Pop, sound.Pop, sound.Pop, sound.Pop)
	panics(t, "bad vote", func() {
		v := 2
		render(func(c *ui.Context) { AnswerActions(c, "x", &v, AnswerActionsOptions{}) })
	})
}

func TestArtifact(t *testing.T) {
	opened, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, open := Artifact(c, "deck.html", "Web page", ArtifactOptions{})
		if root.Changed() {
			changed++
		}
		if open.Clicked() {
			opened++
		}
	}), 500, 300)
	settles(t, tt, "artifact")
	before := baseline()
	click(t, tt, "Open")
	if opened != 1 || changed != 1 {
		t.Errorf("opened %d changed %d", opened, changed)
	}
	cues(t, before, sound.Open)
	panics(t, "empty kind", func() { render(func(c *ui.Context) { Artifact(c, "x", "", ArtifactOptions{}) }) })
}

func TestApproval(t *testing.T) {
	allowed, denied, changed := 0, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, allow, deny := Approval(c, "git push --force origin main", ApprovalOptions{})
		if root.Changed() {
			changed++
		}
		if allow.Clicked() {
			allowed++
		}
		if deny.Clicked() {
			denied++
		}
	}), 500, 300)
	settles(t, tt, "approval")
	if !tt.HasText("Allow this action?") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Allow")
	if allowed != 1 || changed != 1 {
		t.Fatalf("allowed %d changed %d", allowed, changed)
	}
	click(t, tt, "Deny")
	if denied != 1 || changed != 2 {
		t.Errorf("denied %d changed %d", denied, changed)
	}
	cues(t, before, sound.Success, sound.Close)
	panics(t, "empty action", func() { render(func(c *ui.Context) { Approval(c, "", ApprovalOptions{}) }) })
}
