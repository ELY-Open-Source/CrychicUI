package code

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestBreadcrumb(t *testing.T) {
	parts := []string{"cmd", "gallery", "demos_code.go", "demoCodeBlock"}
	at := 3
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Breadcrumb(c, "Path", parts, &at).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "crumbs")
	before := baseline()
	click(t, tt, "gallery")
	if at != 1 || changed != 1 {
		t.Fatalf("at %d changed %d", at, changed)
	}
	click(t, tt, "gallery")
	if changed != 1 {
		t.Errorf("the current crumb fired: changed %d", changed)
	}
	cues(t, before, sound.Tap)
	panics(t, "no parts", func() {
		render(func(c *ui.Context) { Breadcrumb(c, "Path", nil, &at) })
	})
}

func TestCommit(t *testing.T) {
	copied := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, copy := Commit(c, "1849e62abc", "Add code components", "Mei Chen", time.Now().Add(-2*time.Hour), CommitOptions{})
		if copy.Clicked() {
			copied++
		}
	}), 500, 200)
	settles(t, tt, "commit")
	if !tt.HasText("1849e62") || !tt.HasText("Mei Chen committed 2h") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Copy hash")
	if tt.Clipboard() != "1849e62abc" || copied != 1 {
		t.Fatalf("clipboard %q copied %d", tt.Clipboard(), copied)
	}
	cues(t, before, sound.Pop)
	panics(t, "zero time", func() {
		render(func(c *ui.Context) { Commit(c, "abc", "m", "a", time.Time{}, CommitOptions{}) })
	})

	custom := ui.NewTester(padded(func(c *ui.Context) {
		Commit(c, "1849e62abc", "Add code components", "Mei", time.Now().Add(-2*time.Hour),
			CommitOptions{Copy: "Copier", Committed: "%s a validé %s", FormatAge: func(time.Duration) string { return "il y a 2h" }})
	}), 500, 200)
	settles(t, custom, "custom commit")
	if !custom.HasText("Copier") || !custom.HasText("Mei a validé il y a 2h") {
		t.Errorf("custom copy %q", custom.Texts())
	}
}

func demoPR(state PRState, checks []Check) PR {
	return PR{Number: 1842, Title: "Add the code category", Author: "Ana Reyes", Head: "feature/code", Base: "main",
		State: state, Added: 120, Removed: 45, Checks: checks, Reviewers: []string{"Ana Reyes", "Kai Mori"}}
}

func TestPullRequest(t *testing.T) {
	checks := []Check{{Name: "build", State: CheckPass}, {Name: "test", State: CheckPending}, {Name: "lint", State: CheckFail}}
	tt := ui.NewTester(padded(func(c *ui.Context) { PullRequest(c, demoPR(PROpen, checks), PullRequestOptions{}) }), 500, 300)
	tt.Frame()
	if !tt.HasText("Open") || !tt.HasText("1 of 3 passing") || !tt.HasText("#1842") || !tt.HasText("feature/code") {
		t.Errorf("texts %q", tt.Texts())
	}
	tt.Advance(2 * time.Second)
	if !moves(tt, 400*time.Millisecond) {
		t.Error("pending check does not pulse")
	}
	panics(t, "bad state", func() {
		render(func(c *ui.Context) { PullRequest(c, demoPR(prStateCount, nil), PullRequestOptions{}) })
	})
	panics(t, "no number", func() {
		pr := demoPR(PROpen, nil)
		pr.Number = 0
		render(func(c *ui.Context) { PullRequest(c, pr, PullRequestOptions{}) })
	})

	custom := ui.NewTester(padded(func(c *ui.Context) {
		done := []Check{{Name: "build", State: CheckPass}, {Name: "test", State: CheckPass}, {Name: "lint", State: CheckPass}}
		PullRequest(c, demoPR(PRMerged, done), PullRequestOptions{Merged: "Fusionne", Passing: "%d sur %d"})
	}), 500, 300)
	settles(t, custom, "custom pr")
	if !custom.HasText("Fusionne") || !custom.HasText("3 sur 3") {
		t.Errorf("custom copy %q", custom.Texts())
	}
}

func TestChecks(t *testing.T) {
	checks := []Check{{Name: "build", State: CheckPass, Duration: 72 * time.Second}, {Name: "deploy", State: CheckPending}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Checks(c, "CI", checks, ChecksOptions{}) }), 500, 200)
	if !moves(tt, 700*time.Millisecond) {
		t.Error("pending check does not spin")
	}
	if !tt.HasText("build") || !tt.HasText("1m 12s") || !tt.HasText("Running…") || !tt.HasText("1 of 2 passing") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no checks", func() { render(func(c *ui.Context) { Checks(c, "x", nil, ChecksOptions{}) }) })
	panics(t, "twin checks", func() {
		render(func(c *ui.Context) { Checks(c, "x", []Check{{Name: "a"}, {Name: "a"}}, ChecksOptions{}) })
	})
}

