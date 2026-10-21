package typography

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestLists(t *testing.T) {
	t.Run("BulletList", func(t *testing.T) {
		tt := ui.NewTester(func(c *ui.Context) { BulletList(c, []string{"First", "Second"}) }, 300, 200)
		settles(t, tt, "BulletList")
		find(t, tt, "First")
		find(t, tt, "Second")
	})
	t.Run("NumberedList", func(t *testing.T) {
		tt := ui.NewTester(func(c *ui.Context) { NumberedList(c, []string{"First", "Second"}) }, 300, 200)
		settles(t, tt, "NumberedList")
		find(t, tt, "1.")
		find(t, tt, "Second")
	})
	t.Run("DefinitionList", func(t *testing.T) {
		tt := ui.NewTester(func(c *ui.Context) {
			DefinitionList(c, []Definition{{Term: "Name", Description: "Crychic"}})
		}, 400, 200)
		settles(t, tt, "DefinitionList")
		find(t, tt, "Name")
		find(t, tt, "Crychic")
	})
	t.Run("Prose", func(t *testing.T) {
		tt := ui.NewTester(func(c *ui.Context) { Prose(c, []string{"One paragraph", "Two paragraphs"}) }, 300, 200)
		settles(t, tt, "Prose")
		find(t, tt, "One paragraph")
		find(t, tt, "Two paragraphs")
	})
}

func TestDiffText(t *testing.T) {
	deleted, inserted := "remove me", "add me"
	tt := ui.NewTester(func(c *ui.Context) { DiffText(c, deleted, inserted) }, 300, 200)
	deleted, inserted = "gone", "added"
	tt.Frame()
	settles(t, tt, "DiffText")
	if !tt.HasText("gone") || !tt.HasText("added") || tt.HasText("remove me") || tt.HasText("add me") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestDiffTextAppearance(t *testing.T) {
	inserted := ""
	tt := ui.NewTester(func(c *ui.Context) { DiffText(c, "old", inserted) }, 300, 200)
	inserted = "new"
	tt.Frame()
	if !tt.FramePending() {
		t.Error("new insertion did not fade in")
	}
	tt.Advance(2 * time.Second)
	find(t, tt, "new")
}

func TestMonogram(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Monogram(c, "Jane doe") }, 300, 200)
	settles(t, tt, "Monogram")
	find(t, tt, "JD")
	panics(t, "Monogram empty", func() { ui.Render(func(c *ui.Context) { Monogram(c, " ") }, 10, 10, 1) })
}

func TestRuby(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Ruby(c, "base", "reading") }, 300, 200)
	settles(t, tt, "Ruby")
	if find(t, tt, "reading").Y >= find(t, tt, "base").Y {
		t.Error("annotation is not above the base")
	}
}

func TestLoadingDots(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { LoadingDots(c, "Loading") }, 300, 200)
	find(t, tt, "Loading")
	if !tt.FramePending() {
		t.Error("dots did not start their loop")
	}
	tt.Advance(300 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("dots stopped their loop")
	}
}
