package layout

import (
	"fmt"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

const time2s = 2 * time.Second

func lines(c *ui.Context, n int) func() {
	return func() {
		for i := range n {
			ui.Text(c, fmt.Sprintf("Line %d", i)).Padding(4)
		}
	}
}

func TestScrollAreaEdges(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ScrollArea(c, lines(c, 60)).Size(300, 200)
	}, 300, 200)
	tt.Advance(time2s)
	top, bottom := lum(tt.Image(), 290, 1), lum(tt.Image(), 290, 198)
	if bottom >= top {
		t.Fatal("no bottom shadow while content continues")
	}
	tt.Scroll(150, 100, 0, 200)
	tt.Advance(time2s)
	if lum(tt.Image(), 290, 1) >= top {
		t.Error("no top shadow once scrolled")
	}
}

func TestScrollerPages(t *testing.T) {
	forward := 0
	tt := ui.NewTester(func(c *ui.Context) {
		_, _, fwd := Scroller(c, func() {
			for i := range 12 {
				ui.Column(c).Size(90, 40).Children(func() { ui.Text(c, fmt.Sprintf("Tile %d", i)) })
			}
		})
		if fwd.Clicked() {
			forward++
		}
	}, 400, 100)
	taps := count(sound.Tap)
	if err := tt.Click("Scroll back"); err != nil {
		t.Fatal(err)
	}
	if count(sound.Tap) != taps {
		t.Error("disabled back arrow played a cue")
	}
	x5 := find(t, tt, "Tile 5").X
	if err := tt.Click("Scroll forward"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(time2s)
	if count(sound.Tap) != taps+1 || forward != 1 {
		t.Errorf("forward arrow: %d cues, caller saw %d clicks", count(sound.Tap)-taps, forward)
	}
	if x := find(t, tt, "Tile 5").X; x > x5-200 {
		t.Errorf("Tile 5 moved from %v to %v, less than a page", x5, x)
	}
}

func TestStickyHeaderStays(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		StickyHeader(c, func() { ui.Text(c, "Head") }, lines(c, 60)).Size(300, 200)
	}, 300, 200)
	tt.Advance(time2s)
	y, below := find(t, tt, "Head").Y, find(t, tt, "Line 0")
	shade := lum(tt.Image(), 290, int(below.Y)+2)
	tt.Scroll(150, 150, 0, 300)
	tt.Advance(time2s)
	if find(t, tt, "Head").Y != y {
		t.Error("header moved with the scroll")
	}
	if lum(tt.Image(), 290, int(below.Y)+2) >= shade {
		t.Error("no header shadow once scrolled")
	}
}

func TestScrollAreaEdgeFollowsContent(t *testing.T) {
	n := 3
	tt := ui.NewTester(func(c *ui.Context) {
		ScrollArea(c, lines(c, n)).Size(300, 200)
	}, 300, 200)
	tt.Advance(time2s)
	flat := lum(tt.Image(), 290, 198)
	n = 60
	tt.Frame()
	tt.Advance(time2s)
	if lum(tt.Image(), 290, 198) >= flat {
		t.Error("no bottom shadow once content overflows")
	}
}
