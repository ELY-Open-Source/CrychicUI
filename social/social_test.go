package social

import (
	"image"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPost(t *testing.T) {
	state := PostState{}
	tt := ui.NewTester(padded(func(c *ui.Context) { Post(c, testPost(), &state, PostOptions{}) }), 500, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Like")
	if !state.Liked || !tt.HasText("Shipping") {
		t.Errorf("state %+v texts %q", state, tt.Texts())
	}
	cues(t, before, sound.Pop)
}

func TestFeed(t *testing.T) {
	fresh := true
	tt := ui.NewTester(padded(func(c *ui.Context) { Feed(c, []PostData{testPost()}, &fresh, FeedOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "New posts")
	if fresh {
		t.Error("new posts pill remains")
	}
	cues(t, before, sound.Notify)
}

func TestPostComposer(t *testing.T) {
	draft, added, sent := "Hello", false, false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, send := PostComposer(c, &draft, nil, &added, PostComposerOptions{})
		sent = sent || send.Changed()
	}), 500, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Post")
	if !sent || !tt.HasText("5/280") {
		t.Errorf("sent %v texts %q", sent, tt.Texts())
	}
	cues(t, before, sound.Send)
}

func TestPostComposerLimit(t *testing.T) {
	draft, added, sent := "too long", false, false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, send := PostComposer(c, &draft, nil, &added, PostComposerOptions{Limit: 3})
		sent = sent || send.Changed()
	}), 500, 350)
	settles(t, tt)
	if err := tt.Click("Post"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if sent {
		t.Error("over-limit post sent")
	}
}

func TestPostActions(t *testing.T) {
	state := PostState{}
	tt := ui.NewTester(padded(func(c *ui.Context) { PostActions(c, &state, PostActionsOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Share")
	click(t, tt, "Save")
	if state.Action != "share" || !state.Saved {
		t.Errorf("state %+v", state)
	}
	cues(t, before, sound.Pop, sound.Pop)
}

func TestCommentThread(t *testing.T) {
	draft, sent := "Thanks", false
	comments := []Comment{{ID: "c1", Text: "First", Author: testUser(), Replies: []Comment{{ID: "c2", Text: "Nested", Author: testUser()}}}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, send := CommentThread(c, comments, &draft, CommentThreadOptions{})
		sent = sent || send.Changed()
	}), 500, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Reply")
	if !sent || !tt.HasText("Nested") {
		t.Errorf("sent %v texts %q", sent, tt.Texts())
	}
	cues(t, before, sound.Send)
}

func TestProfileHeader(t *testing.T) {
	following := false
	tt := ui.NewTester(padded(func(c *ui.Context) { ProfileHeader(c, testUser(), &following, ProfileHeaderOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Follow")
	if !following {
		t.Error("not following")
	}
	cues(t, before, sound.Toggle)
}

func TestProfileStats(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ProfileStats(c, 12, 1200, 40, ProfileStatsOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("1200") || !tt.HasText("Followers") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestWhoToFollow(t *testing.T) {
	following := map[string]bool{}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		WhoToFollow(c, []User{testUser()}, &following, WhoToFollowOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "ana Follow")
	if !following["ana"] {
		t.Errorf("following %v", following)
	}
	cues(t, before, sound.Toggle)
}

func TestStoryRing(t *testing.T) {
	opened := ""
	stories := []Story{{ID: "s1", User: testUser()}}
	tt := ui.NewTester(padded(func(c *ui.Context) { StoryRing(c, stories, &opened) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Ana")
	if opened != "s1" {
		t.Errorf("opened %q", opened)
	}
	cues(t, before, sound.Open)
}

func TestHashtag(t *testing.T) {
	opened := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Hashtag(c, "#mygo", 42, &opened, HashtagOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "#mygo")
	if !opened || !tt.HasText("42 posts") {
		t.Errorf("opened %v texts %q", opened, tt.Texts())
	}
	cues(t, before, sound.Tap)
}

func TestTrending(t *testing.T) {
	topics := []Topic{{Tag: "#mygo", Posts: 42}, {Tag: "#ui", Posts: 20}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Trending(c, topics, TrendingOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("#mygo") || !tt.HasText("42 posts") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestPollPost(t *testing.T) {
	choices := []PollChoice{{Text: "Yes", Votes: 2}, {Text: "No", Votes: 1}}
	voted := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PollPost(c, "Ship?", &choices, c.Now().Add(time.Hour), &voted, PollPostOptions{})
	}), 500, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Yes")
	if voted != 0 || choices[0].Votes != 3 {
		t.Errorf("voted %d choices %+v", voted, choices)
	}
	cues(t, before, sound.Select)
}

func TestMediaPost(t *testing.T) {
	images := []ui.ImageSource{ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2)))}
	tt := ui.NewTester(padded(func(c *ui.Context) { MediaPost(c, images) }), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("no arrival motion")
	}
}

func TestMentions(t *testing.T) {
	text, selected := "Hi @a", ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Mentions(c, &text, []User{testUser()}, &selected, MentionsOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "@ana")
	tt.Type("!")
	tt.Frame()
	if selected != "ana" || text != "Hi @ana !" {
		t.Errorf("selected %q text %q", selected, text)
	}
	cues(t, before, sound.Select)
}

func TestUserCard(t *testing.T) {
	open := false
	tt := ui.NewTester(padded(func(c *ui.Context) { UserCard(c, testUser(), &open) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Ana")
	if !open || !tt.HasText("Designer") {
		t.Errorf("open %v texts %q", open, tt.Texts())
	}
	cues(t, before, sound.Toggle)
}

func TestFollowers(t *testing.T) {
	following := map[string]bool{}
	tt := ui.NewTester(padded(func(c *ui.Context) { Followers(c, []User{testUser()}, &following, FollowersOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "ana Follow back")
	if !following["ana"] {
		t.Errorf("following %v", following)
	}
	cues(t, before, sound.Toggle)
}

func TestProfileBadges(t *testing.T) {
	badges := []Badge{{ID: "b1", Name: "Creator", Icon: "award"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ProfileBadges(c, badges) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Creator") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestReactionBar(t *testing.T) {
	reactions := []Reaction{{ID: "heart", Name: "Heart", Icon: "heart", Count: 2}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ReactionBar(c, &reactions) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Heart")
	if !reactions[0].Active || reactions[0].Count != 3 {
		t.Errorf("reactions %+v", reactions)
	}
	cues(t, before, sound.Pop)
}

func TestMentionsClosesAfterSpace(t *testing.T) {
	text, selected := "Hi @ana ", ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Mentions(c, &text, []User{testUser()}, &selected, MentionsOptions{})
	}), 500, 300)
	settles(t, tt)
	if _, ok := tt.Find("@ana"); ok {
		t.Errorf("suggestions remain: %q", tt.Texts())
	}
}
