package typography

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestBlockquote(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Blockquote(c, "Quoted words") }, 300, 200)
	settles(t, tt, "Blockquote bar")
	find(t, tt, "Quoted words")
}

func TestPullQuote(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { PullQuote(c, "A quote", "Its author") }, 300, 200)
	settles(t, tt, "PullQuote")
	find(t, tt, "A quote")
	find(t, tt, "Its author")
}

func TestHighlight(t *testing.T) {
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		Highlight(c, "mm")
	}, 300, 200)
	if !near(tt.Image(), 1, 17, th.Background) {
		t.Error("marker showed before its sweep")
	}
	settles(t, tt, "Highlight sweep")
	if near(tt.Image(), 1, 17, th.Background) {
		t.Error("marker never swept in")
	}
}

func TestMark(t *testing.T) {
	var th ui.Theme
	tt := ui.NewTester(func(c *ui.Context) {
		th = *c.Theme()
		Mark(c, "match")
	}, 300, 200)
	if !near(tt.Image(), 1, 1, th.Background) {
		t.Error("mark showed before its fade")
	}
	settles(t, tt, "Mark fade")
	if near(tt.Image(), 1, 1, th.Background) {
		t.Error("mark never faded in")
	}
}

func TestAbbr(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Abbr(c, "API", "Application Programming Interface") }, 400, 200)
	settles(t, tt, "Abbr dots")
	r := find(t, tt, "API")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(time.Second)
	if !tt.HasText("Application Programming Interface") {
		t.Error("full term tooltip missing")
	}
}

func TestTruncate(t *testing.T) {
	const text = "A long title that cannot fit on one line"
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(80).Children(func() { Truncate(c, text) })
	}, 300, 200)
	r := find(t, tt, text)
	if r.W > 80.5 {
		t.Errorf("truncated width %v", r.W)
	}
	before := textCount(tt, text)
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(time.Second)
	if textCount(tt, text) <= before {
		t.Error("full text tooltip missing")
	}
}
