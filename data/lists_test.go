package data

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// swipe drags by dx in frame steps, then releases.
func swipe(tt *ui.Tester, x, y, dx float32) {
	tt.Press(x, y)
	for i := 1; i <= 8; i++ {
		tt.Move(x+dx*float32(i)/8, y)
		tt.Advance(16 * time.Millisecond)
	}
	tt.Release(x+dx, y)
	tt.Frame()
}

func TestVirtualList(t *testing.T) {
	chosen, changed := -1, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if VirtualList(c, "Rows", 10000, &chosen, func(i int) { ui.Text(c, fmt.Sprint("Row ", i)) }, VirtualListOptions{}).Changed() {
			changed++
		}
	}, 400, 400)
	if !tt.HasText("Row 0") || tt.HasText("Row 9999") || !tt.FramePending() {
		t.Fatalf("not a virtual list fading in: %d texts", len(tt.Texts()))
	}
	settles(t, tt, "list")
	before := baseline()
	click(t, tt, "Row 2")
	if chosen != 2 || changed != 1 {
		t.Errorf("click: chosen %d, changed %d", chosen, changed)
	}
	cues(t, before, sound.Select)
	panics(t, "height", func() {
		render(func(c *ui.Context) { VirtualList(c, "x", 1, nil, func(int) {}, VirtualListOptions{Height: -1}) })
	})
}

func TestGroupedList(t *testing.T) {
	groups := []Group{{"Fruit", []string{"Apple", "Pear"}}, {"Greens", []string{"Kale"}}}
	chosen := [2]int{-1, -1}
	tt := ui.NewTester(func(c *ui.Context) { GroupedList(c, "Food", groups, &chosen, GroupedListOptions{}) }, 400, 400)
	before := baseline()
	click(t, tt, "Kale")
	if chosen != [2]int{1, 0} {
		t.Errorf("chose %v", chosen)
	}
	cues(t, before, sound.Select)
	panics(t, "not an item", func() {
		render(func(c *ui.Context) { ch := [2]int{0, 5}; GroupedList(c, "x", groups, &ch, GroupedListOptions{}) })
	})
}

func TestReorderList(t *testing.T) {
	items := []string{"a", "b", "c"}
	tt := ui.NewTester(func(c *ui.Context) { ReorderList(c, "Steps", &items, ReorderListOptions{}) }, 400, 400)
	before := baseline()
	click(t, tt, "b")
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Frame()
	if !slices.Equal(items, []string{"b", "a", "c"}) {
		t.Fatalf("Alt+Up: %v", items)
	}
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Key(ui.Alt, ui.KeyDown)
	tt.Frame()
	if !slices.Equal(items, []string{"a", "c", "b"}) {
		t.Errorf("Alt+Down twice: %v", items)
	}
	cues(t, before, sound.Drop, sound.Drop, sound.Drop)
}

func TestSwipeList(t *testing.T) {
	items, archived := []string{"Mail", "News", "Bills"}, []string(nil)
	tt := ui.NewTester(func(c *ui.Context) { SwipeList(c, "Inbox", &items, &archived) }, 400, 400)
	settles(t, tt, "list")
	before := baseline()
	r := find(t, tt, "News")
	swipe(tt, r.X+r.W/2, r.Y+r.H/2, -300)
	// The row slides off first, then the list closes up.
	tt.Advance(time.Second)
	settles(t, tt, "deleted")
	if !slices.Equal(items, []string{"Mail", "Bills"}) || archived != nil {
		t.Fatalf("swipe left: items %v, archived %v", items, archived)
	}
	r = find(t, tt, "Bills")
	swipe(tt, r.X+r.W/2, r.Y+r.H/2, 300)
	tt.Advance(time.Second)
	settles(t, tt, "archived")
	if !slices.Equal(items, []string{"Mail"}) || !slices.Equal(archived, []string{"Bills"}) {
		t.Fatalf("swipe right: items %v, archived %v", items, archived)
	}
	r = find(t, tt, "Mail")
	// Held still before release: no throw.
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W/2+60, r.Y+r.H/2)
	tt.Advance(300 * time.Millisecond)
	tt.Release(r.X+r.W/2+60, r.Y+r.H/2)
	tt.Frame()
	settles(t, tt, "short swipe")
	if len(items) != 1 {
		t.Fatalf("a short swipe removed the row")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDelete)
	tt.Frame()
	if len(items) != 0 {
		t.Errorf("Delete key: %v", items)
	}
	cues(t, before, sound.Delete, sound.Swipe, sound.Delete)
}

func TestCardList(t *testing.T) {
	cards := []Card{{"Free", "For one", "user"}, {"Team", "For five", "users"}}
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) { CardList(c, "Plans", cards, &chosen) }, 400, 400)
	settles(t, tt, "cards")
	r := find(t, tt, "Team")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("no lift on hover")
	}
	before := baseline()
	click(t, tt, "Team")
	if chosen != 1 {
		t.Errorf("chose %d", chosen)
	}
	cues(t, before, sound.Tap)
}

func TestGridList(t *testing.T) {
	tiles := []Tile{{"Photos", "image"}, {"Music", "music"}, {"Docs", "file-text"}}
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) { GridList(c, "Apps", tiles, &chosen, GridListOptions{}) }, 500, 400)
	before := baseline()
	click(t, tt, "Music")
	if chosen != 1 {
		t.Errorf("chose %d", chosen)
	}
	cues(t, before, sound.Select)
}

func TestColumnPicker(t *testing.T) {
	order := []string{"Name", "Size", "Date"}
	var hidden ui.Selection[string]
	tt := ui.NewTester(func(c *ui.Context) { ColumnPicker(c, "Columns", &order, &hidden) }, 400, 400)
	before := baseline()
	r := find(t, tt, "Size")
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !hidden.Has("Size") {
		t.Fatalf("Size not hidden")
	}
	rows, _ := tt.Find("Date")
	tt.ClickAt(rows.X+rows.W+40, rows.Y+rows.H/2)
	tt.Key(ui.Alt, ui.KeyUp)
	tt.Frame()
	if !slices.Equal(order, []string{"Name", "Date", "Size"}) {
		t.Errorf("Alt+Up: %v", order)
	}
	cues(t, before, sound.Toggle, sound.Drop)
}
