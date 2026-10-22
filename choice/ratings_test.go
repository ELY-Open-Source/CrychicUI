package choice

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestKnob(t *testing.T) {
	v, changes := 50.0, 0
	var r ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			k := Knob(c, "Gain", &v, 0, 100, KnobOptions{Step: 5})
			if k.Changed() {
				changes++
			}
			r = k.LayoutBounds()
		})
	}, 200, 200)
	tt.Frame()
	find(t, tt, "50")
	before := baseline()
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W/2, r.Y+r.H/2-r.W/2)
	if v != 75 {
		t.Fatalf("drag up half a size gave %v, want 75", v)
	}
	tt.Release(r.X+r.W/2, r.Y+r.H/2-r.W/2)
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyDown)
	if v != 95 {
		t.Errorf("End then Down gave %v", v)
	}
	if len(played(before)) != changes || changes < 3 {
		t.Errorf("changes %d, cues %v", changes, played(before))
	}
	for _, cue := range played(before) {
		if cue != sound.Tick {
			t.Errorf("cue %v, want only ticks", cue)
		}
	}
	settles(t, tt, "Knob arc")
	find(t, tt, "95")
}

func TestStarRating(t *testing.T) {
	v := 0
	var r ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			r = StarRating(c, "Rating", &v, StarRatingOptions{}).LayoutBounds()
		})
	}, 300, 120)
	tt.Frame()
	star := r.W / 5
	before := baseline()
	tt.ClickAt(r.X+star*3.5, r.Y+r.H/2)
	if v != 4 {
		t.Fatalf("click on the fourth star gave %d", v)
	}
	settles(t, tt, "StarRating bounce")
	tt.Key(0, ui.KeyLeft)
	tt.ClickAt(r.X+star*2.5, r.Y+r.H/2)
	if v != 0 {
		t.Errorf("click on the set star gave %d, want 0", v)
	}
	cues(t, before, sound.Select, sound.Select, sound.Select)
}

func TestEmojiRating(t *testing.T) {
	v := 0
	tt := ui.NewTester(func(c *ui.Context) { EmojiRating(c, "Mood", &v, EmojiRatingOptions{}) }, 300, 120)
	before := baseline()
	click(t, tt, "Good")
	if v != 4 {
		t.Fatalf("value %d", v)
	}
	settles(t, tt, "EmojiRating face")
	tt.Key(0, ui.KeyRight)
	if v != 5 {
		t.Errorf("arrow gave %d", v)
	}
	cues(t, before, sound.Select, sound.Select)
}

func TestNPS(t *testing.T) {
	score := -1
	tt := ui.NewTester(func(c *ui.Context) { NPS(c, "Recommend", &score, NPSOptions{High: "Surely"}) }, 500, 120)
	find(t, tt, "Not likely")
	find(t, tt, "Surely")
	before := baseline()
	click(t, tt, "9")
	if score != 9 {
		t.Fatalf("score %d", score)
	}
	cues(t, before, sound.Select)
	settles(t, tt, "NPS cell")
	panics(t, "score 11", func() {
		s := 11
		ui.Render(func(c *ui.Context) { NPS(c, "x", &s, NPSOptions{}) }, 50, 50, 1)
	})
}

func TestLikert(t *testing.T) {
	v := -1
	tt := ui.NewTester(func(c *ui.Context) { Likert(c, "The docs are clear.", &v, LikertOptions{}) }, 600, 160)
	before := baseline()
	click(t, tt, "Agree")
	if v != 3 {
		t.Fatalf("value %d", v)
	}
	cues(t, before, sound.Select)
	settles(t, tt, "Likert dot")
}

func TestThumbs(t *testing.T) {
	v := 0
	tt := ui.NewTester(func(c *ui.Context) { Thumbs(c, "Answer", &v, ThumbsOptions{}) }, 200, 80)
	before := baseline()
	click(t, tt, "Helpful")
	if v != 1 {
		t.Fatalf("value %d", v)
	}
	settles(t, tt, "Thumbs bounce")
	click(t, tt, "Not helpful")
	click(t, tt, "Not helpful")
	if v != 0 {
		t.Errorf("second click gave %d, want 0", v)
	}
	cues(t, before, sound.Select, sound.Select, sound.Select)
}

func TestNumberStepper(t *testing.T) {
	v := 0.8
	tt := ui.NewTester(func(c *ui.Context) {
		NumberStepper(c, "Opacity", &v, 0, 1, 0.1, NumberStepperOptions{Format: "%.1f"})
	}, 240, 80)
	before := baseline()
	click(t, tt, "Increase Opacity")
	click(t, tt, "Increase Opacity")
	if v != 1 {
		t.Fatalf("value %v, want 1", v)
	}
	click(t, tt, "Decrease Opacity")
	if v != 0.9 {
		t.Errorf("value %v, want 0.9", v)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	settles(t, tt, "NumberStepper digits")
	find(t, tt, "0.9")
}

func TestWheel(t *testing.T) {
	sel, changes := 0, 0
	spacing := float32(4)
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	var r ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		th := *c.Theme()
		th.Spacing = spacing
		c.SetTheme(&th)
		ui.Column(c).Padding(20).Children(func() {
			w := Wheel(c, "Month", months, &sel)
			if w.Changed() {
				changes++
			}
			r = w.LayoutBounds()
		})
	}, 300, 300)
	tt.Frame()
	before := baseline()
	row := r.H / 5
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	for i := 1; i <= 6; i++ {
		tt.Move(x, y-row*float32(i)/2)
		tt.Advance(16 * time.Millisecond)
	}
	tt.Advance(300 * time.Millisecond)
	tt.Release(x, y-row*3)
	tt.Advance(2 * time.Second)
	if sel != 3 || changes != 1 {
		t.Fatalf("drag three rows gave %d, changes %d", sel, changes)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	tt.Key(0, ui.KeyDown)
	tt.Advance(2 * time.Second)
	if sel != 4 {
		t.Errorf("Down gave %d", sel)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	tt.Scroll(x, y, 0, row*2)
	tt.Advance(2 * time.Second)
	if sel != 5 {
		t.Errorf("scrolling two rows from 4 gave %d, want 5 at the end", sel)
	}
	if n := len(played(before)); n != 5 {
		t.Errorf("scroll played %d cues in all, want 5", n)
	}
	spacing = 5
	tt.Frame()
	tt.Advance(2 * time.Second)
	if sel != 5 || changes != 3 {
		t.Errorf("roomier theme moved the wheel to %d, changes %d", sel, changes)
	}
	sel = 0
	tt.Frame()
	tt.Advance(2 * time.Second)
	if len(played(before)) != 5 {
		t.Errorf("app change played cues %v", played(before))
	}
}
