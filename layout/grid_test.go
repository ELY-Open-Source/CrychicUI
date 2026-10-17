package layout

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func block(c *ui.Context, label string, h float32) func() {
	return func() { ui.Column(c).Height(h).Children(func() { ui.Text(c, label) }) }
}

func TestMasonryFillsShortest(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Masonry(c, 2, []func(){block(c, "I0", 120), block(c, "I1", 20), block(c, "I2", 20), block(c, "I3", 20)})
	}, 400, 400)
	tt.Advance(time2s)
	i0, i1, i2, i3 := find(t, tt, "I0"), find(t, tt, "I1"), find(t, tt, "I2"), find(t, tt, "I3")
	if i1.X == i0.X || i2.X != i1.X || i3.X != i1.X || !(i1.Y < i2.Y && i2.Y < i3.Y) {
		t.Errorf("placement %v %v %v %v", i0, i1, i2, i3)
	}
	panics(t, "0 columns", func() { ui.Render(func(c *ui.Context) { Masonry(c, 0, nil) }, 10, 10, 1) })
}

func TestAutoGridColumns(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		var items []func()
		for i := range 4 {
			items = append(items, block(c, fmt.Sprintf("G%d", i), 30))
		}
		AutoGrid(c, 120, items)
	}, 412, 300)
	tt.Advance(time2s)
	g0, g1, g2, g3 := find(t, tt, "G0"), find(t, tt, "G1"), find(t, tt, "G2"), find(t, tt, "G3")
	if !(g0.X < g1.X && g1.X < g2.X) || g0.Y != g2.Y || g3.X != g0.X || g3.Y <= g0.Y {
		t.Errorf("not 3 columns: %v %v %v %v", g0, g1, g2, g3)
	}
	tt.SetSize(260, 300)
	tt.Advance(time2s)
	if g1, g2 := find(t, tt, "G1"), find(t, tt, "G2"); g2.X != find(t, tt, "G0").X || g1.Y != find(t, tt, "G0").Y {
		t.Error("not 2 columns after resize")
	}
	panics(t, "minWidth 0", func() { ui.Render(func(c *ui.Context) { AutoGrid(c, 0, nil) }, 10, 10, 1) })
}

func TestBentoSpans(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Bento(c, 3, []BentoTile{
			{Cols: 2, Content: func() { ui.Text(c, "Wide") }},
			{Content: func() { ui.Text(c, "One") }},
			{Content: func() { ui.Text(c, "Two") }},
		})
	}, 600, 300)
	w, one, two := find(t, tt, "Wide"), find(t, tt, "One"), find(t, tt, "Two")
	if one.Y != w.Y || one.X < 380 || two.Y <= w.Y || two.X > w.X+1 {
		t.Errorf("spans wrong: %v %v %v", w, one, two)
	}
	tt.Move(500, 40)
	settles(t, tt, "Bento lift")
	panics(t, "span 4 of 3", func() {
		ui.Render(func(c *ui.Context) { Bento(c, 3, []BentoTile{{Cols: 4}}) }, 10, 10, 1)
	})
}

func TestClusterAligns(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Cluster(c, ui.End, []func(){func() { ui.Text(c, "Tail") }})
	}, 300, 100)
	if r := find(t, tt, "Tail"); abs(r.X+r.W-300) > 1 {
		t.Errorf("Tail ends at %v", r.X+r.W)
	}
	panics(t, "Stretch", func() { ui.Render(func(c *ui.Context) { Cluster(c, ui.Stretch, nil) }, 10, 10, 1) })
}

func TestColumnsEqual(t *testing.T) {
	var gap float32
	n := 2
	tt := ui.NewTester(func(c *ui.Context) {
		gap = c.Theme().Space(4)
		items := []func(){func() { ui.Text(c, "C0") }, func() { ui.Text(c, "C1") }, func() { ui.Text(c, "C2") }}
		Columns(c, items[:n])
	}, 400, 100)
	if x := find(t, tt, "C1").X; abs(x-(200+gap/2)) > 0.5 {
		t.Errorf("C1 at %v", x)
	}
	n = 3
	tt.Frame()
	settles(t, tt, "Columns width")
	if x := find(t, tt, "C1").X; abs(x-(400-2*gap)/3-gap) > 0.5 {
		t.Errorf("C1 at %v after a third column", x)
	}
	panics(t, "no items", func() { ui.Render(func(c *ui.Context) { Columns(c, nil) }, 10, 10, 1) })
}

func TestMasonryReplacesWhenAnItemShrinks(t *testing.T) {
	h0 := float32(120)
	tt := ui.NewTester(func(c *ui.Context) {
		Masonry(c, 2, []func(){block(c, "I0", h0), block(c, "I1", 20), block(c, "I2", 20)})
	}, 400, 400)
	tt.Advance(time2s)
	if find(t, tt, "I2").X == find(t, tt, "I0").X {
		t.Fatal("I2 starts under I0")
	}
	x2 := find(t, tt, "I2").X
	h0 = 10
	tt.Frame()
	if find(t, tt, "I2").X != x2 {
		t.Error("I2 jumped instead of sliding")
	}
	tt.Advance(time2s)
	if find(t, tt, "I2").X != find(t, tt, "I0").X {
		t.Error("I2 stayed beside I0 after I0 shrank")
	}
}

func TestMasonrySettlesOnAnEmptyItem(t *testing.T) {
	full := true
	tt := ui.NewTester(func(c *ui.Context) {
		first := func() {
			if full {
				ui.Box(c).Height(120)
			}
		}
		Masonry(c, 2, []func(){first, block(c, "I1", 20), block(c, "I2", 20)})
	}, 400, 400)
	tt.Advance(time2s)
	full = false
	tt.Frame()
	tt.Advance(time2s)
	if tt.FramePending() {
		t.Fatal("Masonry kept rebuilding after an item emptied")
	}
	if find(t, tt, "I2").Y >= find(t, tt, "I1").Y+20 {
		t.Error("I2 did not move into the emptied column")
	}
}

func TestMasonrySettlesAtZeroWidth(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(0).ClipX().Children(func() {
			Masonry(c, 2, []func(){block(c, "I0", 40), block(c, "I1", 20)})
		})
	}, 400, 400)
	tt.Advance(time2s)
	if tt.FramePending() {
		t.Error("zero-width Masonry kept rebuilding")
	}
}
