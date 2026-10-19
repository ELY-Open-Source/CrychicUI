package layout

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestSplitPaneDrags(t *testing.T) {
	size, changes := float32(100), 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SplitPane(c, &size, func() { ui.Text(c, "Left") }, func() { ui.Text(c, "Right") }).Changed() {
			changes++
		}
	}, 400, 200)
	ticks := count(sound.Tick)
	tt.Press(104, 100)
	tt.Move(130, 100)
	tt.Move(154, 100)
	tt.Release(154, 100)
	if size != 150 || changes == 0 || count(sound.Tick) != ticks+1 {
		t.Errorf("size %v, changes %d, ticks %d", size, changes, count(sound.Tick)-ticks)
	}
	if x := find(t, tt, "Right").X; x < 150 {
		t.Errorf("right pane at %v", x)
	}
	tt.Press(154, 100)
	tt.Move(100, 100)
	tt.Move(-50, 100)
	tt.Release(-50, 100)
	if size != minPane {
		t.Errorf("size %v, want clamp to %d", size, minPane)
	}
	panics(t, "nil size", func() { ui.Render(func(c *ui.Context) { SplitPane(c, nil, nil, nil) }, 10, 10, 1) })
}

func TestSplitStackDrags(t *testing.T) {
	size := float32(80)
	tt := ui.NewTester(func(c *ui.Context) {
		SplitStack(c, &size, func() { ui.Text(c, "Top") }, func() { ui.Text(c, "Bottom") })
	}, 300, 300)
	ticks := count(sound.Tick)
	tt.Press(150, 84)
	tt.Move(150, 120)
	tt.Move(150, 144)
	tt.Release(150, 144)
	if size != 140 || count(sound.Tick) != ticks+1 {
		t.Errorf("size %v, ticks %d", size, count(sound.Tick)-ticks)
	}
	tt.Press(150, 144)
	tt.Move(150, 400)
	tt.Move(150, 600)
	tt.Release(150, 600)
	if want := 300 - 8 - float32(minPane); size != want {
		t.Errorf("size %v, want clamp to %v", size, want)
	}
}

func TestSidebarLayoutToggles(t *testing.T) {
	open, changes := true, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SidebarLayout(c, &open, func() { ui.Text(c, "Side") }, func() { ui.Text(c, "Main") }).Changed() {
			changes++
		}
	}, 600, 300)
	before, toggles := find(t, tt, "Main").X, count(sound.Toggle)
	if err := tt.Click("Toggle sidebar"); err != nil {
		t.Fatal(err)
	}
	if open || changes != 1 || count(sound.Toggle) != toggles+1 {
		t.Errorf("open %v, changes %d, toggles %d", open, changes, count(sound.Toggle)-toggles)
	}
	tt.Advance(80 * time.Millisecond)
	mid := find(t, tt, "Main").X
	settles(t, tt, "SidebarLayout width")
	if after := find(t, tt, "Main").X; !(before > mid && mid > after) {
		t.Errorf("main did not slide: %v, %v, %v", before, mid, after)
	}
	panics(t, "nil open", func() { ui.Render(func(c *ui.Context) { SidebarLayout(c, nil, nil, nil) }, 10, 10, 1) })
}

func TestSplitPaneStillPressIsNoChange(t *testing.T) {
	size, changes := float32(100), 0
	tt := ui.NewTester(func(c *ui.Context) {
		if SplitPane(c, &size, func() { ui.Text(c, "Left") }, func() { ui.Text(c, "Right") }).Changed() {
			changes++
		}
	}, 400, 200)
	for range 3 {
		tt.Press(104, 100)
		tt.Frame()
		tt.Release(104, 100)
	}
	if size != 100 || changes != 0 {
		t.Errorf("still presses: size %v, %d changes", size, changes)
	}
}

func TestSplitPaneNarrowKeepsMinimum(t *testing.T) {
	size := float32(minPane)
	tt := ui.NewTester(func(c *ui.Context) {
		SplitPane(c, &size, func() {}, func() {})
	}, 60, 100)
	tt.Press(44, 50)
	tt.Move(50, 50)
	tt.Release(50, 50)
	if size != minPane {
		t.Errorf("size %v, want %d", size, minPane)
	}
}

func TestSplitStackKeepsCallerKey(t *testing.T) {
	order, ids := []string{"a", "b"}, map[string]uint64{}
	sizes := map[string]*float32{"a": new(float32(60)), "b": new(float32(60))}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Children(func() {
			for _, k := range order {
				c.Key(k)
				ids[k] = SplitStack(c, sizes[k], func() {}, func() {}).ID()
			}
		})
	}, 300, 300)
	a, b := ids["a"], ids["b"]
	order = []string{"b", "a"}
	tt.Frame()
	if ids["a"] != a || ids["b"] != b {
		t.Errorf("ids moved with order: a %d→%d, b %d→%d", a, ids["a"], b, ids["b"])
	}
}
