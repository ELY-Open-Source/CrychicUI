package social

import (
	"fmt"
	"image"
	"image/color"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func people() []User {
	return []User{
		{ID: "maya", Name: "Maya Chen", Handle: "@mayachen"},
		{ID: "daniel", Name: "Daniel Okafor", Handle: "@dokafor"},
		{ID: "priya", Name: "Priya Nair", Handle: "@priyanair"},
	}
}

func exact(tt *ui.Tester, s string) bool { return slices.Contains(tt.Texts(), s) }

func panics(t *testing.T, want string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want %q", r, want)
		}
	}()
	ui.Render(view, 600, 600, 1)
}

// swapUnderPress presses name, reorders, then releases there.
func swapUnderPress(t *testing.T, tt *ui.Tester, name string, reorder func()) {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("%q missing: %q", name, tt.Texts())
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	reorder()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func TestFollowListsKeepIdentity(t *testing.T) {
	for _, tc := range []struct {
		name, press string
		view        func(c *ui.Context, users []User, m *map[string]bool)
	}{
		{"WhoToFollow", "Follow Daniel Okafor", func(c *ui.Context, users []User, m *map[string]bool) {
			WhoToFollow(c, users, m, WhoToFollowOptions{})
		}},
		{"Followers", "Follow back Daniel Okafor", func(c *ui.Context, users []User, m *map[string]bool) {
			Followers(c, users, m, FollowersOptions{})
		}},
	} {
		users, following := people(), map[string]bool{}
		tt := ui.NewTester(padded(func(c *ui.Context) { tc.view(c, users, &following) }), 500, 400)
		settles(t, tt)
		swapUnderPress(t, tt, tc.press, func() { users[0], users[1] = users[1], users[0] })
		if following["maya"] {
			t.Errorf("%s: the row that slid under the pointer was followed: %v", tc.name, following)
		}
	}
}

func TestFollowKeepsFocus(t *testing.T) {
	users, following := people(), map[string]bool{}
	tt := ui.NewTester(padded(func(c *ui.Context) { WhoToFollow(c, users, &following, WhoToFollowOptions{}) }), 500, 400)
	settles(t, tt)
	click(t, tt, "Follow Priya Nair")
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if following["priya"] {
		t.Errorf("Space after the first click did not unfollow: %v", following)
	}
}

func TestStoryRingKeepsIdentity(t *testing.T) {
	u := people()
	stories := []Story{{ID: "s1", User: u[0]}, {ID: "s2", User: u[1]}, {ID: "s3", User: u[2]}}
	opened := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { StoryRing(c, stories, &opened) }), 500, 250)
	settles(t, tt)
	swapUnderPress(t, tt, "Daniel Okafor", func() { stories[0], stories[1] = stories[1], stories[0] })
	if opened == "s1" {
		t.Errorf("opened %q, the story that slid in", opened)
	}
}

func TestRepeatsPanic(t *testing.T) {
	u := people()
	post := PostData{ID: "p1", Text: "Kiln day", Author: u[0]}
	cases := map[string]func(c *ui.Context){
		"Feed": func(c *ui.Context) { Feed(c, []PostData{post, post}, new(bool), FeedOptions{}) },
		"WhoToFollow": func(c *ui.Context) {
			m := map[string]bool{}
			WhoToFollow(c, []User{u[0], u[0]}, &m, WhoToFollowOptions{})
		},
		"StoryRing": func(c *ui.Context) { StoryRing(c, []Story{{ID: "s", User: u[0]}, {ID: "s", User: u[1]}}, new(string)) },
		"PollPost": func(c *ui.Context) {
			choices := []PollChoice{{Text: "Oat"}, {Text: "Oat"}}
			voted := -1
			PollPost(c, "Glaze?", &choices, time.Now().Add(time.Hour), &voted, PollPostOptions{})
		},
		"Trending": func(c *ui.Context) { Trending(c, []Topic{{Tag: "#a"}, {Tag: "#a"}}, TrendingOptions{}) },
		"ProfileBadges": func(c *ui.Context) {
			ProfileBadges(c, []Badge{{ID: "b", Name: "A", Icon: "award"}, {ID: "b", Name: "B", Icon: "star"}})
		},
		"ReactionBar": func(c *ui.Context) {
			r := []Reaction{{ID: "love", Name: "Love", Icon: "heart"}, {ID: "love", Name: "Like", Icon: "heart"}}
			ReactionBar(c, &r)
		},
		"Mentions": func(c *ui.Context) { Mentions(c, new(string), []User{u[0], u[0]}, new(string), MentionsOptions{}) },
	}
	for name, view := range cases {
		panics(t, name, view)
	}
}

