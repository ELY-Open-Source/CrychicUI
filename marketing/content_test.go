package marketing

import (
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func count(tt *ui.Tester, s string) int {
	n := 0
	for _, x := range tt.Texts() {
		if x == s {
			n++
		}
	}
	return n
}

func TestHistoryAndSteps(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		History(c, []Milestone{{Year: "2019", Title: "Founded"}, {Year: "2023", Title: "500 studios", Icon: "users"}}, HistoryOptions{})
		HowItWorks(c, []HowStep{{Title: "Add", Icon: "store"}, {Title: "Invite", Icon: "users"}, {Title: "Book", Icon: "flame"}}, HowItWorksOptions{})
	}), 520, 700)
	settles(t, tt)
	for _, s := range []string{"2019", "Founded", "1", "2", "3", "Book"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	panics(t, "Kiln log", func(c *ui.Context) { History(c, []Milestone{{Title: "x"}}, HistoryOptions{Title: "Kiln log"}) })
	panics(t, "2 to 6 steps", func(c *ui.Context) { HowItWorks(c, []HowStep{{Title: "Add", Icon: "store"}}, HowItWorksOptions{}) })
	panics(t, "repeated", func(c *ui.Context) {
		HowItWorks(c, []HowStep{{Title: "Add", Icon: "store"}, {Title: "Add", Icon: "users"}}, HowItWorksOptions{})
	})
}

func testApps() []Integration {
	return []Integration{{Name: "Tillpoint", Icon: "credit-card", Connected: true}, {Name: "Daybook", Icon: "calendar"},
		{Name: "Postbox", Icon: "mail"}}
}

func TestIntegrationsOpen(t *testing.T) {
	apps, opened, changes := testApps(), "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Integrations(c, apps, &opened, IntegrationsOptions{}).Changed() {
			changes++
		}
	}), 520, 600)
	settles(t, tt)
	if !exact(tt, "Connected") {
		t.Errorf("no connected tag: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Daybook")
	cues(t, before, sound.Tap)
	if opened != "Daybook" || changes != 1 {
		t.Errorf("opened %q after %d changes", opened, changes)
	}
	before = baseline()
	click(t, tt, "Daybook")
	cues(t, before, sound.Tap)
	if changes != 1 {
		t.Errorf("reopening the same app reported %d changes", changes)
	}
	opened = ""
	settles(t, tt)
	swapUnderPress(t, tt, "Daybook", func() { apps[0], apps[1] = apps[1], apps[0] })
	if opened == "Tillpoint" {
		t.Error("a reorder under the press opened the app that slid in")
	}
	panics(t, "opened app", func(c *ui.Context) { o := "Nope"; Integrations(c, testApps(), &o, IntegrationsOptions{}) })
}

func TestVersusAndCase(t *testing.T) {
	read := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		VersusTable(c, "Glaze", "Sheets", []Versus{{Feature: "Kiln booking", Us: true}, {Feature: "Free", Us: true, Them: true}}, VersusTableOptions{})
		_, r := CaseStudy(c, Case{Company: "Riverbend", Icon: "flame", Quote: "No double bookings.", Name: "Priya Nair",
			Metrics: []Metric{{Value: "+38%", Label: "fill rate"}}}, CaseStudyOptions{})
		if r.Clicked() {
			read++
		}
	}), 520, 800)
	settles(t, tt)
	for _, s := range []string{"Glaze", "Sheets", "Kiln booking", "+38%", "fill rate"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	before := baseline()
	click(t, tt, "Read the story")
	cues(t, before, sound.Tap)
	if read != 1 {
		t.Errorf("read %d times", read)
	}
	panics(t, "two different names", func(c *ui.Context) { VersusTable(c, "A", "A", []Versus{{Feature: "x"}}, VersusTableOptions{}) })
	five := make([]Metric, 5)
	for i := range five {
		five[i] = Metric{Value: "1", Label: "x"}
	}
	panics(t, "Riverbend", func(c *ui.Context) {
		CaseStudy(c, Case{Company: "Riverbend", Icon: "flame", Quote: "q", Name: "n", Metrics: five}, CaseStudyOptions{})
	})
}

