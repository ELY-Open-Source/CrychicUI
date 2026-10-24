package social

import (
	"image"
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestPostCopyOverrides(t *testing.T) {
	state := PostState{}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Post(c, testPost(), &state, PostOptions{Likes: "%d appreciations"})
	}), 500, 400)
	settles(t, tt)
	if !tt.HasText("2 appreciations") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestReactionMaximum(t *testing.T) {
	reactions := []Reaction{{ID: "heart", Name: "Heart", Icon: "heart", Count: math.MaxInt}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ReactionBar(c, &reactions) }), 500, 180)
	settles(t, tt)
	if err := tt.Click("Heart"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if reactions[0].Count != math.MaxInt || reactions[0].Active {
		t.Errorf("reactions %+v", reactions)
	}
}

func TestCommentComposerKeepsFocusOnInsert(t *testing.T) {
	comments := []Comment{{ID: "social.commentComposerKey:0", Text: "First", Author: testUser()}}
	draft := "abc"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		CommentThread(c, comments, &draft, CommentThreadOptions{})
	}), 500, 400)
	settles(t, tt)
	if err := tt.Click("Write a reply"); err != nil {
		t.Fatal(err)
	}
	comments = append(comments, Comment{ID: "social.commentComposerKey:1", Text: "Second", Author: testUser()})
	tt.Frame()
	tt.Type("!")
	tt.Frame()
	if draft != "abc!" || !tt.Focused("Write a reply") {
		t.Errorf("draft %q focused %v", draft, tt.Focused("Write a reply"))
	}
}

func TestPostComposerKeepsFocusWhenMediaAppears(t *testing.T) {
	draft, added := "Hello", false
	var media ui.ImageSource
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PostComposer(c, &draft, media, &added, PostComposerOptions{})
	}), 500, 400)
	settles(t, tt)
	click(t, tt, "Post")
	media = ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2)))
	tt.Frame()
	if !tt.Focused("Post") {
		t.Error("post button lost focus after media insertion")
	}
}

func TestPollTotalOverflowRejected(t *testing.T) {
	choices := []PollChoice{{Text: "A", Votes: math.MaxInt}, {Text: "B", Votes: math.MaxInt}, {Text: "C", Votes: 2}}
	voted := -1
	defer func() {
		if recover() == nil {
			t.Error("expected poll total overflow panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		PollPost(c, "Ship?", &choices, time.Now().Add(time.Hour), &voted, PollPostOptions{})
	}, 500, 350, 1)
}

func TestPollLargeTotal(t *testing.T) {
	choices := []PollChoice{{Text: "Yes", Votes: math.MaxInt - 1}, {Text: "No", Votes: 1}}
	voted := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PollPost(c, "Ship?", &choices, c.Now().Add(time.Hour), &voted, PollPostOptions{})
	}), 500, 350)
	settles(t, tt)
	click(t, tt, "No")
	if choices[1].Votes != 2 || voted != 1 {
		t.Errorf("choices %+v voted %d", choices, voted)
	}
}
