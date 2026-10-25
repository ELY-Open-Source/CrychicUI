package social

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestRankingBarsFollowScores(t *testing.T) {
	u := people()
	rows := func(second int) []RankedUser {
		return []RankedUser{{User: u[0], Score: 100}, {User: u[1], Score: second}}
	}
	live := rows(50)
	tt := ui.NewTester(padded(func(c *ui.Context) { Ranking(c, live, RankingOptions{}) }), 500, 400)
	settles(t, tt)
	live = rows(25)
	tt.Frame()
	settles(t, tt)
	want := ui.NewTester(padded(func(c *ui.Context) { Ranking(c, rows(25), RankingOptions{}) }), 500, 400)
	settles(t, want)
	if !bytes.Equal(tt.Image().Pix, want.Image().Pix) {
		t.Error("second bar kept its old share after the score fell")
	}
}

func TestReviewOverflowTileVisible(t *testing.T) {
	photos := make([]ui.ImageSource, 6)
	for i := range photos {
		photos[i] = opaque(color.RGBA{uint8(40 * i), 90, 120, 255})
	}
	review := ReviewData{ID: "r", Author: people()[0], Rating: 4, Text: "Lovely", Photos: photos}
	tt := ui.NewTester(padded(func(c *ui.Context) { Review(c, review) }), 500, 500)
	settles(t, tt)
	r, ok := tt.Find("+2")
	if !ok || r.W == 0 || r.H == 0 {
		t.Errorf("+2 not visible: %v %+v", ok, r)
	}
}

func TestCommunityJoinKeepsFocus(t *testing.T) {
	community := CommunityData{ID: "k", Name: "Kiln", Description: "Potters", Members: 12}
	joined := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Community(c, community, &joined, CommunityOptions{}) }), 500, 500)
	settles(t, tt)
	click(t, tt, "Join")
	community.Faces = []string{"Priya Nair", "Daniel Okafor"}
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	tt.Frame()
	if joined {
		t.Error("Space did not leave: Join lost focus when faces appeared")
	}
}

func TestGiftCountsStayInTiles(t *testing.T) {
	gifts := []GiftData{{ID: "a", Name: "Heart", Icon: "heart", Count: math.MaxInt},
		{ID: "b", Name: "Coffee", Icon: "coffee", Count: math.MaxInt - 1}, {ID: "c", Name: "Gem", Icon: "gem", Count: math.MaxInt - 2}}
	selected := -1
	tt := ui.NewTester(narrow(func(c *ui.Context) { Gift(c, &gifts, &selected, GiftOptions{}) }), 500, 400)
	settles(t, tt)
	for _, g := range gifts {
		tile, _ := tt.Find(g.Name)
		count, ok := tt.Find(fmt.Sprintf("×%s", grouped(g.Count)))
		if !ok || count.X < tile.X-0.5 || count.X+count.W > tile.X+tile.W+0.5 {
			t.Errorf("%s count %+v leaves tile %+v", g.Name, count, tile)
		}
	}
}

func TestProfileCompleteExactPercent(t *testing.T) {
	steps := make([]ProfileStep, 100)
	for i := range steps {
		steps[i] = ProfileStep{ID: fmt.Sprint(i), Title: fmt.Sprint("Step ", i), Done: i < 53}
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { ProfileComplete(c, steps, ProfileCompleteOptions{}) }), 500, 600)
	settles(t, tt)
	if !exact(tt, "53%") || !exact(tt, "53/100 complete") {
		t.Errorf("texts %q", tt.Texts()[:4])
	}
}

func TestCirclePanicsNameCaller(t *testing.T) {
	notices := []SocialNotification{{ID: "n1", Kind: "like", Time: "2m", Actors: []string{"Priya Nair"}}}
	panics(t, "one %s", func(c *ui.Context) {
		NotificationsFeed(c, notices, NotificationsFeedOptions{Like: "%s and %s liked your post"})
	})
	panics(t, `Testimonial "Priya Nair"`, func(c *ui.Context) { Testimonial(c, "", "Priya Nair", "Potter") })
	for _, vote := range []*int{new(2), nil} {
		panics(t, `Question "Why glaze?"`, func(c *ui.Context) {
			Question(c, &QuestionData{ID: "q", Title: "Why glaze?", Body: "Asking"}, vote, QuestionOptions{})
		})
	}
}

// grouped writes n with commas, as the components do.
func grouped(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
