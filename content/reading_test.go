package content

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPrevNextSteps(t *testing.T) {
	pages, current, changes := []string{"Quickstart", "Studio setup", "Book a kiln"}, 0, 0
	tt := ui.NewTester(within(480, func(c *ui.Context) {
		if PrevNext(c, pages, &current, PrevNextOptions{}).Changed() {
			changes++
		}
	}), 480, 300)
	settles(t, tt)
	if exact(tt, "Previous") {
		t.Errorf("the first page offers a previous one: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Studio setup")
	click(t, tt, "Book a kiln")
	cues(t, before, sound.Tap, sound.Tap)
	if current != 2 || changes != 2 || exact(tt, "Next") {
		t.Errorf("current %d, changes %d, texts %q", current, changes, tt.Texts())
	}
	panics(t, "current one", func(c *ui.Context) { n := 3; PrevNext(c, pages, &n, PrevNextOptions{}) })
	panics(t, "repeated", func(c *ui.Context) { PrevNext(c, []string{"A", "A"}, new(int), PrevNextOptions{}) })
}

func TestEditLinkOpens(t *testing.T) {
	url := "https://github.com/glaze-studio/docs/edit/main/kilns.md"
	tt := ui.NewTester(padded(func(c *ui.Context) { EditLink(c, url, time.Time{}, EditLinkOptions{}) }), 480, 200)
	settles(t, tt)
	if tt.HasText("Updated") {
		t.Errorf("a zero date still shows: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Edit this page")
	cues(t, before, sound.Tap)
	if !slices.Equal(tt.OpenedURLs(), []string{url}) {
		t.Errorf("opened %q", tt.OpenedURLs())
	}
	panics(t, "needs a url", func(c *ui.Context) { EditLink(c, " ", time.Time{}, EditLinkOptions{}) })
}

func TestWasHelpfulVotes(t *testing.T) {
	vote, changes := Unrated, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if WasHelpful(c, &vote, WasHelpfulOptions{}).Changed() {
			changes++
		}
	}), 480, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Yes")
	click(t, tt, "Yes")
	click(t, tt, "No")
	cues(t, before, sound.Pop, sound.Pop)
	if vote != NotHelped || changes != 2 || !exact(tt, "Thanks for the feedback.") {
		t.Errorf("vote %d, changes %d, texts %q", vote, changes, tt.Texts())
	}
	panics(t, "known vote", func(c *ui.Context) { v := Helpful(7); WasHelpful(c, &v, WasHelpfulOptions{}) })
}

func TestDocCrumbsAndEmbed(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		DocCrumbs(c, []string{"Docs", "API", "Webhooks"})
		Embed(c, "Kiln 2", EmbedOptions{Host: "glaze.studio"}, func() { ui.Text(c, "Peak 2232°F") })
	}), 480, 400)
	settles(t, tt)
	for _, s := range []string{"Webhooks", "glaze.studio", "Peak 2232°F"} {
		if !exact(tt, s) {
			t.Errorf("%q missing: %q", s, tt.Texts())
		}
	}
	panics(t, "blank part", func(c *ui.Context) { DocCrumbs(c, []string{"Docs", ""}) })
	panics(t, "needs a title and a body", func(c *ui.Context) { Embed(c, "Kiln", EmbedOptions{}, nil) })
}

func TestVersionSelectPicks(t *testing.T) {
	versions := []DocVersion{{Name: "v3.2", Latest: true}, {Name: "v3.1"}, {Name: "v3.0"}}
	current, changes := "v3.2", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if VersionSelect(c, versions, &current, VersionSelectOptions{}).Changed() {
			changes++
		}
	}), 400, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Docs version")
	settles(t, tt)
	click(t, tt, "v3.1")
	cues(t, before, sound.Open, sound.Select)
	// Opening and picking each report.
	if current != "v3.1" || changes != 2 {
		t.Errorf("current %q, changes %d", current, changes)
	}
	panics(t, "at most one latest", func(c *ui.Context) {
		s := "a"
		VersionSelect(c, []DocVersion{{Name: "a", Latest: true}, {Name: "b", Latest: true}}, &s, VersionSelectOptions{})
	})
	panics(t, "current version", func(c *ui.Context) { s := "v9"; VersionSelect(c, versions, &s, VersionSelectOptions{}) })
}

func TestVersionSelectLongListScrolls(t *testing.T) {
	var versions []DocVersion
	for i := 40; i >= 1; i-- {
		versions = append(versions, DocVersion{Name: fmt.Sprintf("v%d.0", i)})
	}
	current := "v40.0"
	tt := ui.NewTester(padded(func(c *ui.Context) { VersionSelect(c, versions, &current, VersionSelectOptions{}) }), 400, 420)
	settles(t, tt)
	click(t, tt, "Docs version")
	settles(t, tt)
	r := rectOf(t, tt, "v39.0")
	for range 8 {
		tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 400)
		tt.Advance(500 * time.Millisecond)
	}
	if l, ok := tt.Find("v1.0"); !ok || l.H == 0 || l.Y+l.H > 420 {
		t.Errorf("v1.0 out of reach: %v", l)
	}
}

