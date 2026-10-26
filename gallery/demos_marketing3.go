package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/marketing"
	"github.com/egoist/mygo/ui"
)

var marketingEntries2 = []Entry{
	{ID: 1411, Name: "History", Category: "marketing", Summary: "Company timeline.", Demo: demoMarketingHistory},
	{ID: 1412, Name: "HowItWorks", Category: "marketing", Summary: "Numbered steps.", Demo: demoMarketingHowItWorks},
	{ID: 1413, Name: "Integrations", Category: "marketing", Summary: "Integration tiles; cue Tap.", Demo: demoMarketingIntegrations},
	{ID: 1414, Name: "VersusTable", Category: "marketing", Summary: "Us versus them.", Demo: demoMarketingVersusTable},
	{ID: 1415, Name: "CaseStudy", Category: "marketing", Summary: "Case study card; cue Tap.", Demo: demoMarketingCaseStudy},
	{ID: 1416, Name: "BlogCard", Category: "marketing", Summary: "Post card with cover; cue Tap.", Demo: demoMarketingBlogCard},
	{ID: 1417, Name: "BlogList", Category: "marketing", Summary: "Posts list; cue Tap.", Demo: demoMarketingBlogList},
	{ID: 1418, Name: "Byline", Category: "marketing", Summary: "Author and date.", Demo: demoMarketingByline},
	{ID: 1419, Name: "Changelog", Category: "marketing", Summary: "Release notes entries.", Demo: demoMarketingChangelog},
	{ID: 1420, Name: "PublicRoadmap", Category: "marketing", Summary: "Planned, in progress, done.", Demo: demoMarketingPublicRoadmap},
	{ID: 1421, Name: "CookieBanner", Category: "marketing", Summary: "Accept and settings; cue Tap.", Demo: demoMarketingCookieBanner},
	{ID: 1422, Name: "ContactSection", Category: "marketing", Summary: "Contact info and form; cue Success.", Demo: demoMarketingContactSection},
	{ID: 1423, Name: "SocialProof", Category: "marketing", Summary: "Avatars and \"loved by\" count.", Demo: demoMarketingSocialProof},
	{ID: 1424, Name: "AppRating", Category: "marketing", Summary: "Store rating summary.", Demo: demoMarketingAppRating},
	{ID: 1425, Name: "Download", Category: "marketing", Summary: "Download for platforms; cue Tap.", Demo: demoMarketingDownload},
	{ID: 1426, Name: "PromoBanner", Category: "marketing", Summary: "Promo strip with code copy; cue Pop.", Demo: demoMarketingPromoBanner},
	{ID: 1427, Name: "LaunchCountdown", Category: "marketing", Summary: "Launch countdown.", Demo: demoMarketingLaunchCountdown},
	{ID: 1428, Name: "Awards", Category: "marketing", Summary: "Award badges row.", Demo: demoMarketingAwards},
	{ID: 1429, Name: "Partners", Category: "marketing", Summary: "Partner tiers.", Demo: demoMarketingPartners},
	{ID: 1430, Name: "Careers", Category: "marketing", Summary: "Open roles list; cue Tap.", Demo: demoMarketingCareers},
}

func demoMarketingHistory(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		marketing.History(c, []marketing.Milestone{
			{Year: "2019", Title: "Started in a shared studio in Portland", Text: "Maya kept the kiln schedule on a whiteboard until it ran out of room.", Icon: "flag"},
			{Year: "2021", Title: "Kiln scheduling ships", Text: "Members book shelf space and get a text when the firing is done.", Icon: "flame"},
			{Year: "2023", Title: "500 studios on Glaze", Text: "Classes, waitlists and member billing in one place.", Icon: "users"},
			{Year: "2025", Title: "Member billing and payouts", Icon: "credit-card"},
			{Year: "2026", Title: "Glaze for co-ops", Text: "Shared billing for studios run by their members.", Icon: "sparkles"},
		}, marketing.HistoryOptions{Title: "From one studio to a thousand"})
	})
}

func demoMarketingHowItWorks(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.HowItWorks(c, []marketing.HowStep{
			{Title: "Add your studio", Text: "Kilns, shelves and class rooms in five minutes.", Icon: "store"},
			{Title: "Invite members", Text: "They pick a plan and pay online.", Icon: "users"},
			{Title: "Book the kiln", Text: "Loads fill by cone and everyone sees the plan.", Icon: "flame"},
		}, marketing.HowItWorksOptions{Title: "How Glaze works", Text: "Three steps from whiteboard to calendar."})
	})
}

func demoMarketingIntegrations(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	opened := textState(root, "opened", "")
	root.Children(func() {
		r := marketing.Integrations(c, []marketing.Integration{
			{Name: "Tillpoint", Text: "Card payments at the studio counter.", Icon: "credit-card", Connected: true},
			{Name: "Daybook", Text: "Classes on members' calendars.", Icon: "calendar", Connected: true},
			{Name: "Postbox", Text: "Newsletters to every member.", Icon: "mail"},
			{Name: "Bookkeep", Text: "Invoices straight to your books.", Icon: "receipt"},
		}, opened, marketing.IntegrationsOptions{Title: "Works with your tools"})
		if r.Changed() {
			log("opened " + *opened)
		}
	})
}

func demoMarketingVersusTable(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Children(func() {
		marketing.VersusTable(c, "Glaze", "Spreadsheets", []marketing.Versus{
			{Feature: "Members book the kiln themselves", Us: true},
			{Feature: "Class waitlists fill cancelled seats", Us: true},
			{Feature: "Monthly member billing", Us: true},
			{Feature: "Shelf tags with QR codes", Us: true},
			{Feature: "Free to start", Us: true, Them: true},
		}, marketing.VersusTableOptions{Title: "Glaze versus the spreadsheet"})
	})
}

