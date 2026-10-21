package button

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestLikeButton(t *testing.T) {
	liked, likes, changed := false, 4, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if LikeButton(c, &liked, &likes, LikeButtonOptions{}).Changed() {
			changed++
		}
	}, 300, 200)
	pops := sound.Recent()
	if err := tt.Click("Like"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "LikeButton burst")
	if !liked || likes != 5 || changed != 1 || !appended(pops, sound.Pop) {
		t.Errorf("liked %v, likes %d, changed %d", liked, likes, changed)
	}
	if err := tt.Click("Like"); err != nil {
		t.Fatal(err)
	}
	if liked || likes != 4 || changed != 2 {
		t.Errorf("unlike failed: liked %v, likes %d", liked, likes)
	}
}

func TestStarButton(t *testing.T) {
	on, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if StarButton(c, &on, StarButtonOptions{}).Changed() {
			changed++
		}
	}, 300, 200)
	taps := sound.Recent()
	if err := tt.Click("Star"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "StarButton spin")
	if !on || changed != 1 || !appended(taps, sound.Toggle) {
		t.Errorf("on %v, changed %d", on, changed)
	}
	panics(t, "StarButton state", func() { ui.Render(func(c *ui.Context) { StarButton(c, nil, StarButtonOptions{}) }, 10, 10, 1) })
}

func TestFollowButton(t *testing.T) {
	following, changed := false, 0
	width := float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		b := FollowButton(c, &following, FollowButtonOptions{})
		width = b.LayoutBounds().W
		if b.Changed() {
			changed++
		}
	}, 300, 200)
	tt.Frame()
	idle := width
	taps := sound.Recent()
	if err := tt.Click("Follow"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "FollowButton")
	tt.Frame()
	if width == idle {
		t.Error("label change did not change width")
	}
	if !following || changed != 1 || !appended(taps, sound.Toggle) || !tt.HasText("Following") {
		t.Errorf("following %v, changed %d", following, changed)
	}
	if err := tt.Click("Following"); err != nil {
		t.Fatal(err)
	}
	if following || changed != 2 {
		t.Error("unfollow failed")
	}
}

func TestCounterButton(t *testing.T) {
	n, changed := 2, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if CounterButton(c, "Notifications", &n).Changed() {
			changed++
		}
	}, 300, 200)
	ticks := sound.Recent()
	if err := tt.Click("Notifications"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "CounterButton")
	if n != 3 || changed != 1 || !appended(ticks, sound.Tick) {
		t.Errorf("count %d, changed %d", n, changed)
	}
}

func TestUndoButton(t *testing.T) {
	t.Run("commit", func(t *testing.T) {
		changed := 0
		tt := ui.NewTester(func(c *ui.Context) {
			if UndoButton(c, "Undo", UndoButtonOptions{Duration: 100 * time.Millisecond}).Changed() {
				changed++
			}
		}, 300, 200)
		closes := sound.Recent()
		tt.Advance(200 * time.Millisecond)
		if changed != 1 || !appended(closes) || tt.FramePending() {
			t.Errorf("changed %d, close cues %v, pending %v", changed, sound.Recent(), tt.FramePending())
		}
	})
	t.Run("undo", func(t *testing.T) {
		changed, clicks := 0, 0
		tt := ui.NewTester(func(c *ui.Context) {
			b := UndoButton(c, "Undo", UndoButtonOptions{Duration: 3 * time.Second})
			if b.Clicked() {
				clicks++
			}
			if b.Changed() {
				changed++
			}
		}, 300, 200)
		closes := sound.Recent()
		if err := tt.Click("Undo"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(3 * time.Second)
		if changed != 0 || clicks != 1 || !appended(closes, sound.Close) {
			t.Errorf("changed %d, clicks %d, close cues %v", changed, clicks, sound.Recent())
		}
	})
}

func TestUndoButtonReducedMotion(t *testing.T) {
	changed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if UndoButton(c, "Undo", UndoButtonOptions{Duration: 100 * time.Millisecond}).Changed() {
			changed++
		}
	}, 300, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	if changed != 0 {
		t.Fatal("undo committed immediately")
	}
	tt.Advance(50 * time.Millisecond)
	if changed != 0 {
		t.Fatal("undo committed before its duration")
	}
	tt.Advance(60 * time.Millisecond)
	if changed != 1 {
		t.Error("undo did not commit at its duration")
	}
}

func TestSpeedDial(t *testing.T) {
	open, selected, changed := false, -1, 0
	actions := []SpeedDialAction{{Label: "Edit", Icon: "pencil"}, {Label: "Share", Icon: "share-2"}}
	tt := ui.NewTester(func(c *ui.Context) {
		if SpeedDial(c, &open, actions, &selected, SpeedDialOptions{}).Changed() {
			changed++
		}
	}, 300, 300)
	opens, closes, taps := sound.Recent(), sound.Recent(), sound.Recent()
	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "SpeedDial open")
	if !open || changed != 1 || !appended(opens, sound.Open) {
		t.Error("dial did not open")
	}
	if err := tt.Click("Edit"); err != nil {
		t.Fatal(err)
	}
	if open || selected != 0 || changed != 2 || !appended(taps, sound.Tap) {
		t.Errorf("open %v, selected %d, changed %d", open, selected, changed)
	}
	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	if open || !appended(closes, sound.Close) {
		t.Error("dial did not close")
	}
}

func TestSpeedDialKeepsFocus(t *testing.T) {
	open, selected := false, -1
	actions := []SpeedDialAction{{Label: "Edit", Icon: "pencil"}}
	tt := ui.NewTester(func(c *ui.Context) {
		SpeedDial(c, &open, actions, &selected, SpeedDialOptions{})
	}, 300, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !open || !tt.Focused("Actions") {
		t.Errorf("open %v, focused %v", open, tt.Focused("Actions"))
	}
	tt.Key(0, ui.KeyEnter)
	if open {
		t.Error("second Enter did not close the dial")
	}
}

func TestSplitButton(t *testing.T) {
	selected, changed, mainClicks := -1, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		root, main := SplitButton(c, "Publish", []string{"Draft", "Schedule"}, &selected)
		if root.Changed() {
			changed++
		}
		if main.Clicked() {
			mainClicks++
		}
	}, 300, 200)
	taps, opens := sound.Recent(), sound.Recent()
	if err := tt.Click("Publish"); err != nil {
		t.Fatal(err)
	}
	if mainClicks != 1 || !appended(taps, sound.Tap) || selected != -1 {
		t.Error("main action failed")
	}
	if err := tt.Click("Publish menu"); err != nil {
		t.Fatal(err)
	}
	if len(tt.Menu()) != 2 || !appended(opens, sound.Open) {
		t.Errorf("menu %q", tt.Menu())
	}
	if err := tt.ChooseMenuItem("Draft"); err != nil {
		t.Fatal(err)
	}
	if selected != 0 || changed != 1 {
		t.Errorf("selected %d, changed %d", selected, changed)
	}
}
