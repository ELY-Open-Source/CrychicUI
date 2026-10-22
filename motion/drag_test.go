package motion

import (
	"fmt"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// drag presses, then moves by dx, dy in eight steps.
func drag(tt *ui.Tester, x, y, dx, dy float32) {
	tt.Press(x, y)
	for i := 1; i <= 8; i++ {
		tt.Advance(16 * time.Millisecond)
		tt.Move(x+dx*float32(i)/8, y+dy*float32(i)/8)
	}
	tt.Release(x+dx, y+dy)
}

// spot is the mean position of pixels exactly k.
func spot(tt *ui.Tester, k ui.Color) (float32, float32) {
	img, n := tt.Image(), float32(0)
	var sx, sy float32
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if p := img.RGBAAt(x, y); p.R == k.R && p.G == k.G && p.B == k.B {
				sx, sy, n = sx+float32(x), sy+float32(y), n+1
			}
		}
	}
	return sx / max(n, 1), sy / max(n, 1)
}

func TestCardStack(t *testing.T) {
	top, changes := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CardStack(c, 3, &top, func(i int) { ui.Box(c).Size(160, 80).Label(fmt.Sprint("Card ", i)) }, CardStackOptions{}).Changed() {
			changes++
		}
	}), 400, 300)
	r := find(t, tt, "Card 0")
	drag(tt, r.X+80, r.Y+40, 20, 0)
	tt.Advance(2 * time.Second)
	if top != 0 {
		t.Fatal("a short drag dismissed the card")
	}
	before := baseline()
	drag(tt, r.X+80, r.Y+40, 220, 0)
	tt.Advance(2 * time.Second)
	if top != 1 || changes != 1 {
		t.Fatalf("top %d, changes %d", top, changes)
	}
	cues(t, before, sound.Swipe)
	find(t, tt, "Card 1")
	top = 3
	tt.Frame()
	find(t, tt, "No more cards")
}

func TestSwipeActions(t *testing.T) {
	chosen := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SwipeActions(c, []string{"Archive", "Delete"}, &chosen, func() {
			ui.Box(c).Size(300, 48).Background(c.Theme().Accent).Label("row")
		})
	}), 400, 200)
	th := ui.LightTheme()
	r := find(t, tt, "row")
	// Poses skip Find; the row's right end shows the drag.
	covered := func() bool { return tinted(tt, int(r.X+290), int(r.Y+24), th.Accent) }
	before := baseline()
	drag(tt, r.X+250, r.Y+24, -200, 0)
	tt.Advance(2 * time.Second)
	if covered() {
		t.Fatal("row did not open")
	}
	cues(t, before, sound.Swipe)
	click(t, tt, "Delete")
	tt.Advance(2 * time.Second)
	if chosen != 1 || !covered() {
		t.Errorf("chosen %d, closed %v", chosen, covered())
	}
	cues(t, before, sound.Swipe, sound.Tap)
}

func TestThrowable(t *testing.T) {
	slot, changes := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Throwable(c, "Tile", &slot, ThrowableOptions{}).Changed() {
			changes++
		}
	}), 400, 200)
	r := find(t, tt, "Tile")
	before := baseline()
	drag(tt, r.X+30, r.Y+30, 150, 0)
	tt.Advance(2 * time.Second)
	if slot != 2 || changes != 1 {
		t.Fatalf("slot %d, changes %d", slot, changes)
	}
	cues(t, before, sound.Drop)
	tt.Key(0, ui.KeyLeft)
	tt.Advance(2 * time.Second)
	if slot != 1 {
		t.Errorf("left arrow left slot %d", slot)
	}
	slot = 0
	tt.Frame()
	tt.Advance(2 * time.Second)
	cues(t, before, sound.Drop, sound.Drop)
	if !tinted(tt, int(r.X+4), int(r.Y+4), ui.LightTheme().Accent) {
		t.Error("app slot 0 did not bring the tile back")
	}
	panics(t, "slot past slots", func() {
		render(func(c *ui.Context) { s := 3; Throwable(c, "x", &s, ThrowableOptions{}) })
	})
}

// spaced builds view under a theme with the given spacing.
func spaced(spacing *float32, view func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) {
		th := *c.Theme()
		th.Spacing = *spacing
		c.SetTheme(&th)
		ui.Column(c).Padding(40).Children(func() { view(c) })
	}
}

func TestDragsFollowSpacing(t *testing.T) {
	spacing, slot, chosen := float32(4), 2, -1
	tt := ui.NewTester(spaced(&spacing, func(c *ui.Context) {
		Throwable(c, "Tile", &slot, ThrowableOptions{})
		SwipeActions(c, []string{"A", "B"}, &chosen, func() { ui.Box(c).Size(300, 48).Background(c.Theme().Accent).Label("row") })
	}), 600, 300)
	th := ui.LightTheme()
	r := find(t, tt, "row")
	drag(tt, r.X+250, r.Y+24, -200, 0)
	tt.Advance(2 * time.Second)
	before := baseline()
	spacing = 5
	tt.Frame()
	tt.Advance(2 * time.Second)
	if slot != 2 || !tinted(tt, 40+180+76, 40+40, th.Accent) {
		t.Errorf("tile: slot %d, not at the new third slot", slot)
	}
	if r = find(t, tt, "row"); tinted(tt, int(r.X+120), int(r.Y+24), th.Accent) {
		t.Error("open row did not follow the wider actions")
	}
	cues(t, before)
}

