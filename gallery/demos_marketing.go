package gallery

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/marketing"
	"github.com/egoist/mygo/ui"
)

var marketingEntries = []Entry{
	{ID: 1391, Name: "Hero", Category: "marketing", Summary: "Headline, text, actions; rises in; cue Tap.", Demo: demoMarketingHero},
	{ID: 1392, Name: "HeroSplit", Category: "marketing", Summary: "Hero with media beside.", Demo: demoMarketingHeroSplit},
	{ID: 1393, Name: "FeatureGrid", Category: "marketing", Summary: "Features in a grid; stagger in.", Demo: demoMarketingFeatureGrid},
	{ID: 1394, Name: "FeatureRow", Category: "marketing", Summary: "Text beside media, alternating.", Demo: demoMarketingFeatureRow},
	{ID: 1395, Name: "BentoFeatures", Category: "marketing", Summary: "Features in bento tiles.", Demo: demoMarketingBento},
	{ID: 1396, Name: "LogoCloud", Category: "marketing", Summary: "Customer logos row scrolling.", Demo: demoMarketingLogoCloud},
	{ID: 1397, Name: "StatsBand", Category: "marketing", Summary: "Big numbers counting up.", Demo: demoMarketingStatsBand},
	{ID: 1398, Name: "Testimonials", Category: "marketing", Summary: "Quotes carousel; cue Swipe.", Demo: demoMarketingTestimonials},
	{ID: 1399, Name: "PricingTable", Category: "marketing", Summary: "Plans with features; cue Tap.", Demo: demoMarketingPricingTable},
	{ID: 1400, Name: "PricingCard", Category: "marketing", Summary: "One plan card with highlight.", Demo: demoMarketingPricingCard},
	{ID: 1401, Name: "BillingSwitch", Category: "marketing", Summary: "Monthly and yearly with saving; prices roll; cue Toggle.", Demo: demoMarketingBillingSwitch},
	{ID: 1402, Name: "PlanCompare", Category: "marketing", Summary: "Feature matrix.", Demo: demoMarketingPlanCompare},
	{ID: 1403, Name: "FAQ", Category: "marketing", Summary: "Questions that open; cue Toggle.", Demo: demoMarketingFAQ},
	{ID: 1404, Name: "CTA", Category: "marketing", Summary: "Call to action band; cue Tap.", Demo: demoMarketingCTA},
	{ID: 1405, Name: "Newsletter", Category: "marketing", Summary: "Email signup band; cue Success.", Demo: demoMarketingNewsletter},
	{ID: 1406, Name: "SiteFooter", Category: "marketing", Summary: "Columns, socials, legal.", Demo: demoMarketingSiteFooter},
	{ID: 1407, Name: "SiteNav", Category: "marketing", Summary: "Marketing navbar with menu; cue Open.", Demo: demoMarketingSiteNav},
	{ID: 1408, Name: "AnnouncePill", Category: "marketing", Summary: "\"New\" pill link; cue Tap.", Demo: demoMarketingAnnouncePill},
	{ID: 1409, Name: "Waitlist", Category: "marketing", Summary: "Waitlist signup with position; cue Success.", Demo: demoMarketingWaitlist},
	{ID: 1410, Name: "Team", Category: "marketing", Summary: "Team grid.", Demo: demoMarketingTeam},
}

func demoMarketingHero(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		_, start, demo := marketing.Hero(c, "Run your ceramics studio in one place",
			"Glaze books classes, schedules the kiln and bills members, so you can get back to the wheel.",
			"Start free trial", "Book a demo", marketing.HeroOptions{Eyebrow: "New: kiln scheduling"})
		if start.Clicked() {
			log("start trial")
		}
		if demo.Clicked() {
			log("book demo")
		}
	})
}

func demoMarketingHeroSplit(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		_, b := marketing.HeroSplit(c, "Every firing, on schedule", "Members book shelf space and see when their pieces come out of the kiln.",
			"See how it works", productShots()[1], marketing.HeroSplitOptions{Eyebrow: "Kiln scheduling",
				Points: []string{"Bisque and glaze loads", "Shelf space per member", "Text when pieces are ready"}})
		if b.Clicked() {
			log("see how it works")
		}
	})
}

func demoMarketingFeatureGrid(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.FeatureGrid(c, []marketing.Feature{
			{Icon: "flame", Title: "Kiln schedule", Text: "Plan bisque and glaze loads by cone and shelf."},
			{Icon: "calendar-check", Title: "Class bookings", Text: "Wheel and handbuilding classes with waitlists."},
			{Icon: "credit-card", Title: "Member billing", Text: "Monthly plans, clay fees and studio credit."},
			{Icon: "palette", Title: "Glaze library", Text: "Recipes, test tiles and photos in one place."},
		}, marketing.FeatureGridOptions{Title: "Built for clay studios", Text: "Everything a shared studio runs on, without the spreadsheet."})
	})
}

func demoMarketingFeatureRow(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.FeatureRow(c, []marketing.Showcase{
			{Eyebrow: "Classes", Title: "Fill every wheel", Text: "Open a class, share the link and let the waitlist fill gaps.",
				Points: []string{"Automatic waitlist", "Reminders the day before"}, Media: productShots()[5]},
			{Eyebrow: "Members", Title: "Know who is firing what", Text: "Each member sees their shelf, their pieces and their balance.",
				Points: []string{"Shelf tags by member", "Clay and firing fees"}, Media: productShots()[3]},
		}, marketing.FeatureRowOptions{})
	})
}