func demoMarketingCaseStudy(c *ui.Context, log func(string)) {
	ui.Column(c).Width(480).Children(func() {
		_, read := marketing.CaseStudy(c, marketing.Case{Company: "Riverbend Clay Co.", Icon: "flame",
			Quote: "We stopped double-booking the kiln the week we switched. Members love seeing when their pieces are ready.",
			Name:  "Priya Nair", Role: "Owner, Riverbend Clay Co.", Cover: productShots()[3],
			Metrics: []marketing.Metric{{Value: "+38%", Label: "class fill rate"}, {Value: "6 h", Label: "admin saved weekly"},
				{Value: "0", Label: "double bookings"}}}, marketing.CaseStudyOptions{})
		if read.Clicked() {
			log("read story")
		}
	})
}

// demoPosts are the studio blog's posts.
func demoPosts() []marketing.Post {
	day := func(d int) time.Time { return time.Date(2026, 10, d, 9, 0, 0, 0, time.UTC) }
	return []marketing.Post{
		{Title: "How to plan a glaze firing for forty members", Excerpt: "Cone charts, shelf maps and the one message that stops the 9 pm questions.",
			Tag: "Guides", Author: "Maya Chen", Cover: productShots()[5], Date: day(14), Minutes: 6},
		{Title: "Pricing a membership that covers the kiln", Tag: "Business", Author: "Lucas Martin", Cover: productShots()[1], Date: day(9), Minutes: 8},
		{Title: "What we learned running classes for a year", Tag: "Stories", Author: "Priya Nair", Cover: productShots()[3], Date: day(2), Minutes: 5},
		{Title: "Shelf tags are here", Tag: "Product", Author: "Daniel Okafor", Date: day(1), Minutes: 2},
	}
}

func demoMarketingBlogCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(340).Children(func() {
		if marketing.BlogCard(c, demoPosts()[0], marketing.BlogCardOptions{}).Clicked() {
			log("opened post")
		}
	})
}

func demoMarketingBlogList(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	opened := textState(root, "opened", "")
	root.Children(func() {
		if marketing.BlogList(c, demoPosts(), opened, marketing.BlogListOptions{Title: "From the studio blog"}).Changed() {
			log("opened " + *opened)
		}
	})
}

func demoMarketingByline(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).Children(func() {
		marketing.Byline(c, "Daniel Okafor", nil, time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC),
			marketing.BylineOptions{Role: "Engineering", Minutes: 8})
	})
}

func demoMarketingChangelog(c *ui.Context, log func(string)) {
	day := func(m time.Month, d int) time.Time { return time.Date(2026, m, d, 9, 0, 0, 0, time.UTC) }
	ui.Column(c).Width(520).Children(func() {
		marketing.Changelog(c, []marketing.Release{
			{Version: "2.4.0", Title: "Waitlists for classes", Date: day(10, 20), Changes: []marketing.Change{
				{Kind: marketing.ChangeNew, Text: "Class waitlists fill a seat when someone cancels."},
				{Kind: marketing.ChangeImproved, Text: "The kiln calendar loads twice as fast."},
				{Kind: marketing.ChangeFixed, Text: "Member invoices show the right studio address."}}},
			{Version: "2.3.0", Title: "Shelf tags", Date: day(9, 30), Changes: []marketing.Change{
				{Kind: marketing.ChangeNew, Text: "Print shelf tags with a QR code per member."},
				{Kind: marketing.ChangeFixed, Text: "Firing reminders no longer send twice."}}},
			{Version: "2.2.1", Date: day(9, 12), Changes: []marketing.Change{
				{Kind: marketing.ChangeFixed, Text: "Card payments retry after a dropped connection."}}},
		}, marketing.ChangelogOptions{Title: "What's new in Glaze"})
	})
}

func demoMarketingPublicRoadmap(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.PublicRoadmap(c, []marketing.RoadItem{
			{Title: "Co-op billing splits", Text: "Members share a kiln and split the bill.", Stage: marketing.RoadPlanned, Votes: 128},
			{Title: "Glaze recipe sharing", Stage: marketing.RoadPlanned, Votes: 74},
			{Title: "Kiln temperature sync", Text: "Live cone readings from connected kilns.", Stage: marketing.RoadBuilding, Votes: 211},
			{Title: "Class waitlists", Stage: marketing.RoadShipped},
			{Title: "Shelf tags", Stage: marketing.RoadShipped, Votes: 1},
		}, marketing.PublicRoadmapOptions{Title: "Roadmap", Text: "What we are building next, in the open."})
	})
}

func demoMarketingCareers(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	opened, team := textState(root, "opened", ""), textState(root, "team", "")
	root.Children(func() {
		was := *team
		r := marketing.Careers(c, []marketing.Opening{
			{Title: "Senior Go engineer", Team: "Engineering", Place: "Remote, US", Kind: "Full-time"},
			{Title: "Product designer", Team: "Design", Place: "Portland", Kind: "Full-time"},
			{Title: "Studio success lead", Team: "Customer success", Place: "Remote, EU", Kind: "Full-time"},
			{Title: "Support specialist", Team: "Customer success", Place: "Portland", Kind: "Part-time"},
			{Title: "Data engineer", Team: "Engineering", Place: "Remote", Kind: "Contract"},
		}, opened, marketing.CareersOptions{Title: "Join the Glaze team", Team: team})
		if r.Changed() && *team != was {
			log(fmt.Sprintf("team %q", *team))
		} else if r.Changed() {
			log(fmt.Sprintf("opened %s", *opened))
		}
	})
}
