package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/marketing"
	"github.com/egoist/mygo/ui"
)

func demoMarketingFAQ(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	open := textState(root, "open", "")
	root.Children(func() {
		items := []marketing.QA{
			{Question: "Can members book kiln space themselves?", Answer: "Yes. Members pick a firing, choose shelf space and get a text when their pieces are out."},
			{Question: "Do you take a cut of class fees?", Answer: "No. You pay the plan price; card fees go to the payment provider."},
			{Question: "Can I import my member list?", Answer: "Upload a spreadsheet of names and emails and we send each member an invite."},
			{Question: "What happens if I cancel?", Answer: "Your studio stays readable for 90 days, and you can export everything at any time."},
		}
		if marketing.FAQ(c, items, open, marketing.FAQOptions{}).Changed() {
			log(fmt.Sprintf("open %q", *open))
		}
	})
}

func demoMarketingCTA(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		_, b := marketing.CTA(c, "Fire up your studio today", "Set up classes and the kiln schedule in an afternoon.", "Start free trial",
			marketing.CTAOptions{Note: "14 days free. No card needed."})
		if b.Clicked() {
			log("start trial")
		}
	})
}

func demoMarketingNewsletter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	email, sent := textState(root, "email", ""), boolState(root, "sent", false)
	root.Children(func() {
		r, _ := marketing.Newsletter(c, email, sent, marketing.NewsletterOptions{Text: "Studio tips, glaze recipes and product news, once a month.",
			Note: "No spam. Unsubscribe any time."})
		if r.Changed() && *sent {
			log("subscribed")
		} else if r.Changed() {
			log(fmt.Sprintf("email %d runes", len([]rune(*email))))
		}
	})
}

func demoMarketingSiteFooter(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	opened := textState(root, "opened", "")
	root.Children(func() {
		r := marketing.SiteFooter(c, "Glaze", []marketing.FooterColumn{
			{Title: "Product", Links: []marketing.Link{{Label: "Features", URL: "https://glaze.example/features"},
				{Label: "Pricing", URL: "https://glaze.example/pricing"}, {Label: "Changelog", URL: "https://glaze.example/changelog"}}},
			{Title: "Company", Links: []marketing.Link{{Label: "About", URL: "https://glaze.example/about"},
				{Label: "Careers", URL: "https://glaze.example/careers"}, {Label: "Contact", URL: "https://glaze.example/contact"}}},
		}, marketing.SiteFooterOptions{Tagline: "Studio software for people who make things with clay.",
			Legal: "© 2026 Glaze Labs, Inc. Privacy · Terms",
			Socials: []marketing.Social{{Label: "Instagram", Icon: "camera", URL: "https://glaze.example/instagram"},
				{Label: "YouTube", Icon: "play", URL: "https://glaze.example/youtube"}, {Label: "Blog feed", Icon: "rss", URL: "https://glaze.example/rss"}},
			Opened: opened})
		if r.Changed() {
			log("opened " + *opened)
		}
	})
}

func demoMarketingSiteNav(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	active := textState(root, "active", "Features")
	root.Children(func() {
		r, signIn := marketing.SiteNav(c, "Glaze", []string{"Features", "Pricing", "Customers", "Blog"}, active, "Sign in", marketing.SiteNavOptions{})
		if r.Changed() {
			log("page " + *active)
		}
		if signIn.Clicked() {
			log("sign in")
		}
	})
}

func demoMarketingAnnouncePill(c *ui.Context, log func(string)) {
	ui.Column(c).Width(400).AlignItems(ui.Center).Children(func() {
		if marketing.AnnouncePill(c, "New", "Kiln scheduling is here").Clicked() {
			log("opened announcement")
		}
	})
}

func demoMarketingWaitlist(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(440)
	email := textState(root, "email", "")
	place := intState(root, "place", 0)
	root.Children(func() {
		r, _ := marketing.Waitlist(c, email, place, marketing.WaitlistOptions{Title: "Glaze for potters' co-ops",
			Text: "Shared billing for co-op studios is coming in spring. Join to get early access.", Count: 1204,
			Recent: []string{"Maya Chen", "Lucas Martin", "Jonas Berg"}})
		if r.Changed() {
			log(fmt.Sprintf("place %d", *place))
		}
	})
}

func demoMarketingTeam(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.Team(c, []marketing.Member{{Name: "Maya Chen", Role: "Founder, product"}, {Name: "Daniel Okafor", Role: "Engineering"},
			{Name: "Priya Nair", Role: "Studio partnerships"}, {Name: "Lucas Martin", Role: "Design"}, {Name: "Sofia Rossi", Role: "Customer success"},
			{Name: "Jonas Berg", Role: "Payments"}}, marketing.TeamOptions{Title: "The people behind Glaze"})
	})
}