func testPosts() []Post {
	d := time.Date(2026, 10, 14, 9, 0, 0, 0, time.UTC)
	return []Post{{Title: "Plan a firing", Author: "Maya Chen", Date: d, Minutes: 6, Tag: "Guides"},
		{Title: "Pricing", Author: "Lucas Martin", Date: d, Minutes: 1}, {Title: "Shelf tags", Author: "Daniel Okafor", Date: d}}
}

func TestBlogCardAndByline(t *testing.T) {
	opens := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if BlogCard(c, testPosts()[0], BlogCardOptions{}).Clicked() {
			opens++
		}
		Byline(c, "Daniel Okafor", nil, time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), BylineOptions{Role: "Engineering", Minutes: 8})
	}), 400, 700)
	settles(t, tt)
	for _, s := range []string{"6 min read", "Oct 14, 2026", "8 min read", "Oct 1, 2026", "Engineering"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	before := baseline()
	click(t, tt, "Plan a firing")
	cues(t, before, sound.Tap)
	if opens != 1 {
		t.Errorf("opened %d times", opens)
	}
	panics(t, "Plan a firing", func(c *ui.Context) { BlogCard(c, Post{Title: "Plan a firing", Author: "M"}, BlogCardOptions{}) })
	panics(t, "Daniel", func(c *ui.Context) { Byline(c, "Daniel", nil, time.Time{}, BylineOptions{}) })
}

func TestBlogListOpens(t *testing.T) {
	posts, opened := testPosts(), ""
	tt := ui.NewTester(padded(func(c *ui.Context) { BlogList(c, posts, &opened, BlogListOptions{}) }), 480, 600)
	settles(t, tt)
	if exact(tt, "0 min read") {
		t.Errorf("zero minutes showed a read time: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Pricing")
	cues(t, before, sound.Tap)
	if opened != "Pricing" {
		t.Errorf("opened %q", opened)
	}
	opened = ""
	settles(t, tt)
	swapUnderPress(t, tt, "Pricing", func() { posts[0], posts[1] = posts[1], posts[0] })
	if opened == "Plan a firing" {
		t.Error("a reorder under the press opened the post that slid in")
	}
	none := ""
	et := ui.NewTester(padded(func(c *ui.Context) { BlogList(c, nil, &none, BlogListOptions{}) }), 480, 400)
	settles(t, et)
	if !exact(et, "No posts yet") {
		t.Errorf("no empty state: %q", et.Texts())
	}
	panics(t, "repeated", func(c *ui.Context) { o := ""; BlogList(c, append(testPosts(), testPosts()[0]), &o, BlogListOptions{}) })
}

func TestChangelogAndRoadmap(t *testing.T) {
	d := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Changelog(c, []Release{{Version: "2.4.0", Date: d, Changes: []Change{{Kind: ChangeNew, Text: "Waitlists"}, {Kind: ChangeFixed, Text: "Invoices"}}},
			{Version: "2.3.0", Date: d, Changes: []Change{{Kind: ChangeImproved, Text: "Faster"}}}}, ChangelogOptions{})
		PublicRoadmap(c, []RoadItem{{Title: "Splits", Votes: 1}, {Title: "Sync", Stage: RoadBuilding, Votes: 211}}, PublicRoadmapOptions{})
	}), 520, 1200)
	settles(t, tt)
	for _, s := range []string{"New", "Fixed", "Improved", "Latest", "1 vote", "211 votes", "Nothing here yet"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	if n := count(tt, "Latest"); n != 1 {
		t.Errorf("%d releases marked latest", n)
	}
	panics(t, "known kind", func(c *ui.Context) {
		Changelog(c, []Release{{Version: "1", Date: d, Changes: []Change{{Kind: 3, Text: "x"}}}}, ChangelogOptions{})
	})
	panics(t, "known stage", func(c *ui.Context) { PublicRoadmap(c, []RoadItem{{Title: "x", Stage: 3}}, PublicRoadmapOptions{}) })
	if slices.Contains(tt.Texts(), "0 votes") {
		t.Error("zero votes showed")
	}
}
