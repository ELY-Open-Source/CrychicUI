package marketing

import (
	"slices"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func testPlans() []Plan {
	return []Plan{{Name: "Starter", Features: []string{"1 kiln"}}, {Name: "Studio", Monthly: 24, Yearly: 19, Highlight: true, Badge: "Popular"},
		{Name: "Collective", Monthly: 59, Yearly: 47}}
}

func TestPricingTablePick(t *testing.T) {
	plans, picked := testPlans(), "Starter"
	yearly := false
	tt := ui.NewTester(padded(func(c *ui.Context) { PricingTable(c, plans, &picked, PricingTableOptions{Yearly: yearly}) }), 620, 600)
	settles(t, tt)
	if !exact(tt, "Free") || !exact(tt, "$24") || !exact(tt, "Current plan") {
		t.Errorf("prices: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Choose Studio")
	cues(t, before, sound.Tap)
	if picked != "Studio" {
		t.Errorf("picked %q", picked)
	}
	yearly = true
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "$19") {
		t.Errorf("yearly: %q", tt.Texts())
	}
	swapUnderPress(t, tt, "Choose Collective", func() { plans[0], plans[2] = plans[2], plans[0] })
	if picked != "Studio" {
		t.Errorf("a reorder under the press picked %q", picked)
	}
	panics(t, "picked name among them", func(c *ui.Context) { p := "Gold"; PricingTable(c, testPlans(), &p, PricingTableOptions{}) })
	panics(t, "prices of 0 to 1e9", func(c *ui.Context) {
		p := ""
		PricingTable(c, []Plan{{Name: "A", Monthly: -1}}, &p, PricingTableOptions{})
	})
}

func TestPricingCardCurrent(t *testing.T) {
	current := true
	tt := ui.NewTester(padded(func(c *ui.Context) { PricingCard(c, testPlans()[1], PricingCardOptions{Current: current}) }), 360, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Current plan")
	cues(t, before)
	current = false
	tt.Frame()
	before = baseline()
	click(t, tt, "Choose Studio")
	cues(t, before, sound.Tap)
	if !exact(tt, "Popular") || !exact(tt, "Get started") {
		t.Errorf("card: %q", tt.Texts())
	}
}

func TestBillingSwitchRolls(t *testing.T) {
	yearly := false
	monthly, perYear := 24.0, 228.0
	tt := ui.NewTester(padded(func(c *ui.Context) { BillingSwitch(c, monthly, perYear, &yearly, BillingSwitchOptions{}) }), 420, 300)
	settles(t, tt)
	if !exact(tt, "Save 20%") || !exact(tt, "$24") {
		t.Errorf("monthly: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Yearly")
	cues(t, before, sound.Toggle)
	settles(t, tt)
	if !yearly || !exact(tt, "$19") || !exact(tt, "$228 billed yearly") {
		t.Errorf("yearly: %v %q", yearly, tt.Texts())
	}
	before = baseline()
	click(t, tt, "Yearly")
	cues(t, before)
	perYear = 300
	tt.Frame()
	if exact(tt, "Save 20%") || slices.ContainsFunc(tt.Texts(), func(s string) bool { return len(s) > 4 && s[:4] == "Save" }) {
		t.Errorf("a dearer year shows a saving: %q", tt.Texts())
	}
	panics(t, "at least 0.01", func(c *ui.Context) { y := false; BillingSwitch(c, 0, 10, &y, BillingSwitchOptions{}) })
}

func TestPlanCompareMatrix(t *testing.T) {
	plans := []string{"Starter", "Studio", "Collective", "Guild", "Campus", "Network"}
	row := CompareRow{Feature: "Kilns", Cells: []Cell{{Text: "1"}, {Text: "3"}, {Included: true}, {}, {Text: "9"}, {Text: "Unlimited"}}}
	view := padded(func(c *ui.Context) { PlanCompare(c, plans, []CompareRow{row}, PlanCompareOptions{Highlight: "Studio"}) })
	tt := ui.NewTester(view, 320, 400)
	settles(t, tt)
	// Each plan column keeps its minimum; the matrix scrolls instead.
	if got := rectOf(t, tt, "Starter").W; got < 79.5 {
		t.Errorf("a narrow matrix squeezed Starter's column to %.1f", got)
	}
	if _, ok := tt.Find("Included"); !ok {
		t.Errorf("no included mark: %q", tt.Texts())
	}
	panics(t, "has 1 cells for 2 plans", func(c *ui.Context) {
		PlanCompare(c, []string{"A", "B"}, []CompareRow{{Feature: "X", Cells: []Cell{{}}}}, PlanCompareOptions{})
	})
	panics(t, "highlight among the plans", func(c *ui.Context) {
		PlanCompare(c, []string{"A"}, []CompareRow{{Feature: "X", Cells: []Cell{{}}}}, PlanCompareOptions{Highlight: "Z"})
	})
}

func TestFAQOneOpen(t *testing.T) {
	items := []QA{{"Can members book?", "Yes, they can."}, {"Any fees?", "None at all."}}
	open := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { FAQ(c, items, &open, FAQOptions{}) }), 480, 500)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Can members book?")
	cues(t, before, sound.Toggle)
	settles(t, tt)
	if open != "Can members book?" || !exact(tt, "Yes, they can.") {
		t.Errorf("open: %q %q", open, tt.Texts())
	}
	click(t, tt, "Any fees?")
	settles(t, tt)
	if open != "Any fees?" || exact(tt, "Yes, they can.") {
		t.Errorf("one open: %q %q", open, tt.Texts())
	}
	click(t, tt, "Any fees?")
	if open != "" {
		t.Errorf("close: %q", open)
	}
	panics(t, "needs an answer", func(c *ui.Context) { o := ""; FAQ(c, []QA{{"Q", ""}}, &o, FAQOptions{}) })
}

func TestCTAAction(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { CTA(c, "Fire up your studio", "Today.", "Start trial", CTAOptions{}) }), 520, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Start trial")
	cues(t, before, sound.Tap)
	panics(t, "a title and an action", func(c *ui.Context) { CTA(c, "", "", "Go", CTAOptions{}) })
}

