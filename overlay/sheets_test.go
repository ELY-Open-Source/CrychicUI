package overlay

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestDrawer(t *testing.T) {
	open, changed := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Drawer(c, "Filters", &open, DrawerOptions{}, func() { ui.Text(c, "Price") }).Changed() {
			changed++
		}
	}, 800, 500)
	before := baseline()
	open = true
	tt.Frame()
	if !tt.HasText("Price") || !tt.FramePending() {
		t.Fatalf("did not slide in: %q", tt.Texts())
	}
	settles(t, tt, "Drawer")
	if r := find(t, tt, "Filters"); r.X < 400 {
		t.Errorf("not at the end: %+v", r)
	}
	tt.ClickAt(40, 250)
	tt.Frame()
	if open || changed != 1 {
		t.Fatalf("outside click: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Open, sound.Close)

	open = true
	tr := ui.NewTester(func(c *ui.Context) { Drawer(c, "Filters", &open, DrawerOptions{}, func() { ui.Text(c, "Price") }) }, 800, 500)
	tr.SetDirection(ui.RightToLeft)
	settles(t, tr, "right to left")
	if r := find(t, tr, "Filters"); r.X > 400 {
		t.Errorf("right to left: not at the start: %+v", r)
	}
}

func TestBottomSheet(t *testing.T) {
	open, snap, changed := true, 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if BottomSheet(c, "Details", &open, &snap, BottomSheetOptions{}, func() { ui.Text(c, "Notes") }).Changed() {
			changed++
		}
	}, 600, 700)
	settles(t, tt, "BottomSheet")
	low := find(t, tt, "Notes").Y
	before := baseline()
	tt.Key(0, ui.KeyUp)
	settles(t, tt, "Up")
	if snap != 1 || find(t, tt, "Notes").Y >= low {
		t.Fatalf("Up: snap %d, notes at %v (was %v)", snap, find(t, tt, "Notes").Y, low)
	}
	r := find(t, tt, "Details")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	for i := 1; i <= 8; i++ {
		tt.Advance(16 * time.Millisecond)
		tt.Move(r.X+r.W/2, r.Y+r.H/2+float32(i)*15)
	}
	tt.Advance(300 * time.Millisecond)
	tt.Release(r.X+r.W/2, r.Y+r.H/2+120)
	settles(t, tt, "dragged down")
	if snap != 0 || !open {
		t.Fatalf("drag: snap %d, open %v", snap, open)
	}
	r = find(t, tt, "Details")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	for i := 1; i <= 8; i++ {
		tt.Advance(16 * time.Millisecond)
		tt.Move(r.X+r.W/2, r.Y+r.H/2+float32(i)*40)
	}
	tt.Release(r.X+r.W/2, r.Y+r.H/2+320)
	settles(t, tt, "dragged away")
	if open || changed != 3 {
		t.Errorf("close: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Swipe, sound.Swipe, sound.Swipe)
	panics(t, "falling heights", func() {
		render(func(c *ui.Context) {
			o, s := true, 0
			BottomSheet(c, "x", &o, &s, BottomSheetOptions{Heights: []float32{300, 100}}, func() {})
		})
	})
}

func TestTopSheet(t *testing.T) {
	open := false
	tt := ui.NewTester(func(c *ui.Context) { TopSheet(c, "Search", &open, TopSheetOptions{}, func() { ui.Text(c, "Query") }) }, 600, 500)
	before := baseline()
	open = true
	tt.Frame()
	if !tt.HasText("Query") || !tt.FramePending() {
		t.Fatalf("did not drop: %q", tt.Texts())
	}
	settles(t, tt, "TopSheet")
	if r := find(t, tt, "Search"); r.Y > 100 {
		t.Errorf("not at the top: %+v", r)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open {
		t.Error("Escape did not close")
	}
	cues(t, before, sound.Open)
}

func TestInspector(t *testing.T) {
	open, changed := false, 0
	props := []Prop{{"Width", "320"}, {"Height", "240"}}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Text(c, "Canvas")
		if Inspector(c, "Layer", props, &open, InspectorOptions{}).Changed() {
			changed++
		}
	}, 800, 500)
	before := baseline()
	open = true
	tt.Frame()
	if !tt.HasText("240") || !tt.FramePending() {
		t.Fatalf("did not slide in: %q", tt.Texts())
	}
	settles(t, tt, "Inspector")
	if r := find(t, tt, "Layer"); r.X < 400 {
		t.Errorf("not on the right: %+v", r)
	}
	click(t, tt, "Close")
	if open || changed != 1 {
		t.Errorf("close: open %v, changed %d", open, changed)
	}
	cues(t, before, sound.Open)
	panics(t, "nameless prop", func() {
		render(func(c *ui.Context) { o := true; Inspector(c, "x", []Prop{{"", "1"}}, &o, InspectorOptions{}) })
	})
}
