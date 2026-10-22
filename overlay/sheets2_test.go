package overlay

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// pull drags from a point in frame steps, then releases.
func pull(tt *ui.Tester, x, y, dx, dy float32, steps int) {
	tt.Press(x, y)
	for i := 1; i <= steps; i++ {
		tt.Move(x+dx*float32(i)/float32(steps), y+dy*float32(i)/float32(steps))
		tt.Advance(16 * time.Millisecond)
	}
	tt.Release(x+dx, y+dy)
	tt.Frame()
}

func TestShareSheet(t *testing.T) {
	targets := []ShareTarget{{"Mail", "mail"}, {"Messages", "message-circle"}}
	chosen, open, changed := -1, true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ShareSheet(c, "Share", "https://crychic.dev/a", targets, &chosen, &open, ShareSheetOptions{}).Changed() {
			changed++
		}
	}, 500, 500)
	settles(t, tt, "sheet")
	before := baseline()
	click(t, tt, "Copy link")
	tt.Frame()
	if tt.Clipboard() != "https://crychic.dev/a" || !tt.HasText("Copied") || changed != 1 {
		t.Fatalf("copy: clipboard %q, changed %d, %q", tt.Clipboard(), changed, tt.Texts())
	}
	tt.Advance(2 * time.Second)
	if tt.HasText("Copied") {
		t.Error("Copied stayed")
	}
	click(t, tt, "Messages")
	if chosen != 1 || open || changed != 2 {
		t.Errorf("share: chosen %d, open %v, changed %d", chosen, open, changed)
	}
	cues(t, before, sound.Pop, sound.Pop)
	panics(t, "no link", func() {
		render(func(c *ui.Context) { ch, o := -1, true; ShareSheet(c, "x", "", targets, &ch, &o, ShareSheetOptions{}) })
	})
}

func TestActionSheet(t *testing.T) {
	actions := []MenuItem{{"Archive", "archive", ""}, {"Delete", "trash-2", ""}}
	chosen, open := -1, true
	tt := ui.NewTester(func(c *ui.Context) {
		ActionSheet(c, "Message", actions, &chosen, &open, ActionSheetOptions{})
	}, 500, 500)
	settles(t, tt, "sheet")
	before := baseline()
	click(t, tt, "Delete")
	if chosen != 1 || open {
		t.Fatalf("pick: chosen %d, open %v", chosen, open)
	}
	open = true
	tt.Frame()
	settles(t, tt, "reopened")
	click(t, tt, "Cancel")
	if chosen != 1 || open {
		t.Errorf("cancel: chosen %d, open %v", chosen, open)
	}
	cues(t, before, sound.Select, sound.Close)
}

func TestMiniPlayer(t *testing.T) {
	th := ui.LightTheme()
	playing, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if MiniPlayer(c, "Song", "Band", &playing, 0.5, MiniPlayerOptions{}).Changed() {
			changed++
		}
	}, 400, 200)
	if !tt.FramePending() {
		t.Fatal("no slide up")
	}
	settles(t, tt, "player")
	r := find(t, tt, "Song")
	// The bar sits on the bottom border's row above.
	lit, y := 0, int(r.Y+r.H)-2
	for x := int(r.X) + 8; x < int(r.X+r.W)-8; x++ {
		if p := tt.Image().RGBAAt(x, y); p.R == th.Accent.R && p.G == th.Accent.G && p.B == th.Accent.B {
			lit++
		}
	}
	if w := int(r.W) - 16; lit < w*4/10 || lit > w*6/10 {
		t.Errorf("progress lit %d of %d px, want half", lit, w)
	}
	before := baseline()
	click(t, tt, "Play")
	click(t, tt, "Pause")
	if playing || changed != 2 {
		t.Errorf("toggle twice: playing %v, changed %d", playing, changed)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
	panics(t, "progress 2", func() { render(func(c *ui.Context) { p := false; MiniPlayer(c, "x", "", &p, 2, MiniPlayerOptions{}) }) })
}

func TestFloatingPanel(t *testing.T) {
	at, changed := ui.Point{}, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if FloatingPanel(c, "Tools", &at, FloatingPanelOptions{Height: 300}, func() { ui.Text(c, "Brush") }).Changed() {
			changed++
		}
	}, 400, 300)
	settles(t, tt, "panel")
	r := find(t, tt, "Brush")
	before := baseline()
	pull(tt, r.X+4, r.Y+4, 100, 50, 10)
	settles(t, tt, "moved")
	if at.X < 98 || at.X > 102 || at.Y < 48 || at.Y > 52 || changed != 1 {
		t.Fatalf("drag: at %+v, changed %d", at, changed)
	}
	pull(tt, r.X+104, r.Y+54, 600, 0, 10)
	settles(t, tt, "sprung back")
	if edge := 400 - ui.LightTheme().Space(60); at.X < edge-1 || at.X > edge+1 {
		t.Errorf("past the edge: at %+v, want x %v", at, edge)
	}
	cues(t, before, sound.Drop, sound.Drop)
}

func TestPiPBox(t *testing.T) {
	corner, changed := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if PiPBox(c, "Call", &corner, PiPBoxOptions{Height: 300}, func() { ui.Text(c, "Ada") }).Changed() {
			changed++
		}
	}, 400, 300)
	settles(t, tt, "box")
	before := baseline()
	pull(tt, 92, 60, 200, 170, 10)
	settles(t, tt, "thrown")
	if corner != 3 || changed != 1 {
		t.Fatalf("drag to bottom right: corner %d, changed %d", corner, changed)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	tt.Frame()
	if corner != 2 {
		t.Errorf("Left: corner %d", corner)
	}
	cues(t, before, sound.Drop, sound.Drop)
	panics(t, "corner 4", func() { render(func(c *ui.Context) { k := 4; PiPBox(c, "x", &k, PiPBoxOptions{}, func() {}) }) })
}

func TestSheetStack(t *testing.T) {
	titles := []string{"Settings", "Account", "Password"}
	depth, changed := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SheetStack(c, titles, &depth, func(i int) {
			ui.Text(c, titles[i]+" body")
			if i+1 < len(titles) && ui.ButtonBase(c).Label("Deeper "+titles[i]).Size(120, 32).Clicked() {
				depth++
			}
		}, SheetStackOptions{}).Changed() {
			changed++
		}
	}, 500, 600)
	before := baseline()
	depth = 1
	tt.Frame()
	settles(t, tt, "first sheet")
	click(t, tt, "Deeper Settings")
	if depth != 2 || !tt.FramePending() {
		t.Fatalf("push: depth %d, moving %v", depth, tt.FramePending())
	}
	settles(t, tt, "second sheet")
	if !tt.HasText("Account body") {
		t.Fatalf("second sheet %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if depth != 1 || changed != 1 {
		t.Fatalf("Escape: depth %d, changed %d", depth, changed)
	}
	settles(t, tt, "popped")
	click(t, tt, "Close")
	if depth != 0 || changed != 2 {
		t.Errorf("close: depth %d, changed %d", depth, changed)
	}
	cues(t, before, sound.Open, sound.Open, sound.Close, sound.Close)
	panics(t, "too deep", func() {
		render(func(c *ui.Context) { d := 4; SheetStack(c, titles, &d, func(int) {}, SheetStackOptions{}) })
	})
}