func TestCardStackSpacingMidFlight(t *testing.T) {
	spacing, top := float32(4), 0
	tt := ui.NewTester(spaced(&spacing, func(c *ui.Context) {
		CardStack(c, 2, &top, func(i int) { ui.Box(c).Size(160, 80).Label(fmt.Sprint("Card ", i)) }, CardStackOptions{})
	}), 400, 300)
	r := find(t, tt, "Card 0")
	before := baseline()
	drag(tt, r.X+80, r.Y+40, 220, 0)
	tt.Advance(50 * time.Millisecond)
	spacing = 5
	tt.Frame()
	tt.Advance(2 * time.Second)
	if top != 1 {
		t.Errorf("spacing change mid-flight left top %d", top)
	}
	cues(t, before, sound.Swipe)
}

func TestThrowableDirectionChange(t *testing.T) {
	slot, changes := 2, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Throwable(c, "Tile", &slot, ThrowableOptions{}).Changed() {
			changes++
		}
	}), 400, 200)
	before := baseline()
	tt.SetDirection(ui.RightToLeft)
	tt.Frame()
	tt.Advance(2 * time.Second)
	if slot != 2 || changes != 0 {
		t.Errorf("direction change moved slot to %d with %d changes", slot, changes)
	}
	cues(t, before)
}

// catch throws by dx, then holds it mid-flight.
func catch(tt *ui.Tester, x, y, dx float32) {
	drag(tt, x, y, dx, 0)
	tt.Advance(16 * time.Millisecond)
	tt.Press(x+dx, y)
	tt.Advance(time.Second)
	tt.Move(x+dx+1, y)
	tt.Advance(100 * time.Millisecond)
}

func TestCaughtDragsWait(t *testing.T) {
	top, chosen, slot, changes := 0, -1, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ui.Column(c).Gap(40).Children(func() {
			CardStack(c, 2, &top, func(i int) { ui.Box(c).Size(160, 80).Label(fmt.Sprint("Card ", i)) }, CardStackOptions{})
			SwipeActions(c, []string{"A", "B"}, &chosen, func() { ui.Box(c).Size(300, 48).Label("row") })
			if Throwable(c, "Tile", &slot, ThrowableOptions{}).Changed() {
				changes++
			}
		})
	}), 500, 500)
	card, row, tile := find(t, tt, "Card 0"), find(t, tt, "row"), find(t, tt, "Tile")
	before := baseline()
	for _, at := range []struct{ x, y, dx float32 }{
		{card.X + 80, card.Y + 40, 20}, {row.X + 250, row.Y + 24, -20}, {tile.X + 30, tile.Y + 30, 110},
	} {
		catch(tt, at.x, at.y, at.dx)
		if top != 0 || slot != 0 || changes != 0 {
			t.Errorf("held at %v: top %d, slot %d, changes %d", at, top, slot, changes)
		}
		cues(t, before)
		tt.Release(at.x+at.dx+1, at.y)
		tt.Advance(2 * time.Second)
	}
	if top != 0 || slot == 0 || changes != 1 {
		t.Errorf("after release: top %d, slot %d, changes %d", top, slot, changes)
	}
	cues(t, before, sound.Drop)
}
func TestPullRefresh(t *testing.T) {
	busy := false
	tt := ui.NewTester(func(c *ui.Context) {
		PullRefresh(c, &busy, PullRefreshOptions{}, func() { ui.Box(c).Size(300, 200).Label("feed") })
	}, 300, 300)
	r := find(t, tt, "feed")
	drag(tt, 150, 100, 0, 30)
	tt.Advance(time.Second)
	if busy || abs(find(t, tt, "feed").Y-r.Y) > 0.5 {
		t.Fatal("a short pull refreshed or did not spring back")
	}
	before := baseline()
	tt.Press(150, 100)
	for i := 1; i <= 8; i++ {
		tt.Advance(16 * time.Millisecond)
		tt.Move(150, 100+float32(i)*24)
	}
	if !tt.HasText("Release to refresh") {
		t.Errorf("past the threshold shows %q", tt.Texts())
	}
	tt.Release(150, 292)
	tt.Advance(time.Second)
	if !busy || !tt.HasText("Refreshing…") || find(t, tt, "feed").Y <= r.Y {
		t.Fatalf("busy %v, texts %q", busy, tt.Texts())
	}
	cues(t, before, sound.Pop)
	busy = false
	tt.Frame()
	tt.Advance(time.Second)
	if abs(find(t, tt, "feed").Y-r.Y) > 0.5 {
		t.Error("did not spring back when done")
	}
}

func TestMagnet(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Magnet(c, func() { dot(c) }) }), 300, 300)
	r := find(t, tt, "dot")
	tt.Move(r.X-20, r.Y-20)
	tt.Frame()
	tt.Advance(time.Second)
	if p := find(t, tt, "dot"); p.X >= r.X || p.Y >= r.Y {
		t.Errorf("not pulled toward the pointer: %v from %v", p, r)
	}
	tt.Move(290, 290)
	tt.Frame()
	settles(t, tt, "spring back")
	if p := find(t, tt, "dot"); abs(p.X-r.X) > 0.5 || abs(p.Y-r.Y) > 0.5 {
		t.Errorf("rests at %v, want %v", p, r)
	}
}

func TestFollower(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Follower(c, FollowerOptions{}) }, 300, 200)
	th := ui.LightTheme()
	tt.Advance(time.Second)
	cx, cy := spot(tt, th.Accent)
	if cx < 140 || cx > 160 || cy < 70 || cy > 90 {
		t.Errorf("idle dot at %v,%v, want the center", cx, cy)
	}
	tt.Move(50, 40)
	tt.Frame()
	settles(t, tt, "follow")
	if x, y := spot(tt, th.Accent); x < 45 || x > 55 || y < 35 || y > 45 {
		t.Errorf("dot at %v,%v, want 50,40", x, y)
	}
}
