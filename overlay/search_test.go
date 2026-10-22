package overlay

import (
	"image"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCommandPalette(t *testing.T) {
	commands := []MenuItem{{"Open file", "", "Ctrl O"}, {"Save file", "", ""}, {"Save all", "", ""}}
	chosen, open := -1, false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Text(c, "Editor")
		CommandPalette(c, "Commands", commands, &chosen, &open, CommandPaletteOptions{})
	}, 600, 500)
	before := baseline()
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Frame()
	if !open || !tt.HasText("Open file") {
		t.Fatalf("Cmd+K: open %v, %q", open, tt.Texts())
	}
	tt.Type("save")
	tt.Frame()
	if tt.HasText("Open file") {
		t.Fatalf("filter %q", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 2 || open {
		t.Fatalf("chose %d, open %v", chosen, open)
	}
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if chosen != 0 {
		t.Errorf("reopened palette kept its query or highlight: chose %d", chosen)
	}
	cues(t, before, sound.Open, sound.Select, sound.Open, sound.Select)
	tt.Key(ui.Cmd, ui.KeyK)
	tt.Type("zzz")
	tt.Frame()
	if !tt.HasText("No matching commands") {
		t.Errorf("empty %q", tt.Texts())
	}
}

func TestSpotlightSearch(t *testing.T) {
	groups := []ResultGroup{{"Files", []string{"main.go", "notes.md"}}, {"People", []string{"Ada", "Alan"}}}
	chosen, open, changed := [2]int{-1, -1}, true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SpotlightSearch(c, "Search", groups, &chosen, &open, SpotlightSearchOptions{}).Changed() {
			changed++
		}
	}, 800, 600)
	if !tt.HasText("Files") || !tt.HasText("People") {
		t.Fatalf("groups %q", tt.Texts())
	}
	before := baseline()
	tt.Type("a")
	tt.Frame()
	if tt.HasText("notes.md") || !tt.HasText("main.go") {
		t.Fatalf("filter a %q", tt.Texts())
	}
	tt.Type("d")
	tt.Frame()
	if tt.HasText("Files") || !tt.HasText("Ada") {
		t.Fatalf("filter ad %q", tt.Texts())
	}
	click(t, tt, "Ada")
	if chosen != [2]int{1, 0} || open || changed != 1 {
		t.Errorf("chose %v, open %v, changed %d", chosen, open, changed)
	}
	cues(t, before, sound.Select)
}

func TestLightbox(t *testing.T) {
	imgs := []*ui.Bitmap{ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 8, 8))), ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 8, 8)))}
	at, open, changed := 0, true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Lightbox(c, "Photos", imgs, []string{"Beach", "Hills"}, &at, &open, LightboxOptions{}).Changed() {
			changed++
		}
	}, 800, 600)
	before := baseline()
	click(t, tt, "Next")
	if at != 1 || !tt.FramePending() || !tt.HasText("Hills") {
		t.Fatalf("next: at %d, %q", at, tt.Texts())
	}
	settles(t, tt, "slide")
	tt.Key(0, ui.KeyLeft)
	tt.Frame()
	if at != 0 {
		t.Errorf("Left: at %d", at)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || changed != 3 {
		t.Errorf("Escape: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Swipe, sound.Swipe)
	panics(t, "short captions", func() {
		render(func(c *ui.Context) {
			a, o := 0, true
			Lightbox(c, "x", imgs, []string{"one"}, &a, &o, LightboxOptions{})
		})
	})
}