func TestNewsletterSend(t *testing.T) {
	email, sent := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) { Newsletter(c, &email, &sent, NewsletterOptions{}) }), 480, 400)
	settles(t, tt)
	click(t, tt, "Email address")
	tt.Type("maya")
	tt.Frame()
	before := baseline()
	click(t, tt, "Subscribe")
	cues(t, before)
	click(t, tt, "Email address")
	tt.Type("@kiln.co")
	tt.Frame()
	before = baseline()
	click(t, tt, "Subscribe")
	cues(t, before, sound.Success)
	if !exact(tt, "You're on the list") {
		t.Errorf("sent: %q", tt.Texts())
	}
	before = baseline()
	click(t, tt, "Subscribe")
	cues(t, before)
	click(t, tt, "Email address")
	tt.Type("x")
	tt.Frame()
	if exact(tt, "You're on the list") {
		t.Errorf("an edit kept the sent note: %q", tt.Texts())
	}
}

func TestSiteFooterOpens(t *testing.T) {
	opened := ""
	cols := []FooterColumn{{Title: "Product", Links: []Link{{"Features", "https://a.example/f"}, {"Pricing", "https://a.example/p"}}}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SiteFooter(c, "Glaze", cols, SiteFooterOptions{Opened: &opened, Socials: []Social{{"Blog feed", "rss", "https://a.example/rss"}}})
	}), 520, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Pricing")
	click(t, tt, "Blog feed")
	cues(t, before, sound.Tap, sound.Tap)
	if opened != "https://a.example/rss" || !slices.Equal(tt.OpenedURLs(), []string{"https://a.example/p", "https://a.example/rss"}) {
		t.Errorf("opened %q %q", opened, tt.OpenedURLs())
	}
	panics(t, "needs a URL", func(c *ui.Context) {
		SiteFooter(c, "G", []FooterColumn{{Title: "P", Links: []Link{{"F", ""}}}}, SiteFooterOptions{})
	})
}

func TestSiteNavFolds(t *testing.T) {
	links := []string{"Features", "Pricing", "Customers", "Blog"}
	active := "Features"
	view := padded(func(c *ui.Context) { SiteNav(c, "Glaze", links, &active, "Sign in", SiteNavOptions{}) })
	wide := ui.NewTester(view, 560, 300)
	settles(t, wide)
	if _, ok := wide.Find("Menu"); ok {
		t.Error("a wide bar shows the menu")
	}
	before := baseline()
	click(t, wide, "Pricing")
	cues(t, before, sound.Tap)
	narrow := ui.NewTester(view, 320, 400)
	settles(t, narrow)
	if _, ok := narrow.Find("Customers"); ok {
		t.Errorf("a narrow bar shows its links: %q", narrow.Texts())
	}
	before = baseline()
	click(t, narrow, "Menu")
	cues(t, before, sound.Open)
	settles(t, narrow)
	click(t, narrow, "Blog")
	if active != "Blog" {
		t.Errorf("active %q", active)
	}
	panics(t, "active link among them", func(c *ui.Context) { a := "X"; SiteNav(c, "G", links, &a, "Go", SiteNavOptions{}) })
}

func TestAnnouncePillTap(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { AnnouncePill(c, "New", "Kiln scheduling is here") }), 420, 200)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Kiln scheduling is here")
	cues(t, before, sound.Tap)
	panics(t, "a tag and text", func(c *ui.Context) { AnnouncePill(c, "", "x") })
}

func TestWaitlistJoins(t *testing.T) {
	email, place := "", 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Waitlist(c, &email, &place, WaitlistOptions{Count: 1204, Recent: []string{"Maya Chen"}})
	}), 460, 500)
	settles(t, tt)
	if !exact(tt, "1,204 people waiting") {
		t.Errorf("count: %q", tt.Texts())
	}
	click(t, tt, "Email address")
	tt.Type("maya@kiln.co")
	tt.Frame()
	before := baseline()
	click(t, tt, "Join")
	cues(t, before, sound.Success)
	settles(t, tt)
	if place != 1205 || !exact(tt, "#1,205") || !exact(tt, "1,205 people waiting") {
		t.Errorf("joined: %d %q", place, tt.Texts())
	}
	panics(t, "place of 0 to count+1", func(c *ui.Context) { e, p := "", 9; Waitlist(c, &e, &p, WaitlistOptions{Count: 2}) })
}

func TestTeamGrid(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Team(c, []Member{{Name: "Maya Chen", Role: "Founder"}, {Name: "Jonas Berg"}}, TeamOptions{Title: "Team"})
	}), 520, 400)
	settles(t, tt)
	if !exact(tt, "Maya Chen") || !exact(tt, "Founder") {
		t.Errorf("team: %q", tt.Texts())
	}
	panics(t, "blank or repeated", func(c *ui.Context) { Team(c, []Member{{Name: "A"}, {Name: "A"}}, TeamOptions{}) })
}
