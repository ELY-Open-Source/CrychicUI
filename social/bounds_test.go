package social

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

// narrow puts a view in a 320-DIP column.
func narrow(view func(c *ui.Context)) func(c *ui.Context) {
	return padded(func(c *ui.Context) { ui.Column(c).Width(320).Children(func() { view(c) }) })
}

func inside(t *testing.T, tt *ui.Tester, name string) ui.Rect {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok || r.W == 0 {
		t.Fatalf("%q missing: %q", name, tt.Texts())
	}
	if r.X < 10 || r.X+r.W > 330.5 {
		t.Errorf("%q spans %.1f..%.1f, outside 10..330", name, r.X, r.X+r.W)
	}
	return r
}

func TestLikeCountNoOverflow(t *testing.T) {
	state := PostState{}
	post := PostData{ID: "p1", Text: "Kiln day", Author: people()[0], Likes: math.MaxInt}
	tt := ui.NewTester(padded(func(c *ui.Context) { Post(c, post, &state, PostOptions{}) }), 600, 500)
	settles(t, tt)
	click(t, tt, "Like")
	settles(t, tt)
	if !exact(tt, "9,223,372,036,854,775,808") {
		t.Errorf("liked MaxInt: %q", tt.Texts())
	}
}

func TestCommentIDsUniqueInTree(t *testing.T) {
	u := people()
	comments := []Comment{{ID: "c1", Text: "First", Author: u[0], Replies: []Comment{{ID: "c1", Text: "Echo", Author: u[1]}}}}
	panics(t, "CommentThread", func(c *ui.Context) { CommentThread(c, comments, new(string), CommentThreadOptions{}) })
}

func TestMentionAtCaret(t *testing.T) {
	text, selected := " coffee later", ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Mentions(c, &text, people(), &selected, MentionsOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Write with @mentions")
	tt.Key(0, ui.KeyHome)
	tt.Type("@pr")
	settles(t, tt)
	click(t, tt, "@priyanair")
	settles(t, tt)
	if text != "@priyanair coffee later" || selected != "priya" {
		t.Errorf("text %q selected %q", text, selected)
	}
}

func TestHugeCountsStayInside(t *testing.T) {
	state := PostState{}
	post := PostData{ID: "p1", Text: "Kiln day", Author: people()[0], Likes: math.MaxInt - 2, Comments: math.MaxInt - 1, Shares: math.MaxInt}
	tt := ui.NewTester(narrow(func(c *ui.Context) { Post(c, post, &state, PostOptions{}) }), 600, 700)
	settles(t, tt)
	for _, n := range []string{"9,223,372,036,854,775,805", "9,223,372,036,854,775,806", "9,223,372,036,854,775,807"} {
		inside(t, tt, n)
	}

	tt = ui.NewTester(narrow(func(c *ui.Context) { ProfileStats(c, math.MaxInt-2, math.MaxInt-1, math.MaxInt, ProfileStatsOptions{}) }), 600, 400)
	settles(t, tt)
	var prev ui.Rect
	for i, n := range []string{"9,223,372,036,854,775,805", "9,223,372,036,854,775,806", "9,223,372,036,854,775,807"} {
		r := inside(t, tt, n)
		if i > 0 && r.X < prev.X+prev.W {
			t.Errorf("%q starts at %.1f, inside the previous count ending %.1f", n, r.X, prev.X+prev.W)
		}
		prev = r
	}
}

func TestLongHashtagStaysInside(t *testing.T) {
	opened := false
	tt := ui.NewTester(narrow(func(c *ui.Context) {
		Hashtag(c, "#handmadeceramicsandcontemporarydesign", 12480, &opened, HashtagOptions{})
	}), 600, 200)
	settles(t, tt)
	inside(t, tt, "#handmadeceramicsandcontemporarydesign")
}

func TestPanicsNameComponent(t *testing.T) {
	bad := User{Name: "Maya Chen"}
	panics(t, "ProfileHeader", func(c *ui.Context) { ProfileHeader(c, bad, new(bool), ProfileHeaderOptions{}) })
	panics(t, "Feed", func(c *ui.Context) {
		Feed(c, []PostData{{ID: "p1", Author: people()[0]}}, new(bool), FeedOptions{})
	})
}
