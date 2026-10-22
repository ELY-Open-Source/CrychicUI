package navigation

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCrumbsFold(t *testing.T) {
	path := []string{"Home", "Work", "Projects", "Crychic", "Docs", "Readme"}
	chosen := -1
	tt := ui.NewTester(func(c *ui.Context) { Crumbs(c, "Path", path, &chosen, CrumbsOptions{}) }, 700, 100)
	if tt.HasText("Projects") || !tt.HasText("Docs") || !tt.HasText("Readme") {
		t.Fatalf("shown %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Home")
	if chosen != 0 {
		t.Errorf("chosen %d", chosen)
	}
	click(t, tt, "More")
	if !slices.Equal(tt.Menu(), []string{"Work", "Projects", "Crychic"}) {
		t.Fatalf("menu %q", tt.Menu())
	}
	if err := tt.ChooseMenuItem("Projects"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if chosen != 2 {
		t.Errorf("from the menu: chosen %d", chosen)
	}
	cues(t, before, sound.Tap, sound.Tap)
}

func TestPaginationSlots(t *testing.T) {
	for _, c := range []struct {
		page, pages int
		want        []int
	}{{1, 5, []int{1, 2, 3, 4, 5}}, {10, 20, []int{1, 0, 9, 10, 11, 0, 20}}, {2, 20, []int{1, 2, 3, 4, 5, 0, 20}},
		{19, 20, []int{1, 0, 16, 17, 18, 19, 20}}, {1, 1, []int{1}}} {
		if got := slots(c.page, c.pages); !slices.Equal(got, c.want) {
			t.Errorf("slots(%d, %d) = %v, want %v", c.page, c.pages, got, c.want)
		}
	}
	page := 10
	tt := ui.NewTester(func(c *ui.Context) { Pagination(c, "Pages", &page, 20, PaginationOptions{}) }, 600, 100)
	before := baseline()
	click(t, tt, "Page 11")
	if page != 11 || !tt.FramePending() {
		t.Errorf("page %d, popping %v", page, tt.FramePending())
	}
	click(t, tt, "Next")
	if page != 12 {
		t.Errorf("next: page %d", page)
	}
	cues(t, before, sound.Tick, sound.Tick)
	settles(t, tt, "Pagination")
	panics(t, "page 0", func() { render(func(c *ui.Context) { p := 0; Pagination(c, "x", &p, 3, PaginationOptions{}) }) })
}

func TestPager(t *testing.T) {
	page := 1
	tt := ui.NewTester(func(c *ui.Context) { Pager(c, "Pages", &page, 5, PagerOptions{}) }, 500, 100)
	before := baseline()
	click(t, tt, "Next")
	if page != 2 || !tt.HasText("Page 2 of 5") {
		t.Errorf("page %d, %q", page, tt.Texts())
	}
	cues(t, before, sound.Tick)
	settles(t, tt, "Pager")
}

func TestLoadMore(t *testing.T) {
	shown := 10
	tt := ui.NewTester(func(c *ui.Context) {
		LoadMore(c, "Posts", &shown, 25, func(i int) { ui.Text(c, fmt.Sprintf("Post %d", i)) }, LoadMoreOptions{})
	}, 400, 1200)
	settles(t, tt, "LoadMore")
	before := baseline()
	click(t, tt, "Show more")
	if shown != 20 || !tt.FramePending() {
		t.Errorf("shown %d, staggering %v", shown, tt.FramePending())
	}
	settles(t, tt, "more")
	click(t, tt, "Show more")
	settles(t, tt, "the rest")
	if shown != 25 || tt.HasText("Show more") || !tt.HasText("Post 24") {
		t.Errorf("shown %d, %q", shown, tt.Texts())
	}
	cues(t, before, sound.Tap, sound.Tap)
}

func TestInfiniteListAsksOnce(t *testing.T) {
	count, loading, asks := 30, false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		InfiniteList(c, "Feed", count, loading, func(i int) { ui.Box(c).Height(30).Label(fmt.Sprintf("Row %d", i)) },
			func() { asks++; loading = true }, InfiniteListOptions{}).Height(200)
	}, 300, 300)
	settles(t, tt, "InfiniteList")
	if asks != 0 {
		t.Fatalf("asked %d times before scrolling", asks)
	}
	r := find(t, tt, "Feed")
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 5000)
	tt.Advance(time.Second)
	if asks != 1 || !tt.HasText("Loading") {
		t.Fatalf("near the end: asks %d, %q", asks, tt.Texts())
	}
	tt.Advance(time.Second)
	if asks != 1 {
		t.Errorf("asked again: %d", asks)
	}
	count, loading = 40, false
	tt.Frame()
	tt.Advance(time.Second)
	if asks != 1 {
		t.Errorf("asked for more while away from the end: %d", asks)
	}
	short, shortAsks := 2, 0
	ts := ui.NewTester(func(c *ui.Context) {
		InfiniteList(c, "Feed", short, false, func(i int) { ui.Box(c).Height(30) }, func() { shortAsks++ }, InfiniteListOptions{}).Height(200)
	}, 300, 300)
	ts.Advance(time.Second)
	if shortAsks != 1 {
		t.Errorf("a short list asked %d times", shortAsks)
	}
}

func TestInfiniteListRowsFromAClick(t *testing.T) {
	count, loading, asks := 30, false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if ui.ButtonBase(c).Label("Done").Size(40, 20).Clicked() {
			count, loading = count+10, false
		}
		InfiniteList(c, "Feed", count, loading, func(i int) { ui.Box(c).Height(30) },
			func() { asks++; loading = true }, InfiniteListOptions{}).Height(200)
	}, 300, 300)
	settles(t, tt, "InfiniteList")
	r := find(t, tt, "Feed")
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 5000)
	tt.Advance(time.Second)
	click(t, tt, "Done")
	tt.Advance(time.Second)
	if asks != 1 || count != 40 {
		t.Errorf("rows added by a click: asks %d, count %d", asks, count)
	}
}

func TestSteps(t *testing.T) {
	current := 0
	tt := ui.NewTester(func(c *ui.Context) { Steps(c, "Checkout", []string{"Cart", "Address", "Pay"}, &current) }, 600, 150)
	before := baseline()
	click(t, tt, "Pay")
	if current != 2 || !tt.FramePending() {
		t.Errorf("current %d, filling %v", current, tt.FramePending())
	}
	cues(t, before, sound.Select)
	settles(t, tt, "Steps")
	panics(t, "one step", func() { render(func(c *ui.Context) { s := 0; Steps(c, "x", []string{"A"}, &s) }) })
}
