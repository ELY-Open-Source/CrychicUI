package agent

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestSuggestions(t *testing.T) {
	draft, changed := "", 0
	prompts := []Suggestion{{Text: "Summarize the thread", Icon: "list"}, {Text: "Draft a reply", Icon: "pen-line"}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Suggestions(c, "Starters", prompts, &draft).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "suggestions")
	before := baseline()
	click(t, tt, "Draft a reply")
	if draft != "Draft a reply" || changed != 1 {
		t.Fatalf("draft %q changed %d", draft, changed)
	}
	click(t, tt, "Draft a reply")
	if changed != 1 {
		t.Error("refilling the same prompt reported a change")
	}
	cues(t, before, sound.Tap)
	long := []Suggestion{{Text: "Summarize every message in this thread and list each open decision with its owner", Icon: "list"}}
	narrow := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(240).Children(func() { Suggestions(c, "Long", long, &draft) })
	}, 400, 200)
	narrow.Frame()
	if r := find(t, narrow, long[0].Text); r.X+r.W > 240 {
		t.Errorf("long chip overflows to %v", r.X+r.W)
	}
	panics(t, "no icon", func() {
		render(func(c *ui.Context) { Suggestions(c, "S", []Suggestion{{Text: "x"}}, &draft) })
	})
}

func TestPromptLibrary(t *testing.T) {
	query, chosen, changed := "", -1, 0
	prompts := []SavedPrompt{
		{Title: "Release note", Body: "Write a release note from these commits.", Tag: "Writing"},
		{Title: "Code review", Body: "Review this pull request for bugs.", Tag: "Engineering"},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if root, _ := PromptLibrary(c, "Library", prompts, &query, &chosen, PromptLibraryOptions{}); root.Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "library")
	if !tt.HasText("Engineering") || !tt.HasText("Review this pull request for bugs.") {
		t.Fatalf("texts %q", tt.Texts())
	}
	query = "pull request"
	tt.Frame()
	if tt.HasText("Release note") || !tt.HasText("Code review") {
		t.Errorf("body search %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Code review")
	if chosen != 1 || changed != 1 {
		t.Errorf("chosen %d changed %d", chosen, changed)
	}
	cues(t, before, sound.Select)
	query = "zzz"
	tt.Frame()
	if !tt.HasText("No prompts match") {
		t.Errorf("empty %q", tt.Texts())
	}
	reordered := ui.NewTester(padded(func(c *ui.Context) {
		PromptLibrary(c, "Library", prompts, &query, &chosen, PromptLibraryOptions{})
	}), 500, 400)
	query, chosen = "", -1
	settles(t, reordered, "library")
	r := find(t, reordered, "Release note")
	reordered.Press(r.X+r.W/2, r.Y+r.H/2)
	reordered.Frame()
	prompts[0], prompts[1] = prompts[1], prompts[0]
	reordered.Frame()
	r = find(t, reordered, "Release note")
	reordered.Release(r.X+r.W/2, r.Y+r.H/2)
	reordered.Frame()
	if chosen != 1 {
		t.Errorf("chosen %d, want Release note at 1", chosen)
	}
	panics(t, "no body", func() {
		render(func(c *ui.Context) {
			PromptLibrary(c, "L", []SavedPrompt{{Title: "x"}}, &query, &chosen, PromptLibraryOptions{})
		})
	})
}

func TestArtifact(t *testing.T) {
	opened := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if _, open := Artifact(c, "launch-deck.html", Web, "HTML · 18 KB", ArtifactOptions{}); open.Clicked() {
			opened++
		}
	}), 500, 300)
	settles(t, tt, "artifact")
	if !tt.HasText("HTML · 18 KB") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Open")
	if opened != 1 {
		t.Errorf("opened %d", opened)
	}
	cues(t, before, sound.Open)
	panics(t, "bad kind", func() {
		render(func(c *ui.Context) { Artifact(c, "x", Web+1, "", ArtifactOptions{}) })
	})
}

func TestApproval(t *testing.T) {
	for _, tc := range []struct {
		button string
		want   Decision
		cue    sound.Cue
		shown  string
	}{{"Allow", Allowed, sound.Success, "Allowed"}, {"Deny", Denied, sound.Close, "Denied"}} {
		decision, changed := Undecided, 0
		tt := ui.NewTester(padded(func(c *ui.Context) {
			if root, _, _ := Approval(c, "shell", "git push origin main", "Publishes 3 commits.", &decision, ApprovalOptions{}); root.Changed() {
				changed++
			}
		}), 500, 400)
		settles(t, tt, "approval")
		before := baseline()
		click(t, tt, tc.button)
		if decision != tc.want || changed != 1 {
			t.Fatalf("%s: decision %d changed %d", tc.button, decision, changed)
		}
		settles(t, tt, "decided approval")
		if !tt.HasText(tc.shown) || tt.HasText("Deny") {
			t.Errorf("%s: texts %q", tc.button, tt.Texts())
		}
		cues(t, before, tc.cue)
	}
	panics(t, "nil decision", func() {
		render(func(c *ui.Context) { Approval(c, "shell", "ls", "", nil, ApprovalOptions{}) })
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
