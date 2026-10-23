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
		if Breadcrumb(c, parts, &at).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "crumbs")
	before := baseline()
	click(t, tt, "gallery")
	if at != 1 || changed != 1 {
		t.Fatalf("at %d changed %d", at, changed)
	}
	cues(t, before, sound.Tap)
	panics(t, "no parts", func() {
		render(func(c *ui.Context) { Breadcrumb(c, nil, &at) })
	})
}

func TestCommit(t *testing.T) {
	copied := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, copy := Commit(c, "1849e62abc", "Add code components", "Mei", time.Now().Add(-2*time.Hour), CommitOptions{})
		if copy.Clicked() {
			copied++
		}
	}), 500, 200)
	settles(t, tt, "commit")
	if !tt.HasText("1849e62") || !tt.HasText("2h") {
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
			CommitOptions{Copy: "Copier", FormatAge: func(time.Duration) string { return "vor 2h" }})
	}), 500, 200)
	settles(t, custom, "custom commit")
	if !custom.HasText("Copier") || !custom.HasText("vor 2h") {
		t.Errorf("custom copy %q", custom.Texts())
	}
}

func TestPullRequest(t *testing.T) {
	checks := []Check{{"build", CheckPass}, {"test", CheckPending}, {"lint", CheckFail}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PullRequest(c, "Add the code category", PROpen, checks, []string{"Ana Reyes", "Kai Mori"}, PullRequestOptions{})
	}), 500, 300)
	tt.Frame()
	if !tt.HasText("Open") || !tt.HasText("1 of 3 passing") {
		t.Errorf("texts %q", tt.Texts())
	}
	if !tt.FramePending() {
		t.Error("pending check does not pulse")
	}
	panics(t, "bad state", func() {
		render(func(c *ui.Context) { PullRequest(c, "x", prStateCount, nil, nil, PullRequestOptions{}) })
	})

	custom := ui.NewTester(padded(func(c *ui.Context) {
		done := []Check{{Name: "build", State: CheckPass}, {Name: "test", State: CheckPass}, {Name: "lint", State: CheckPass}}
		PullRequest(c, "Titre", PRMerged, done, nil, PullRequestOptions{Merged: "Fusionne", Passing: "%d sur %d"})
	}), 500, 300)
	settles(t, custom, "custom pr")
	if !custom.HasText("Fusionne") || !custom.HasText("3 sur 3") {
		t.Errorf("custom copy %q", custom.Texts())
	}
}

func TestChecks(t *testing.T) {
	checks := []Check{{"build", CheckPass}, {"deploy", CheckPending}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Checks(c, "CI", checks)
	}), 500, 200)
	if !moves(tt, 700*time.Millisecond) {
		t.Error("pending check does not spin")
	}
	if !tt.HasText("build") || !tt.HasText("deploy") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no checks", func() {
		render(func(c *ui.Context) { Checks(c, "x", nil) })
	})
}

func TestCommitGraph(t *testing.T) {
	commits := []GraphCommit{
		{"aaaaaaa1", "Merge the branch", 0},
		{"bbbbbbb2", "Work on the branch", 1},
		{"ccccccc3", "Start the branch", 0},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		CommitGraph(c, "Graph", commits)
	}), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("rows do not arrive")
	}
	settles(t, tt, "graph")
	if !tt.HasText("aaaaaaa") || !tt.HasText("Work on the branch") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "lane too big", func() {
		render(func(c *ui.Context) { CommitGraph(c, "x", []GraphCommit{{Hash: "a", Message: "m", Lane: 8}}) })
	})

	dup := []GraphCommit{{Hash: "abcdefg1", Message: "first"}, {Hash: "abcdefg2", Message: "second"}}
	same := ui.NewTester(padded(func(c *ui.Context) {
		CommitGraph(c, "Graph", dup)
	}), 500, 300)
	settles(t, same, "same prefix")
	if !same.HasText("first") || !same.HasText("second") {
		t.Errorf("same-prefix commits %q", same.Texts())
	}
}

func TestBlame(t *testing.T) {
	lines := []BlameLine{
		{"Mei", "package main", time.Now().Add(-48 * time.Hour)},
		{"Kai", "func main() {}", time.Now().Add(-2 * time.Hour)},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Blame(c, "Blame", lines, BlameOptions{})
	}), 500, 200)
	settles(t, tt, "blame")
	if !tt.HasText("2d") || !tt.HasText("2h") || !tt.HasText("func main() {}") {
		t.Errorf("texts %q", tt.Texts())
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
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Language(c, "Languages", shares)
	}), 500, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("segments do not grow")
	}
	settles(t, tt, "language")
	if !tt.HasText("Go 72%") || !tt.HasText("Shell 18%") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "over one", func() {
		render(func(c *ui.Context) { Language(c, "x", []LangShare{{"a", 0.7}, {"b", 0.7}}) })
	})
}
