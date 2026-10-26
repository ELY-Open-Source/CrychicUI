package marketing

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// within builds view in a fixed-width column.
func within(w float32, view func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) { ui.Column(c).Width(w).Children(func() { view(c) }) }
}

func TestBlogMetaFits(t *testing.T) {
	p := Post{Title: "Pricing a membership", Author: "Aleksandra Papadopoulos", Tag: "Business", Minutes: 8,
		Date: time.Date(2026, 10, 14, 9, 0, 0, 0, time.UTC)}
	long := "Monday, January 2, 2006 at 15:04:05 MST"
	opened := ""
	tt := ui.NewTester(within(320, func(c *ui.Context) {
		BlogList(c, []Post{p}, &opened, BlogListOptions{PostWords: PostWords{Date: long}})
	}), 320, 400)
	settles(t, tt)
	inside(t, "read time", rectOf(t, tt, "8 min read"), 0, 320)
	if r, ok := tt.Find(p.Date.Format(long)); ok && r.X+r.W > 320.5 {
		t.Errorf("date ends at %.1f past the card", r.X+r.W)
	}
}

func TestContactSendFits(t *testing.T) {
	var form ContactForm
	sent := false
	label := "Send your inquiry to the Glaze studio support team"
	tt := ui.NewTester(within(480, func(c *ui.Context) {
		ContactSection(c, nil, &form, &sent, ContactSectionOptions{Send: label})
	}), 480, 600)
	settles(t, tt)
	b, name := rectOf(t, tt, label), rectOf(t, tt, "Your name")
	// Name frame: 40 DIPs before the input, 13 after.
	inside(t, "Send", b, name.X-40, name.X+name.W+13)
}

func TestPromoCopiedFollowsCode(t *testing.T) {
	code := "SPRING20"
	tt := ui.NewTester(padded(func(c *ui.Context) { PromoBanner(c, "Spring sale", code, PromoBannerOptions{}) }), 520, 300)
	settles(t, tt)
	click(t, tt, "Copy code SPRING20")
	code = "SUMMER30"
	tt.Frame()
	if exact(tt, "Copied") || !exact(tt, "SUMMER30") {
		t.Errorf("a new code still says copied: %q", tt.Texts())
	}
}

func TestRoadmapVotes(t *testing.T) {
	tt := ui.NewTester(within(520, func(c *ui.Context) {
		PublicRoadmap(c, []RoadItem{{Title: "Splits", Votes: 1}, {Title: "Sync", Stage: RoadBuilding, Votes: math.MaxInt},
			{Title: "Tags", Stage: RoadShipped}}, PublicRoadmapOptions{Vote: "%d voto", Votes: "%d votos"})
	}), 520, 600)
	settles(t, tt)
	if !exact(tt, "1 voto") {
		t.Errorf("singular not overridden: %q", tt.Texts())
	}
	// A card title fills its card's inner width.
	sync := rectOf(t, tt, "Sync")
	if r, ok := tt.Find("9223372036854775807 votos"); !ok || r.X+r.W > sync.X+sync.W+0.5 {
		t.Errorf("huge vote count at %v runs past its card at %.1f", r, sync.X+sync.W)
	}
}

func TestStaggerDelays(t *testing.T) {
	steps := []HowStep{{Title: "Add", Icon: "store"}, {Title: "Invite", Icon: "users"}, {Title: "Book the kiln", Icon: "flame"}}
	d := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	rels := []Release{{Version: "3", Date: d, Changes: []Change{{Text: "a"}}}, {Version: "2", Date: d, Changes: []Change{{Text: "b"}}},
		{Version: "1", Date: d, Changes: []Change{{Text: "Card payments retry"}}}}
	for _, c := range []struct {
		what, probe string
		view        func(c *ui.Context)
	}{
		{"HowItWorks", "Book the kiln", func(c *ui.Context) { HowItWorks(c, steps, HowItWorksOptions{}) }},
		{"Changelog", "Card payments retry", func(c *ui.Context) { Changelog(c, rels, ChangelogOptions{}) }},
	} {
		st := ui.NewTester(within(520, c.view), 520, 600)
		settles(t, st)
		r := rectOf(t, st, c.probe)
		tt := ui.NewTester(within(520, c.view), 520, 600)
		tt.Frame()
		tt.Advance(20 * time.Millisecond)
		img, ink := tt.Image(), 0
		for y := int(r.Y); y < int(r.Y+r.H); y++ {
			for x := int(r.X); x < int(r.X+r.W); x++ {
				if p := img.RGBAAt(x, y); p.R < 235 || p.G < 235 || p.B < 235 {
					ink++
				}
			}
		}
		if ink > 0 {
			t.Errorf("%s: third item shows %d inked pixels at 20 ms, before its 80 ms stagger", c.what, ink)
		}
	}
}

func TestCareersReportsOnlyChanges(t *testing.T) {
	roles, opened, team, changes := testRoles(), "", "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Careers(c, roles, &opened, CareersOptions{Team: &team}).Changed() {
			changes++
		}
	}), 480, 700)
	settles(t, tt)
	if changes != 0 {
		t.Fatalf("%d changes at rest", changes)
	}
	click(t, tt, "Design")
	settles(t, tt)
	if team != "Design" || changes != 1 {
		t.Fatalf("filter gave team %q after %d changes", team, changes)
	}
	roles = roles[:1]
	tt.Frame()
	settles(t, tt)
	if team != "" || changes != 2 {
		t.Errorf("vanished team left %q after %d changes", team, changes)
	}
}

func TestLongReadAndVersionFit(t *testing.T) {
	p := Post{Title: "Pricing a membership", Author: "Maya Chen", Tag: "Business", Minutes: 8,
		Date: time.Date(2026, 10, 14, 9, 0, 0, 0, time.UTC)}
	read := "Approximate reading time: about %d minutes"
	opened := ""
	tt := ui.NewTester(within(320, func(c *ui.Context) {
		BlogList(c, []Post{p}, &opened, BlogListOptions{PostWords: PostWords{Read: read}})
	}), 320, 400)
	settles(t, tt)
	r, d := rectOf(t, tt, "Approximate reading time: about 8 minutes"), rectOf(t, tt, "Oct 14, 2026")
	if d.W <= 0 || r.X+r.W > d.X {
		t.Errorf("read time %v runs into the date %v", r, d)
	}
	rel := Release{Version: "2026.10.0-rc.1+build.12345", Title: "Co-op billing", Date: p.Date,
		Changes: []Change{{Kind: 0, Text: "Split one kiln bill across members"}}}
	tt = ui.NewTester(within(320, func(c *ui.Context) { Changelog(c, []Release{rel}, ChangelogOptions{}) }), 320, 500)
	settles(t, tt)
	v, h := rectOf(t, tt, rel.Version), rectOf(t, tt, rel.Title)
	if v.X+v.W > h.X {
		t.Errorf("version %v runs into the title %v", v, h)
	}
}