func TestTook(t *testing.T) {
	for d, want := range map[time.Duration]string{800 * time.Millisecond: "0.8s", 45 * time.Second: "45s", 72 * time.Second: "1m 12s"} {
		if got := took(d); got != want {
			t.Errorf("took(%v) = %q, want %q", d, got, want)
		}
	}
}

func TestCommitGraph(t *testing.T) {
	commits := []GraphCommit{
		{Hash: "a1b2c3d4", Message: "Merge the branch", Lane: 0, Parents: []string{"c3d4e5f6", "b2c3d4e5"}, Refs: []string{"main"}},
		{Hash: "b2c3d4e5", Message: "Work on the branch", Lane: 1, Parents: []string{"c3d4e5f6"}},
		{Hash: "c3d4e5f6", Message: "Start the branch", Lane: 0},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { CommitGraph(c, "Graph", commits) }), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("rows do not arrive")
	}
	settles(t, tt, "graph")
	if !tt.HasText("a1b2c3d") || !tt.HasText("main") || !tt.HasText("Work on the branch") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "lane too big", func() {
		render(func(c *ui.Context) { CommitGraph(c, "x", []GraphCommit{{Hash: "a", Message: "m", Lane: 8}}) })
	})
	panics(t, "parent above", func() {
		render(func(c *ui.Context) {
			CommitGraph(c, "x", []GraphCommit{{Hash: "a", Message: "m"}, {Hash: "b", Message: "n", Parents: []string{"a"}}})
		})
	})

	dup := []GraphCommit{{Hash: "abcdefg1", Message: "first"}, {Hash: "abcdefg2", Message: "second"}}
	same := ui.NewTester(padded(func(c *ui.Context) { CommitGraph(c, "Graph", dup) }), 500, 300)
	settles(t, same, "same prefix")
	if !same.HasText("first") || !same.HasText("second") {
		t.Errorf("same-prefix commits %q", same.Texts())
	}
}

func TestBlame(t *testing.T) {
	now := time.Now()
	lines := []BlameLine{
		{Author: "Mei", Text: "package main", At: now.Add(-48 * time.Hour)},
		{Author: "Mei", Text: "", At: now.Add(-48 * time.Hour)},
		{Author: "Kai", Text: "func main() {}", At: now.Add(-2 * time.Hour)},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { Blame(c, "Blame", lines, BlameOptions{Lang: "go"}) }), 500, 200)
	settles(t, tt, "blame")
	if !tt.HasText("2d") || !tt.HasText("2h") || !tt.HasText("func main() {}") {
		t.Errorf("texts %q", tt.Texts())
	}
	meis := 0
	for _, s := range tt.Texts() {
		if s == "Mei" {
			meis++
		}
	}
	if meis != 1 {
		t.Errorf("a run of lines repeats its author %d times", meis)
	}
	panics(t, "no author", func() {
		render(func(c *ui.Context) {
			Blame(c, "x", []BlameLine{{Author: "", Text: "y", At: time.Now()}}, BlameOptions{})
		})
	})

	arrive := ui.NewTester(padded(func(c *ui.Context) {
		Blame(c, "Blame", lines, BlameOptions{FormatAge: func(time.Duration) string { return "alt" }})
	}), 500, 200)
	arrive.Frame()
	if !arrive.FramePending() {
		t.Error("rows do not arrive")
	}
	settles(t, arrive, "arrive")
	if !arrive.HasText("alt") {
		t.Errorf("custom age %q", arrive.Texts())
	}
}

func TestLanguageClipsAnimatedSegments(t *testing.T) {
	shares := []LangShare{{"Go", 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Language(c, "Languages", shares) }), 500, 200)
	settles(t, tt, "one language")
	r := find(t, tt, "Languages")
	x, y := int(r.X+r.W+2), int(r.Y+4)
	outside := tt.Image().RGBAAt(x, y)
	shares = []LangShare{{"Go", 0.5}, {"JS", 0.5}}
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	if got := tt.Image().RGBAAt(x, y); got != outside {
		t.Errorf("segments paint outside bar: got %v want %v", got, outside)
	}
}

func TestLanguage(t *testing.T) {
	shares := []LangShare{{"Go", 0.72}, {"Shell", 0.18}, {"JS", 0.1}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Language(c, "Languages", shares) }), 500, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("segments do not grow")
	}
	settles(t, tt, "language")
	if !tt.HasText("Go 72%") || !tt.HasText("Shell 18%") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "over one", func() { render(func(c *ui.Context) { Language(c, "x", []LangShare{{"a", 0.7}, {"b", 0.7}}) }) })
	panics(t, "twin names", func() { render(func(c *ui.Context) { Language(c, "x", []LangShare{{"a", 0.1}, {"a", 0.1}}) }) })
}
