package files

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestColumnBrowserScrollsLongColumns(t *testing.T) {
	var names []string
	for i := range 20 {
		names = append(names, fmt.Sprintf("invoice-%02d.pdf", i+1))
	}
	selected := []int{-1}
	tt := ui.NewTester(padded(func(c *ui.Context) { ColumnBrowser(c, "Browse", [][]string{names}, &selected) }), 500, 400)
	settles(t, tt)
	first := find(t, tt, names[0])
	tt.Scroll(first.X+first.W/2, first.Y+first.H/2, 0, 2000)
	settles(t, tt)
	if r := find(t, tt, names[19]); r.H == 0 {
		t.Fatal("the last row has no height after scrolling")
	}
	click(t, tt, names[19])
	if selected[0] != 19 {
		t.Errorf("selected %v", selected)
	}
}

func TestColumnBrowserKeysWholePaths(t *testing.T) {
	columns := [][]string{{"Finance", "Marketing"}, {"Reports"}, {"budget.xlsx"}}
	selected := []int{0, 0, -1}
	tt := ui.NewTester(padded(func(c *ui.Context) { ColumnBrowser(c, "Browse", columns, &selected) }), 600, 300)
	settles(t, tt)
	r := find(t, tt, "budget.xlsx")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	selected[0] = 1
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
	if selected[2] != -1 {
		t.Errorf("a file under a new ancestor took the held click: %v", selected)
	}
}

func TestCopyNeverBecomesAKey(t *testing.T) {
	render := func(view func(c *ui.Context)) { ui.Render(view, 500, 400, 1) }
	render(func(c *ui.Context) {
		Storage(c, "Storage", []StoragePart{{Name: "Free", Bytes: 1}}, 10, StorageOptions{})
	})
	render(func(c *ui.Context) { Compress(c, "Zip", 10, 5, CompressOptions{Before: "Size", After: "Size"}) })
	f := FileItem{Name: "a.txt"}
	choice := ConflictNone
	render(func(c *ui.Context) {
		Conflict(c, "a.txt", f, f, &choice, ConflictOptions{Yours: "Copy", Theirs: "Copy"})
	})
}

func TestZipListTotals(t *testing.T) {
	opened := -1
	half := []ZipEntry{{Name: "a", Size: math.MaxInt64 / 2, Compressed: math.MaxInt64 / 4}, {Name: "b", Size: math.MaxInt64 / 2, Compressed: math.MaxInt64 / 4}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ZipList(c, "big.zip", half, &opened, ZipListOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("−50%") {
		t.Errorf("savings %q", tt.Texts())
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			ZipList(c, "x", []ZipEntry{{Name: "a", Size: math.MaxInt64}, {Name: "b", Size: 1}}, &opened, ZipListOptions{})
		}, 500, 300, 1)
	})
}

func TestDownloadRowLongWait(t *testing.T) {
	opened := false
	item := FileItem{Name: "archive.tar", Size: 10_000_000_000}
	tt := ui.NewTester(padded(func(c *ui.Context) { DownloadRow(c, item, 0, 1, &opened, DownloadRowOptions{}) }), 500, 200)
	settles(t, tt)
	if strings.Contains(strings.Join(tt.Texts(), " "), "-") || !tt.HasText("115740d left") {
		t.Errorf("texts %q", tt.Texts())
	}
	for secs, want := range map[float64]string{45.9: "45s", 724: "12m 4s", 5400: "1h 30m", 1e10: "115740d"} {
		if got := leftText(secs); got != want {
			t.Errorf("leftText(%v) = %q, want %q", secs, got, want)
		}
	}
}

func TestChipsStayInsideAndHover(t *testing.T) {
	long := strings.Repeat("contract", 25)
	flags := []bool{false}
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(300).Padding(10).Gap(8).Children(func() {
			FileTags(c, "Tags", []string{long}, &flags)
			FileTypeFilter(c, "Types", []TypeCount{{Kind: long, Sample: "a.pdf", Count: 3}}, &chosen)
		})
	}, 600, 300)
	settles(t, tt)
	tag, kind := find(t, tt, long), find(t, tt, long+" 3")
	if tag.X+tag.W > 311 || kind.X+kind.W > 311 {
		t.Errorf("chips run to %v and %v past 310", tag.X+tag.W, kind.X+kind.W)
	}
	tt.Move(tag.X+tag.W/2, tag.Y+tag.H/2)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("a tag does not answer hover")
	}
}
