package data

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPinnedTable(t *testing.T) {
	var cols []Column
	for i := range 6 {
		cols = append(cols, Column{Title: fmt.Sprintf("C%d", i), Width: 120})
	}
	tt := ui.NewTester(func(c *ui.Context) {
		PinnedTable(c, "Wide", cols, 3, func(r, col int) { ui.Text(c, fmt.Sprintf("r%dc%d", r, col)) })
	}, 400, 300)
	pin, far := find(t, tt, "r1c0"), find(t, tt, "r1c1")
	tt.Scroll(far.X, far.Y, 300, 0)
	tt.Frame()
	if p := find(t, tt, "r1c0"); p.X != pin.X {
		t.Errorf("pinned cell moved from %v to %v", pin.X, p.X)
	}
	if f, ok := tt.Find("r1c1"); ok && f.W > 0 && f.X >= far.X {
		t.Errorf("rest did not scroll: %v then %v", far, f)
	}
	if a, b := find(t, tt, "r2c0"), find(t, tt, "r2c5"); a.Y != b.Y {
		t.Errorf("row 2 misaligned: %v and %v", a.Y, b.Y)
	}
	panics(t, "unsized", func() {
		render(func(c *ui.Context) { PinnedTable(c, "x", []Column{{Title: "a"}, {Title: "b"}}, 1, func(int, int) {}) })
	})
}

func TestTotalsRow(t *testing.T) {
	sums := []int64{10, 7}
	tt := ui.NewTester(func(c *ui.Context) {
		TotalsRow(c, "Totals", []Column{{Title: "Item"}, {Title: "Qty"}, {Title: "Cost"}}, sums, TotalsRowOptions{})
	}, 500, 100)
	if !tt.HasText("Total") || !tt.HasText("10") {
		t.Fatalf("shows %q", tt.Texts())
	}
	sums = []int64{11, 7}
	tt.Frame()
	if !tt.FramePending() {
		t.Error("sum did not roll")
	}
	settles(t, tt, "roll")
	panics(t, "count", func() {
		render(func(c *ui.Context) { TotalsRow(c, "x", []Column{{Title: "a"}}, []int64{1}, TotalsRowOptions{}) })
	})
}

func TestHeatTable(t *testing.T) {
	th := ui.LightTheme()
	values := [][]float64{{0, 10}}
	tt := ui.NewTester(func(c *ui.Context) {
		HeatTable(c, "Load", []string{"Mon"}, []string{"AM", "PM"}, values, HeatTableOptions{})
	}, 500, 200)
	settles(t, tt, "heat")
	am, pm := find(t, tt, "Mon, AM"), find(t, tt, "Mon, PM")
	hot := th.Background.Mix(th.Accent, 0.76)
	if tinted(tt, am, hot) || !tinted(tt, pm, hot) {
		t.Errorf("the high cell is not the hot one")
	}
	values = [][]float64{{10, 0}}
	tt.Frame()
	if tinted(tt, am, hot) {
		t.Error("tint jumped instead of animating")
	}
	settles(t, tt, "tint")
	if !tinted(tt, am, hot) {
		t.Error("AM did not heat")
	}
	panics(t, "NaN", func() {
		render(func(c *ui.Context) {
			HeatTable(c, "x", []string{"a"}, []string{"b"}, [][]float64{{math.NaN()}}, HeatTableOptions{})
		})
	})
}

func TestMatrix(t *testing.T) {
	var marks ui.Selection[[2]int]
	tt := ui.NewTester(func(c *ui.Context) {
		Matrix(c, "Access", []string{"Ada", "Grace"}, []string{"Read", "Write"}, &marks, MatrixOptions{})
	}, 500, 300)
	before := baseline()
	click(t, tt, "Grace, Write")
	if !marks.Has([2]int{1, 1}) || marks.Len() != 1 {
		t.Fatalf("marks %v", marks)
	}
	click(t, tt, "Grace, Write")
	if marks.Len() != 0 {
		t.Errorf("mark stayed")
	}
	cues(t, before, sound.Toggle, sound.Toggle)
}

func TestPaginated(t *testing.T) {
	page := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Paginated(c, "People", who, 25, func(r, col int) { ui.Text(c, fmt.Sprintf("p%d-%d", r, col)) }, &page, PaginatedOptions{})
	}, 500, 600)
	if !tt.HasText("p9-0") || tt.HasText("p10-0") || !tt.HasText("Page 1 of 3") {
		t.Fatalf("page 1 %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Next page")
	th := ui.LightTheme()
	// The old page shows first; both fade by 50 ms.
	tt.Advance(50 * time.Millisecond)
	if page != 1 || tinted(tt, find(t, tt, "p10-0"), th.Text) {
		t.Fatalf("page %d did not fade in: %q", page, tt.Texts())
	}
	settles(t, tt, "fade")
	if !tinted(tt, find(t, tt, "p10-0"), th.Text) {
		t.Fatal("page 2 never showed")
	}
	click(t, tt, "Next page")
	click(t, tt, "Next page")
	if page != 2 || !tt.HasText("p24-0") {
		t.Errorf("last page %d", page)
	}
	cues(t, before, sound.Tick, sound.Tick)
	panics(t, "page", func() {
		render(func(c *ui.Context) { p := 3; Paginated(c, "x", who, 25, func(int, int) {}, &p, PaginatedOptions{}) })
	})
}

func TestDiffListAndInventory(t *testing.T) {
	th := ui.LightTheme()
	tt := ui.NewTester(func(c *ui.Context) {
		DiffList(c, "Changes", []DiffItem{{ChangeAdded, "logo.svg"}, {ChangeRemoved, "old.png"}}, DiffListOptions{})
		Inventory(c, "Stock", []Stock{{"A-1", "Bolts", 2, 100}, {"A-2", "Nuts", 90, 100}}, InventoryOptions{})
	}, 700, 400)
	settles(t, tt, "lists")
	if !tt.HasText("logo.svg") || !tt.HasText("2 / 100") {
		t.Fatalf("shows %q", tt.Texts())
	}
	if !tinted(tt, find(t, tt, "Added"), th.Success) {
		t.Error("added mark is not green")
	}
	b, n := find(t, tt, "Bolts"), find(t, tt, "Nuts")
	if band(tt, b.Y-4, b.Y+b.H+4, th.Danger) == 0 || band(tt, n.Y-4, n.Y+n.H+4, th.Success) == 0 {
		t.Error("stock bars are not colored by level")
	}
	panics(t, "over", func() {
		render(func(c *ui.Context) { Inventory(c, "x", []Stock{{"a", "b", 5, 4}}, InventoryOptions{}) })
	})
	panics(t, "kind", func() {
		render(func(c *ui.Context) { DiffList(c, "x", []DiffItem{{9, "a"}}, DiffListOptions{}) })
	})
}