func TestPostCountsIncludeViewersLike(t *testing.T) {
	state := PostState{}
	post := PostData{ID: "p1", Text: "Kiln day", Author: people()[0], Likes: 1214}
	tt := ui.NewTester(padded(func(c *ui.Context) { Post(c, post, &state, PostOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "1,214") {
		t.Fatalf("texts %q", tt.Texts())
	}
	click(t, tt, "Like")
	settles(t, tt)
	if !exact(tt, "1,215") {
		t.Errorf("liked count: %q", tt.Texts())
	}
}

func TestComposerMedia(t *testing.T) {
	draft, added := "Kiln day", false
	photo := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 4, 4)))
	tt := ui.NewTester(padded(func(c *ui.Context) { PostComposer(c, &draft, photo, &added, PostComposerOptions{}) }), 500, 500)
	settles(t, tt)
	if _, ok := tt.Find("Remove media"); ok {
		t.Fatal("preview shown before media was added")
	}
	click(t, tt, "Add media")
	settles(t, tt)
	before := baseline()
	click(t, tt, "Remove media")
	settles(t, tt)
	if added {
		t.Error("remove left media added")
	}
	if _, ok := tt.Find("Remove media"); ok {
		t.Error("preview stayed after remove")
	}
	cues(t, before, sound.Pop)

	tt = ui.NewTester(padded(func(c *ui.Context) { PostComposer(c, &draft, nil, &added, PostComposerOptions{}) }), 500, 500)
	settles(t, tt)
	if err := tt.Click("Add media"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if added {
		t.Error("added media with no image")
	}
}

func TestPollResultsAfterVote(t *testing.T) {
	choices := []PollChoice{{Text: "Speckled oat"}, {Text: "Deep sea blue"}}
	voted := -1
	var ends time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ends.IsZero() {
			ends = c.Now().Add(90 * time.Minute)
		}
		PollPost(c, "Glaze?", &choices, ends, &voted, PollPostOptions{})
	}), 500, 400)
	settles(t, tt)
	if exact(tt, "0%") || !exact(tt, "0 votes") || !exact(tt, "1h left") {
		t.Fatalf("before voting: %q", tt.Texts())
	}
	click(t, tt, "Speckled oat")
	settles(t, tt)
	if !exact(tt, "100%") || !exact(tt, "0%") || !exact(tt, "1 vote") {
		t.Errorf("after voting: %q", tt.Texts())
	}
	tt.Advance(2 * time.Hour)
	if !exact(tt, "Final results") {
		t.Errorf("after the end: %q", tt.Texts())
	}
}

func TestShortRemaining(t *testing.T) {
	for d, want := range map[time.Duration]string{
		50 * time.Hour: "2d", 90 * time.Minute: "1h", 59*time.Minute + 59*time.Second: "59m", 1500 * time.Millisecond: "2s",
	} {
		if got := shortRemaining(d); got != want {
			t.Errorf("%v: %q, want %q", d, got, want)
		}
	}
}

func TestMentionsEmptyAndMark(t *testing.T) {
	text, selected := "Thanks @zz", ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Mentions(c, &text, people(), &selected, MentionsOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "No one matches") {
		t.Errorf("no empty state: %q", tt.Texts())
	}
	text = "Thanks @DO"
	tt = ui.NewTester(padded(func(c *ui.Context) { Mentions(c, &text, people(), &selected, MentionsOptions{}) }), 500, 400)
	settles(t, tt)
	if _, ok := tt.Find("@dokafor"); !ok || exact(tt, "No one matches") {
		t.Errorf("case-folded match: %q", tt.Texts())
	}
}

func TestUserCardOnHover(t *testing.T) {
	open := false
	u := people()[0]
	u.Bio = "Product designer"
	tt := ui.NewTester(padded(func(c *ui.Context) { UserCard(c, u, &open) }), 500, 400)
	settles(t, tt)
	if exact(tt, "Product designer") {
		t.Fatal("card shown before hover")
	}
	r, _ := tt.Find("Maya Chen")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt)
	if !exact(tt, "Product designer") {
		t.Errorf("hover: %q", tt.Texts())
	}
}

func TestReactionToggleBack(t *testing.T) {
	reactions := []Reaction{{ID: "love", Name: "Love", Icon: "heart", Count: 1200, Active: true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ReactionBar(c, &reactions) }), 500, 200)
	settles(t, tt)
	click(t, tt, "Love")
	settles(t, tt)
	if reactions[0].Active || reactions[0].Count != 1199 || !exact(tt, "1,199") {
		t.Errorf("reactions %+v texts %q", reactions, tt.Texts())
	}
}

// opaque is a solid test photo.
func opaque(col color.RGBA) ui.ImageSource {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for i := range 64 {
		img.SetRGBA(i%8, i/8, col)
	}
	return ui.NewBitmap(img)
}

func TestMediaPostReadingOrder(t *testing.T) {
	cols := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {255, 255, 0, 255}}
	images := make([]ui.ImageSource, len(cols))
	for i, col := range cols {
		images[i] = opaque(col)
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { ui.Column(c).Width(400).Children(func() { MediaPost(c, images) }) }), 500, 400)
	settles(t, tt)
	img := tt.Image()
	// The second photo sits top right.
	at := img.RGBAAt(10+300, 10+60)
	if at.G < 200 || at.R > 60 {
		t.Errorf("top right is %v, want the second photo", at)
	}
}