func TestLinkCardOpens(t *testing.T) {
	url := "https://docs.glaze.studio/api"
	tt := ui.NewTester(padded(func(c *ui.Context) { LinkCard(c, "API reference", "Every endpoint.", url, LinkCardOptions{}) }), 480, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "API reference")
	cues(t, before, sound.Tap)
	if !exact(tt, "docs.glaze.studio") || !slices.Equal(tt.OpenedURLs(), []string{url}) {
		t.Errorf("texts %q, opened %q", tt.Texts(), tt.OpenedURLs())
	}
	panics(t, "absolute url", func(c *ui.Context) { LinkCard(c, "API", "", "/api", LinkCardOptions{}) })
}

func TestCaptionedImageChecks(t *testing.T) {
	img := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 4, 4)))
	tt := ui.NewTester(padded(func(c *ui.Context) {
		CaptionedImage(c, img, "A mug after a cone 6 firing.", CaptionedImageOptions{Credit: "Photo: Daniel"})
	}), 480, 500)
	settles(t, tt)
	if !exact(tt, "Photo: Daniel") {
		t.Errorf("credit missing: %q", tt.Texts())
	}
	panics(t, "needs an image", func(c *ui.Context) { CaptionedImage(c, nil, "Mug", CaptionedImageOptions{}) })
	panics(t, "0.25 to 4", func(c *ui.Context) { CaptionedImage(c, img, "Mug", CaptionedImageOptions{Ratio: 5}) })
}

func TestReadingTimeMinutes(t *testing.T) {
	for _, c := range []struct {
		words    int
		progress float32
		want     []string
	}{
		{1611, 0, []string{"8 min read"}},
		{0, 0, []string{"1 min read"}},
		{1611, 0.4, []string{"8 min read", "5 min left"}},
		{1611, 1, []string{"8 min read", "0 min left"}},
	} {
		tt := ui.NewTester(padded(func(ctx *ui.Context) { ReadingTime(ctx, c.words, ReadingTimeOptions{Progress: c.progress}) }), 400, 200)
		settles(t, tt)
		for _, w := range c.want {
			if !exact(tt, w) {
				t.Errorf("%d words at %v: %q missing from %q", c.words, c.progress, w, tt.Texts())
			}
		}
		if c.progress == 0 && tt.HasText("left") {
			t.Errorf("no progress still shows time left: %q", tt.Texts())
		}
	}
	panics(t, "50 to 1000 wpm", func(c *ui.Context) { ReadingTime(c, 10, ReadingTimeOptions{WPM: 10}) })
	panics(t, "progress 0 to 1", func(c *ui.Context) { ReadingTime(c, 10, ReadingTimeOptions{Progress: float32(math.NaN())}) })
}

func TestStickyNoteKeepsTilt(t *testing.T) {
	shot := func(tilt float32) []byte {
		tt := ui.NewTester(padded(func(c *ui.Context) { StickyNote(c, "Wax the feet.", StickyNoteOptions{Tilt: tilt}) }), 300, 300)
		settles(t, tt)
		return tt.Image().Pix
	}
	if bytes.Equal(shot(8), shot(-8)) {
		t.Error("tilts 8 and -8 render the same")
	}
	panics(t, "tint 0 to 3", func(c *ui.Context) { StickyNote(c, "x", StickyNoteOptions{Tint: 4}) })
	panics(t, "tilt -8 to 8", func(c *ui.Context) { StickyNote(c, "x", StickyNoteOptions{Tilt: 9}) })
}

func TestQuoteAndTimelineChecks(t *testing.T) {
	d := func(day int) time.Time { return time.Date(2026, 10, day, 9, 0, 0, 0, time.UTC) }
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Quote(c, "We stopped double-booking.", "Priya Nair", QuoteOptions{Role: "Owner"})
		StoryTimeline(c, []Beat{{Date: d(2), Title: "Thrown", Icon: "hand"}, {Date: d(9), Title: "Fired", Icon: "flame"}}, StoryTimelineOptions{})
	}), 480, 500)
	settles(t, tt)
	for _, s := range []string{"Owner", "Oct 2026", "Fired"} {
		if !exact(tt, s) {
			t.Errorf("%q missing: %q", s, tt.Texts())
		}
	}
	panics(t, "a source", func(c *ui.Context) { Quote(c, "x", "", QuoteOptions{}) })
	panics(t, "comes before", func(c *ui.Context) {
		StoryTimeline(c, []Beat{{Date: d(9), Title: "B", Icon: "flame"}, {Date: d(2), Title: "A", Icon: "hand"}}, StoryTimelineOptions{})
	})
}

func TestPrintBarCopies(t *testing.T) {
	link := "https://docs.glaze.studio/kilns"
	tt := ui.NewTester(padded(func(c *ui.Context) { PrintBar(c, link, PrintBarOptions{}) }), 480, 200)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Print")
	click(t, tt, "Copy link")
	cues(t, before, sound.Tap, sound.Tap)
	if tt.Clipboard() != link || !exact(tt, "Link copied") {
		t.Errorf("clipboard %q, texts %q", tt.Clipboard(), tt.Texts())
	}
	link = "https://docs.glaze.studio/logs"
	tt.Frame()
	if exact(tt, "Link copied") {
		t.Error("a new link still says copied")
	}
	link = "https://docs.glaze.studio/kilns"
	tt.Advance(2100 * time.Millisecond)
	if exact(tt, "Link copied") {
		t.Error("the copied note never clears")
	}
	panics(t, "needs a link", func(c *ui.Context) { PrintBar(c, "", PrintBarOptions{}) })
}
