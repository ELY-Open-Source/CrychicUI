package kit

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestStaggerDelaysOnlyYoungLists(t *testing.T) {
	rows := 3
	tt := ui.NewTester(func(c *ui.Context) {
		list := ui.Column(c).Padding(10)
		list.Children(func() {
			for i := range rows {
				c.Key(i)
				Stagger(c, list, ui.Box(c).Size(100, 20).Background(c.Theme().Accent), i)
			}
		})
	}, 200, 200)
	tt.Advance(40 * time.Millisecond)
	first, third := tt.Image().RGBAAt(60, 20), tt.Image().RGBAAt(60, 60)
	if first == third {
		t.Error("rows arrive together while the list is young")
	}
	tt.Advance(2 * time.Second)
	rows = 4
	tt.Frame()
	tt.Advance(250 * time.Millisecond)
	if tt.FramePending() {
		t.Error("a late row waits for a stagger delay")
	}
}

func TestCardHeadAndChip(t *testing.T) {
	var accent ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		accent = c.Theme().Accent
		Card(c).Width(300).Children(func() {
			CardHead(c, "file", "Invoices", func() { Chip(c, "3 new", accent) })
		})
	}, 400, 200)
	tt.Frame()
	if !tt.HasText("Invoices") || !tt.HasText("3 new") {
		t.Errorf("texts %q", tt.Texts())
	}
	title, chip := mustFind(t, tt, "Invoices"), mustFind(t, tt, "3 new")
	if chip.X <= title.X || chip.Y > title.Y+title.H {
		t.Errorf("chip not trailing the title: %v %v", title, chip)
	}
}

func mustFind(t *testing.T, tt *ui.Tester, s string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(s)
	if !ok {
		t.Fatalf("%q missing", s)
	}
	return r
}
