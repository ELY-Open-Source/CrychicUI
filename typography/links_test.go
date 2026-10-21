package typography

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestTextLink(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if TextLink(c, "Documentation", "https://example.com/docs").Clicked() {
			clicks++
		}
	}, 400, 200)
	r := find(t, tt, "Documentation")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "TextLink underline")
	taps := count(sound.Tap)
	if err := tt.Click("Documentation"); err != nil {
		t.Fatal(err)
	}
	opened := tt.OpenedURLs()
	if clicks != 1 || count(sound.Tap) != taps+1 || len(opened) != 1 || opened[0] != "https://example.com/docs" {
		t.Errorf("clicks %d, cues %d, opened %q", clicks, count(sound.Tap)-taps, opened)
	}
}

func TestExternalLink(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ExternalLink(c, "Repository", "https://example.com/repo")
	}, 400, 200)
	r := find(t, tt, "Repository")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "ExternalLink hover")
	taps := count(sound.Tap)
	if err := tt.Click("Repository"); err != nil {
		t.Fatal(err)
	}
	opened := tt.OpenedURLs()
	if count(sound.Tap) != taps+1 || len(opened) != 1 || opened[0] != "https://example.com/repo" {
		t.Errorf("cues %d, opened %q", count(sound.Tap)-taps, opened)
	}
}

func TestAnchorHeading(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		_, copyAnchor := AnchorHeading(c, "Introduction", "intro")
		if copyAnchor.Clicked() {
			clicks++
		}
	}, 400, 200)
	settles(t, tt, "AnchorHeading")
	r := find(t, tt, "Introduction")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "AnchorHeading hash")
	pops := count(sound.Pop)
	if err := tt.Click("Copy anchor"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 || tt.Clipboard() != "#intro" || count(sound.Pop) != pops+1 {
		t.Errorf("clicks %d, clipboard %q, cues %d", clicks, tt.Clipboard(), count(sound.Pop)-pops)
	}
	panics(t, "AnchorHeading anchor", func() { ui.Render(func(c *ui.Context) { AnchorHeading(c, "Title", "") }, 10, 10, 1) })
}

func TestFootnote(t *testing.T) {
	open, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Footnote(c, 1, "A supporting note", &open).Changed() {
			changed++
		}
	}, 400, 200)
	opens := count(sound.Open)
	if err := tt.Click("Footnote 1"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(2 * time.Second)
	if !open || changed != 1 || count(sound.Open) != opens+1 || !tt.HasText("A supporting note") {
		t.Errorf("open %v, changed %d, cues %d", open, changed, count(sound.Open)-opens)
	}
	tt.ClickAt(390, 190)
	tt.Advance(2 * time.Second)
	if open || changed != 2 || tt.HasText("A supporting note") {
		t.Errorf("dismissed %v, changed %d, texts %q", open, changed, tt.Texts())
	}
}

func TestLinkValidation(t *testing.T) {
	panics(t, "TextLink label", func() { ui.Render(func(c *ui.Context) { TextLink(c, "", "https://example.com") }, 10, 10, 1) })
	panics(t, "TextLink URL", func() { ui.Render(func(c *ui.Context) { TextLink(c, "Link", "") }, 10, 10, 1) })
	panics(t, "Footnote number", func() { ui.Render(func(c *ui.Context) { Footnote(c, 0, "Note", new(bool)) }, 10, 10, 1) })
	panics(t, "Footnote open", func() { ui.Render(func(c *ui.Context) { Footnote(c, 1, "Note", nil) }, 10, 10, 1) })
}
