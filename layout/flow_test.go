package layout

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCollapse(t *testing.T) {
	open := true
	height := float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		height = Collapse(c, open, func() { ui.Box(c).Label("Body").Size(100, 100) }).LayoutBounds().H
	}, 300, 200)
	tt.Frame()
	if height != 100 {
		t.Fatalf("open height %v", height)
	}
	open = false
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	if height <= 0 || height >= 100 {
		t.Errorf("closing height jumped to %v", height)
	}
	settles(t, tt, "Collapse exit")
	if tt.HasText("Body") {
		t.Error("closed content remains")
	}
	open = true
	tt.Frame()
	tt.Advance(50 * time.Millisecond)
	if height <= 0 || height >= 100 {
		t.Errorf("opening height jumped to %v", height)
	}
	settles(t, tt, "Collapse enter")
	find(t, tt, "Body")
}

func TestReveal(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Reveal(c, func() { ui.Text(c, "Revealed") })
	}, 300, 200)
	settles(t, tt, "Reveal")
	find(t, tt, "Revealed")
}

func TestDivider(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Children(func() {
			Divider(c, DividerOptions{Label: "or"})
			Divider(c, DividerOptions{})
		})
	}, 300, 200)
	settles(t, tt, "Divider growth")
	if !tt.HasText("or") {
		t.Error("label missing")
	}
}

func TestVDivider(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Height(30).Children(func() {
			ui.Text(c, "Left")
			VDivider(c)
			ui.Text(c, "Right")
		})
	}, 300, 200)
	settles(t, tt, "VDivider growth")
	if find(t, tt, "Right").X <= find(t, tt, "Left").X {
		t.Error("divider disturbed the row")
	}
}

func TestGutter(t *testing.T) {
	n := float32(2)
	var gap float32
	tt := ui.NewTester(func(c *ui.Context) {
		gap = c.Theme().Space(n)
		ui.Row(c).Children(func() {
			ui.Text(c, "Left")
			Gutter(c, n)
			ui.Text(c, "Right")
		})
	}, 300, 200)
	check := func() {
		t.Helper()
		l, r := find(t, tt, "Left"), find(t, tt, "Right")
		if got := r.X - (l.X + l.W); abs(got-gap) > 0.5 {
			t.Errorf("gap %v, want %v", got, gap)
		}
	}
	check()
	n = 6
	tt.Frame()
	settles(t, tt, "Gutter resize")
	check()
}

func TestBar(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Bar(c, func() { ui.Text(c, "Start") }, func() { ui.Text(c, "Center") }, func() { ui.Text(c, "End") })
	}, 400, 200)
	settles(t, tt, "Bar slots")
	start, center, end := find(t, tt, "Start"), find(t, tt, "Center"), find(t, tt, "End")
	if !(start.X < center.X && center.X < end.X) {
		t.Errorf("slot order %v %v %v", start, center, end)
	}
	if got := center.X + center.W/2; abs(got-200) > 1 {
		t.Errorf("center at %v", got)
	}
}

func TestPageHeader(t *testing.T) {
	clicked := false
	tt := ui.NewTester(func(c *ui.Context) {
		PageHeader(c, "Projects", PageHeaderOptions{
			Subtitle: "All projects",
			Actions:  func() { clicked = ui.Button(c, "New").Clicked() || clicked },
		})
	}, 500, 200)
	settles(t, tt, "PageHeader")
	find(t, tt, "Projects")
	find(t, tt, "All projects")
	if err := tt.Click("New"); err != nil {
		t.Fatal(err)
	}
	if !clicked {
		t.Error("action slot did not click")
	}
}

func TestPageFooter(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		PageFooter(c, func() { ui.Text(c, "Terms") })
	}, 300, 200)
	settles(t, tt, "PageFooter")
	find(t, tt, "Terms")
}

func TestPane(t *testing.T) {
	open, changed := true, 0
	height := float32(0)
	tt := ui.NewTester(func(c *ui.Context) {
		root := Pane(c, "Properties", &open, func() { ui.Box(c).Label("Width").Size(100, 100) })
		height = root.LayoutBounds().H
		if root.Changed() {
			changed++
		}
	}, 300, 200)
	tt.Frame()
	openHeight := height
	header := find(t, tt, "Properties")
	closedHeight := header.Y + header.H
	taps := count(sound.Toggle)
	if err := tt.Click("Properties"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(50 * time.Millisecond)
	if height <= closedHeight || height >= openHeight {
		t.Errorf("closing pane height jumped to %v, want %v to %v", height, closedHeight, openHeight)
	}
	settles(t, tt, "Pane close")
	if open || changed != 1 || count(sound.Toggle) != taps+1 || tt.HasText("Width") {
		t.Errorf("open %v, changed %d, cues %d", open, changed, count(sound.Toggle)-taps)
	}
	if err := tt.Click("Properties"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(50 * time.Millisecond)
	if height <= closedHeight || height >= openHeight {
		t.Errorf("opening pane height jumped to %v, want %v to %v", height, closedHeight, openHeight)
	}
	settles(t, tt, "Pane open")
	if !open || changed != 2 || count(sound.Toggle) != taps+2 {
		t.Errorf("reopen: open %v, changed %d", open, changed)
	}
	find(t, tt, "Width")
}
