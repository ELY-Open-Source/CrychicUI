package social

import (
	"image/color"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func joined(tt *ui.Tester) string { return strings.Join(tt.Texts(), "\n") }

func TestCircleListsKeepIdentity(t *testing.T) {
	users, invited := people(), map[string]bool{}
	tt := ui.NewTester(padded(func(c *ui.Context) { Invite(c, "https://kiln.co/join", users, &invited, InviteOptions{}) }), 500, 500)
	settles(t, tt)
	swapUnderPress(t, tt, "Invite Daniel Okafor", func() { users[0], users[1] = users[1], users[0] })
	if invited["maya"] {
		t.Errorf("Invite: the row that slid in was invited: %v", invited)
	}

	gifts := []GiftData{{ID: "heart", Name: "Heart", Icon: "heart"}, {ID: "coffee", Name: "Coffee", Icon: "coffee"}}
	selected := -1
	tt = ui.NewTester(padded(func(c *ui.Context) { Gift(c, &gifts, &selected, GiftOptions{}) }), 500, 300)
	settles(t, tt)
	swapUnderPress(t, tt, "Coffee", func() { gifts[0], gifts[1] = gifts[1], gifts[0] })
	for _, g := range gifts {
		if g.ID == "heart" && g.Count != 0 {
			t.Errorf("Gift: the tile that slid in was sent: %+v", gifts)
		}
	}

	links := []SocialLink{{ID: "site", Name: "Website", Icon: "globe", URL: "https://kiln.co"}, {ID: "mail", Name: "Email", Icon: "mail", URL: "mailto:a@kiln.co"}}
	opened := ""
	tt = ui.NewTester(padded(func(c *ui.Context) { SocialLinks(c, links, &opened) }), 500, 200)
	settles(t, tt)
	swapUnderPress(t, tt, "Email", func() { links[0], links[1] = links[1], links[0] })
	if opened == "https://kiln.co" {
		t.Errorf("SocialLinks opened the link that slid in")
	}
}

func TestCircleRepeatsPanic(t *testing.T) {
	u := people()
	cases := map[string]func(c *ui.Context){
		"Invite": func(c *ui.Context) {
			m := map[string]bool{}
			Invite(c, "https://a", []User{u[0], u[0]}, &m, InviteOptions{})
		},
		"Ranking": func(c *ui.Context) {
			Ranking(c, []RankedUser{{User: u[0], Score: 2}, {User: u[0], Score: 1}}, RankingOptions{})
		},
		"ProfileComplete": func(c *ui.Context) {
			ProfileComplete(c, []ProfileStep{{ID: "a", Title: "A"}, {ID: "a", Title: "B"}}, ProfileCompleteOptions{})
		},
		"SocialLinks": func(c *ui.Context) {
			l := SocialLink{ID: "a", Name: "A", Icon: "globe", URL: "https://a"}
			SocialLinks(c, []SocialLink{l, l}, new(string))
		},
		"Gift": func(c *ui.Context) {
			g := []GiftData{{ID: "a", Name: "A", Icon: "gift"}, {ID: "a", Name: "B", Icon: "gift"}}
			Gift(c, &g, new(int), GiftOptions{})
		},
		"NotificationsFeed": func(c *ui.Context) {
			n := SocialNotification{ID: "n", Kind: "like", Time: "1m", Actors: []string{"Maya Chen"}}
			NotificationsFeed(c, []SocialNotification{n, n}, NotificationsFeedOptions{})
		},
		"Question": func(c *ui.Context) {
			q := QuestionData{ID: "q", Title: "Why?", Body: "Because", Tags: []string{"glaze", "glaze"}}
			Question(c, &q, new(int), QuestionOptions{})
		},
	}
	for name, view := range cases {
		panics(t, name, view)
	}
	panics(t, "Ranking", func(c *ui.Context) {
		Ranking(c, []RankedUser{{User: u[0], Score: 1}, {User: u[1], Score: 2}}, RankingOptions{})
	})
}

func TestCircleGroupedCounts(t *testing.T) {
	u := people()
	views := map[string]func(c *ui.Context){
		"18,342 members": func(c *ui.Context) {
			Community(c, CommunityData{ID: "p", Name: "Potters", Description: "Clay", Members: 18342}, new(bool), CommunityOptions{})
		},
		"4,210 members": func(c *ui.Context) {
			GroupCard(c, GroupData{ID: "g", Name: "Glaze", Privacy: "Public", Members: 4210}, GroupCardOptions{})
		},
		"2,480 pts": func(c *ui.Context) { Ranking(c, []RankedUser{{User: u[0], Score: 2480}}, RankingOptions{}) },
		"1,204 watching": func(c *ui.Context) {
			Live(c, LiveData{ID: "l", Title: "Mugs", Streamer: "Priya Nair", Viewers: 1204, Cover: opaque(color.RGBA{200, 200, 200, 255})}, LiveOptions{})
		},
		"1 share": func(c *ui.Context) { ShareCount(c, 1, ShareCountOptions{}) },
		"1 answer": func(c *ui.Context) {
			Question(c, &QuestionData{ID: "q", Title: "Why?", Body: "Because", Answers: 1}, new(int), QuestionOptions{})
		},
		"3,214 shares": func(c *ui.Context) { ShareCount(c, 3214, ShareCountOptions{}) },
	}
	for want, view := range views {
		tt := ui.NewTester(padded(view), 500, 500)
		settles(t, tt)
		if !exact(tt, want) {
			t.Errorf("want %q in %q", want, tt.Texts())
		}
	}
}

func TestReviewMorePhotos(t *testing.T) {
	photos := make([]ui.ImageSource, 6)
	for i := range photos {
		photos[i] = opaque(color.RGBA{uint8(40 * i), 90, 120, 255})
	}
	review := ReviewData{ID: "r", Author: people()[0], Rating: 4, Text: "Lovely", Photos: photos}
	tt := ui.NewTester(padded(func(c *ui.Context) { Review(c, review) }), 500, 500)
	settles(t, tt)
	if !exact(tt, "+2") {
		t.Errorf("no overflow tile: %q", tt.Texts())
	}
	panics(t, "Review", func(c *ui.Context) { Review(c, ReviewData{ID: "r", Author: people()[0], Rating: 6, Text: "x"}) })
}

func TestVoteArrowsToggleAndLimit(t *testing.T) {
	score, vote := 1284, 0
	tt := ui.NewTester(padded(func(c *ui.Context) { VoteArrows(c, &score, &vote, VoteArrowsOptions{}) }), 500, 300)
	settles(t, tt)
	click(t, tt, "Upvote")
	settles(t, tt)
	if score != 1285 || !exact(tt, "1,285") {
		t.Fatalf("up: score %d texts %q", score, tt.Texts())
	}
	click(t, tt, "Upvote")
	if score != 1284 || vote != 0 {
		t.Errorf("second up should clear: score %d vote %d", score, vote)
	}
	score = math.MaxInt
	tt.Frame()
	if err := tt.Click("Upvote"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if score != math.MaxInt || vote != 0 {
		t.Errorf("overflowed: score %d vote %d", score, vote)
	}
}

func TestNotificationPhrases(t *testing.T) {
	notices := []SocialNotification{
		{ID: "n1", Kind: "like", Time: "2m", Actors: []string{"Priya Nair", "Daniel Okafor", "Sofia Rossi"}},
		{ID: "n2", Kind: "like", Time: "1h", Actors: []string{"Lucas Martin", "Jonas Berg"}},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { NotificationsFeed(c, notices, NotificationsFeedOptions{}) }), 500, 500)
	settles(t, tt)
	all := joined(tt)
	for _, want := range []string{"Priya Nair and 2 others liked your post", "Lucas Martin and Jonas Berg liked your post"} {
		if !strings.Contains(all, want) {
			t.Errorf("want %q in %q", want, tt.Texts())
		}
	}
	if strings.Contains(all, "New followers") {
		t.Errorf("empty group shows its header: %q", tt.Texts())
	}
	panics(t, "NotificationsFeed", func(c *ui.Context) {
		NotificationsFeed(c, notices, NotificationsFeedOptions{Like: "liked your post"})
	})
}

func TestProfileCompleteDone(t *testing.T) {
	steps := []ProfileStep{{ID: "a", Title: "Add a photo", Done: true}, {ID: "b", Title: "Write a bio", Done: true}, {ID: "c", Title: "Share a post"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ProfileComplete(c, steps, ProfileCompleteOptions{}) }), 500, 500)
	settles(t, tt)
	if !exact(tt, "66%") || exact(tt, "All set") {
		t.Fatalf("partial: %q", tt.Texts())
	}
	steps[2].Done = true
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "100%") || !exact(tt, "All set") {
		t.Errorf("done: %q", tt.Texts())
	}
}

func TestLivePulses(t *testing.T) {
	live := LiveData{ID: "l", Title: "Mugs", Streamer: "Priya Nair", Viewers: 3, Cover: opaque(color.RGBA{200, 200, 200, 255})}
	tt := ui.NewTester(padded(func(c *ui.Context) { Live(c, live, LiveOptions{}) }), 500, 500)
	tt.Advance(5 * time.Second)
	if !tt.FramePending() {
		t.Error("live dot stopped pulsing")
	}
}

func TestGiftAtMaximum(t *testing.T) {
	gifts := []GiftData{{ID: "gem", Name: "Gem", Icon: "gem", Count: math.MaxInt}}
	selected := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { Gift(c, &gifts, &selected, GiftOptions{}) }), 500, 300)
	settles(t, tt)
	if err := tt.Click("Gem"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if gifts[0].Count != math.MaxInt || selected != -1 {
		t.Errorf("gifts %+v selected %d", gifts, selected)
	}
}
