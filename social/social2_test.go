package social

import (
	"image"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestRepost(t *testing.T) {
	action := ""
	changes := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Repost(c, &action, RepostOptions{}).Changed() {
			changes++
		}
	}), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Repost")
	if changes != 0 {
		t.Errorf("menu open reported %d actions", changes)
	}
	click(t, tt, "Quote")
	click(t, tt, "Repost")
	if action != "quote" || changes != 1 {
		t.Errorf("action %q changes %d", action, changes)
	}
	cues(t, before, sound.Pop, sound.Pop, sound.Pop)
}

func TestBookmark(t *testing.T) {
	saved := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Bookmark(c, &saved, BookmarkOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Save")
	if !saved {
		t.Error("not saved")
	}
	cues(t, before, sound.Toggle)
}

func TestCommunity(t *testing.T) {
	joined := false
	community := CommunityData{ID: "ui", Name: "UI", Description: "Build", Members: 42}
	tt := ui.NewTester(padded(func(c *ui.Context) { Community(c, community, &joined, CommunityOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Join")
	if !joined || !tt.HasText("42 members") {
		t.Errorf("joined %v texts %q", joined, tt.Texts())
	}
	cues(t, before, sound.Toggle)
}

func TestGroupCard(t *testing.T) {
	group := GroupData{ID: "ui", Name: "UI", Privacy: "Public", Members: 42}
	tt := ui.NewTester(padded(func(c *ui.Context) { GroupCard(c, group, GroupCardOptions{}) }), 500, 220)
	settles(t, tt)
	if !tt.HasText("Public") || !tt.HasText("42 members") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestInvite(t *testing.T) {
	invited := map[string]bool{}
	link := "https://example.test/i"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Invite(c, link, []User{testUser()}, &invited, InviteOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Copy link")
	if tt.Clipboard() != link {
		t.Errorf("clipboard %q", tt.Clipboard())
	}
	link = "https://example.test/new"
	tt.Frame()
	if !tt.HasText("Copy link") || tt.HasText("Copied") {
		t.Errorf("stale copy state: %q", tt.Texts())
	}
	click(t, tt, "ana Invite")
	if !invited["ana"] {
		t.Errorf("invited %v", invited)
	}
	cues(t, before, sound.Pop, sound.Pop)
}

func TestRanking(t *testing.T) {
	rows := []RankedUser{{User: testUser(), Score: 100}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Ranking(c, rows, RankingOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("100 pts") || !tt.HasText("Ana") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestReview(t *testing.T) {
	img := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2)))
	review := ReviewData{ID: "r1", Author: testUser(), Rating: 5, Text: "Great work", Photos: []ui.ImageSource{img}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Review(c, review) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("Great work") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestTestimonial(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Testimonial(c, "Native feel", "Ana", "Designer") }), 500, 250)
	settles(t, tt)
	if !tt.HasText("Native feel") || !tt.HasText("Designer") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestQuestion(t *testing.T) {
	question := QuestionData{ID: "q1", Title: "Native?", Body: "How?", Answers: 2}
	vote := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { Question(c, &question, &vote, QuestionOptions{}) }), 500, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Upvote")
	if question.Votes != 1 || vote != 1 || !tt.HasText("2 answers") {
		t.Errorf("question %+v vote %d texts %q", question, vote, tt.Texts())
	}
	cues(t, before, sound.Tick)
}

func TestQuestionVoteCopy(t *testing.T) {
	question := QuestionData{ID: "q1", Title: "Native?", Body: "How?"}
	vote := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Question(c, &question, &vote, QuestionOptions{Up: "赞成", Down: "反对"})
	}), 500, 350)
	settles(t, tt)
	click(t, tt, "赞成")
	if question.Votes != 1 || !tt.HasText("反对") {
		t.Errorf("votes %d texts %q", question.Votes, tt.Texts())
	}
}

func TestVoteArrows(t *testing.T) {
	score, vote := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) { VoteArrows(c, &score, &vote, VoteArrowsOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Upvote")
	click(t, tt, "Downvote")
	if score != -1 || vote != -1 {
		t.Errorf("score %d vote %d", score, vote)
	}
	cues(t, before, sound.Tick, sound.Tick)
}

func TestNotificationsFeed(t *testing.T) {
	notices := []SocialNotification{{ID: "n1", Kind: "like", Actors: []string{"Ana"}, Time: "now"}, {ID: "n2", Kind: "follow", Actors: []string{"Ben"}, Time: "1m"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { NotificationsFeed(c, notices, NotificationsFeedOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("Likes") || !tt.HasText("New followers") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestProfileComplete(t *testing.T) {
	steps := []ProfileStep{{ID: "photo", Title: "Add photo", Done: true}, {ID: "bio", Title: "Add bio"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ProfileComplete(c, steps, ProfileCompleteOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("1/2 complete") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestVerified(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Verified(c, "Verified", VerifiedOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Verified") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSocialLinks(t *testing.T) {
	opened := ""
	links := []SocialLink{{ID: "gh", Name: "GitHub", Icon: "link", URL: "https://github.com"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { SocialLinks(c, links, &opened) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "GitHub")
	if opened != "https://github.com" {
		t.Errorf("opened %q", opened)
	}
	cues(t, before, sound.Tap)
}

func TestLive(t *testing.T) {
	img := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2)))
	live := LiveData{ID: "l1", Title: "UI review", Streamer: "Ana", Viewers: 42, Cover: img}
	tt := ui.NewTester(padded(func(c *ui.Context) { Live(c, live, LiveOptions{}) }), 500, 400)
	settles(t, tt)
	if !tt.HasText("42 viewers") || !tt.HasText("LIVE") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestGift(t *testing.T) {
	gifts := []GiftData{{ID: "flower", Name: "Flower", Icon: "flower"}}
	selected := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Gift(c, &gifts, &selected, GiftOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Flower")
	if selected != 0 || gifts[0].Count != 1 {
		t.Errorf("gifts %+v selected %d", gifts, selected)
	}
	cues(t, before, sound.Pop)
}

func TestShareCount(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { ShareCount(c, 42, ShareCountOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("42 shares") {
		t.Errorf("texts %q", tt.Texts())
	}
}
