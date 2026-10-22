package overlay

import (
	"slices"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// pressable is a button counting its clicks into hits.
func pressable(c *ui.Context, label string, hits *int) {
	if ui.ButtonBase(c).Label(label).Size(120, 40).Clicked() {
		*hits++
	}
}

func TestBusyVeil(t *testing.T) {
	busy, hits := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		BusyVeil(c, busy, "Saving", func() { pressable(c, "Edit", &hits) })
	}, 400, 300)
	click(t, tt, "Edit")
	if hits != 1 || tt.HasText("Saving") {
		t.Fatalf("idle: hits %d, %q", hits, tt.Texts())
	}
	busy = true
	tt.Frame()
	if !tt.HasText("Saving") || !tt.FramePending() {
		t.Fatalf("busy: no veil fading in: %q", tt.Texts())
	}
	click(t, tt, "Edit")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if hits != 1 {
		t.Errorf("the veil let a click or key through: hits %d", hits)
	}
	busy = false
	tt.Frame()
	if !tt.FramePending() {
		t.Error("the veil did not fade out")
	}
	settles(t, tt, "veil gone")
	if _, ok := tt.Find("Saving"); ok {
		t.Errorf("veil stayed: %q", tt.Texts())
	}
	panics(t, "no message", func() { render(func(c *ui.Context) { BusyVeil(c, true, "", func() {}) }) })
}

func TestDropOverlay(t *testing.T) {
	var files, got []string
	over, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if dropOverlay(c, ui.Box(c), &files, DropOverlayOptions{}, func() { ui.Text(c, "Inbox") }, over, got).Changed() {
			changed++
		}
		got = nil
	}, 400, 300)
	if tt.HasText("Drop files here") {
		t.Fatal("veil shown with nothing dragged")
	}
	over = true
	tt.Frame()
	if !tt.HasText("Drop files here") {
		t.Fatalf("drag over: %q", tt.Texts())
	}
	before := baseline()
	over, got = false, []string{"a.png", "b.txt"}
	tt.Frame()
	if !slices.Equal(files, []string{"a.png", "b.txt"}) || changed != 1 {
		t.Errorf("drop: files %q, changed %d", files, changed)
	}
	cues(t, before, sound.Drop)
	settles(t, tt, "veil gone")
	render(func(c *ui.Context) { DropOverlay(c, &files, DropOverlayOptions{}, func() {}) })
	panics(t, "no files", func() { render(func(c *ui.Context) { DropOverlay(c, nil, DropOverlayOptions{}, func() {}) }) })
}

func TestLockOverlay(t *testing.T) {
	locked, hits, changed := true, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if LockOverlay(c, "Vault", &locked, LockOverlayOptions{Text: "Unlock to see notes."}, func() { pressable(c, "Edit", &hits) }).Changed() {
			changed++
		}
	}, 400, 400)
	settles(t, tt, "locked")
	if !tt.HasText("Unlock to see notes.") {
		t.Fatalf("no lock panel: %q", tt.Texts())
	}
	r := find(t, tt, "Edit")
	tt.ClickAt(r.X+4, r.Y+4)
	tt.Frame()
	if hits != 0 {
		t.Errorf("locked content took a click")
	}
	before := baseline()
	// Tab skips the locked content to Unlock.
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if locked || changed != 1 {
		t.Fatalf("unlock: locked %v, changed %d", locked, changed)
	}
	cues(t, before, sound.Success)
	settles(t, tt, "unlocked")
	click(t, tt, "Edit")
	if hits != 1 {
		t.Errorf("unlocked content took no click")
	}
}

func TestPaywall(t *testing.T) {
	hits, upgrades := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if _, up := Paywall(c, "Go Pro", "Charts need Pro.", PaywallOptions{}, func() { pressable(c, "Chart", &hits) }); up.Clicked() {
			upgrades++
		}
	}, 400, 400)
	settles(t, tt, "paywall")
	r := find(t, tt, "Chart")
	tt.ClickAt(r.X+4, r.Y+4)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	before := baseline()
	click(t, tt, "Upgrade")
	if hits != 0 || upgrades != 2 {
		t.Errorf("hits %d, upgrades %d", hits, upgrades)
	}
	cues(t, before, sound.Open)
	panics(t, "no title", func() { render(func(c *ui.Context) { Paywall(c, "", "", PaywallOptions{}, func() {}) }) })
}

func TestEmojiPicker(t *testing.T) {
	groups := []EmojiGroup{
		{"Smileys", []Emoji{{"😀", "grinning"}, {"😂", "joy"}}},
		{"Animals", []Emoji{{"🐶", "dog"}, {"🐱", "cat"}}},
	}
	chosen, changed := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		if EmojiPicker(c, "Emoji", groups, &chosen, EmojiPickerOptions{}).Changed() {
			changed++
		}
	}, 400, 400)
	before := baseline()
	click(t, tt, "dog")
	if chosen != "🐶" || changed != 1 {
		t.Fatalf("pick: %q, changed %d", chosen, changed)
	}
	cues(t, before, sound.Select)
	click(t, tt, "Emoji")
	tt.Type("ca")
	tt.Frame()
	if _, ok := tt.Find("dog"); ok || !tt.HasText("Animals") {
		t.Errorf("search ca: %q", tt.Texts())
	}
	if _, ok := tt.Find("Smileys"); ok {
		t.Errorf("a group without matches stayed")
	}
	find(t, tt, "cat")
	tt.Type("zz")
	tt.Frame()
	if !tt.HasText("No emoji found") {
		t.Errorf("no matches: %q", tt.Texts())
	}
	panics(t, "twin names", func() {
		render(func(c *ui.Context) {
			s := ""
			EmojiPicker(c, "x", []EmojiGroup{{"A", []Emoji{{"🐶", "dog"}}}, {"B", []Emoji{{"🐕", "dog"}}}}, &s, EmojiPickerOptions{})
		})
	})
}