func demoMarketingBento(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.BentoFeatures(c, []marketing.Bento{
			{Icon: "flame", Title: "Kiln log", Text: "Every firing with cone, load and notes.", Wide: true, Media: productShots()[4]},
			{Icon: "users", Title: "Member portal", Text: "Bookings, shelves and invoices."},
			{Icon: "chart-line", Title: "Studio reports", Text: "Classes, revenue and clay use."},
			{Icon: "palette", Title: "Glaze tests", Text: "Photo every test tile and share recipes with the studio.", Wide: true},
		}, marketing.BentoFeaturesOptions{Title: "One place for the whole studio"})
	})
}

func demoMarketingLogoCloud(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.LogoCloud(c, []marketing.Logo{{Name: "Clayworks", Icon: "flame"}, {Name: "Greenware", Icon: "leaf"},
			{Name: "Ridgeline", Icon: "mountain"}, {Name: "Tidal Clay", Icon: "waves"}, {Name: "Sunfire", Icon: "sun"},
			{Name: "Featherstone", Icon: "feather"}}, marketing.LogoCloudOptions{})
	})
}

func demoMarketingStatsBand(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.StatsBand(c, []marketing.Stat{{Value: 1240, Suffix: "+", Label: "Studios"},
			{Value: 3.2, Decimals: 1, Suffix: "M", Label: "Firings logged"}, {Value: 4.9, Decimals: 1, Suffix: "/5", Label: "Average rating"}},
			marketing.StatsBandOptions{Title: "Studios fire with Glaze"})
	})
}

func demoMarketingTestimonials(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	at := intState(root, "at", 0)
	root.Children(func() {
		quotes := []marketing.Quote{
			{Text: "We stopped double-booking the kiln the week we switched. Members love seeing when their pieces are done.",
				Name: "Priya Nair", Role: "Owner, Kiln & Co.", Rating: 5},
			{Text: "Class waitlists fill themselves now. I spend Sunday nights throwing instead of answering emails.",
				Name: "Daniel Okafor", Role: "Teacher, Ridgeline Clay", Rating: 5},
			{Text: "The glaze library alone was worth it. Every test tile has a photo and a recipe.",
				Name: "Sofia Rossi", Role: "Glaze lead, Tidal Clay", Rating: 4},
		}
		if marketing.Testimonials(c, quotes, at, marketing.TestimonialsOptions{Title: "Loved by studio owners"}).Changed() {
			log(fmt.Sprintf("quote %d", *at+1))
		}
	})
}

var demoPlans = []marketing.Plan{
	{Name: "Starter", Text: "For a home studio.", Features: []string{"1 kiln", "20 members", "Class bookings"}},
	{Name: "Studio", Text: "For a busy shared studio.", Monthly: 24, Yearly: 19, Highlight: true, Badge: "Most popular",
		Features: []string{"3 kilns", "150 members", "Member billing", "Glaze library"}},
	{Name: "Collective", Text: "For studios with several sites.", Monthly: 59, Yearly: 47,
		Features: []string{"Unlimited kilns", "Unlimited members", "Reports", "Priority help"}},
}

func demoMarketingPricingTable(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	picked := textState(root, "plan", "Starter")
	root.Children(func() {
		if marketing.PricingTable(c, demoPlans, picked, marketing.PricingTableOptions{}).Changed() {
			log("plan " + *picked)
		}
	})
}

func demoMarketingPricingCard(c *ui.Context, log func(string)) {
	ui.Column(c).Width(340).Children(func() {
		_, choose := marketing.PricingCard(c, demoPlans[1], marketing.PricingCardOptions{Yearly: true})
		if choose.Clicked() {
			log("chose Studio")
		}
	})
}

func demoMarketingBillingSwitch(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(400)
	yearly := boolState(root, "yearly", false)
	root.Children(func() {
		if marketing.BillingSwitch(c, 24, 228, yearly, marketing.BillingSwitchOptions{}).Changed() {
			log(fmt.Sprintf("yearly %v", *yearly))
		}
	})
}

func demoMarketingPlanCompare(c *ui.Context, log func(string)) {
	yes, no := marketing.Cell{Included: true}, marketing.Cell{}
	ui.Column(c).Width(520).Children(func() {
		marketing.PlanCompare(c, []string{"Starter", "Studio", "Collective"}, []marketing.CompareRow{
			{Feature: "Kilns", Cells: []marketing.Cell{{Text: "1"}, {Text: "3"}, {Text: "Unlimited"}}},
			{Feature: "Members", Cells: []marketing.Cell{{Text: "20"}, {Text: "150"}, {Text: "Unlimited"}}},
			{Feature: "Class bookings", Cells: []marketing.Cell{yes, yes, yes}},
			{Feature: "Member billing", Cells: []marketing.Cell{no, yes, yes}},
			{Feature: "Glaze library", Cells: []marketing.Cell{no, yes, yes}},
			{Feature: "Reports", Cells: []marketing.Cell{no, no, yes}},
		}, marketing.PlanCompareOptions{Highlight: "Studio"})
	})
}
