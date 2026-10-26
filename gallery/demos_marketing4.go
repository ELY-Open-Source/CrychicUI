package gallery

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/marketing"
	"github.com/egoist/mygo/ui"
)

func demoMarketingCookieBanner(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	consent := ui.Local(root, "consent", func() marketing.Consent { return marketing.Consent{} })
	open := boolState(root, "open", false)
	root.Children(func() {
		was := *open
		if marketing.CookieBanner(c, consent, marketing.CookieBannerOptions{Open: open}).Changed() {
			if *open && !was {
				log("settings open")
			}
			log(fmt.Sprintf("decided %v, analytics %v, marketing %v", consent.Decided, consent.Analytics, consent.Marketing))
		}
	})
}

func demoMarketingContactSection(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	form := ui.Local(root, "form", func() marketing.ContactForm { return marketing.ContactForm{} })
	sent := boolState(root, "sent", false)
	root.Children(func() {
		r, _ := marketing.ContactSection(c, []marketing.ContactInfo{
			{Icon: "mail", Label: "Email", Value: "hello@glaze.example"},
			{Icon: "phone", Label: "Phone", Value: "+1 503 555 0142"},
			{Icon: "map-pin", Label: "Studio", Value: "1820 SE Division St, Portland"},
		}, form, sent, marketing.ContactSectionOptions{Text: "Questions about plans, co-ops or moving off a spreadsheet? Write to us."})
		if r.Changed() && *sent {
			log("message sent")
		} else if r.Changed() {
			log(fmt.Sprintf("form %d runes", len([]rune(form.Name+form.Email+form.Text))))
		}
	})
}

func demoMarketingSocialProof(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Children(func() {
		marketing.SocialProof(c, []marketing.Face{{Name: "Maya Chen"}, {Name: "Daniel Okafor"}, {Name: "Priya Nair"},
			{Name: "Lucas Martin"}, {Name: "Sofia Rossi"}}, 1240, marketing.SocialProofOptions{Rating: 4.9})
	})
}

func demoMarketingAppRating(c *ui.Context, log func(string)) {
	ui.Column(c).Width(440).Children(func() {
		marketing.AppRating(c, 4.8, [5]int{1820, 240, 61, 18, 27}, marketing.AppRatingOptions{})
	})
}

func demoMarketingDownload(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(520)
	opened := textState(root, "opened", "")
	root.Children(func() {
		r := marketing.Download(c, []marketing.Platform{
			{Name: "macOS", Icon: "laptop", Note: "Version 3.2, 84 MB", URL: "https://glaze.example/download/mac"},
			{Name: "Windows", Icon: "monitor", Note: "Version 3.2, 91 MB", URL: "https://glaze.example/download/windows"},
			{Name: "iPhone and iPad", Icon: "tablet-smartphone", Note: "iOS 17 or later", URL: "https://glaze.example/download/ios"},
			{Name: "Android", Icon: "smartphone", Note: "Android 11 or later", URL: "https://glaze.example/download/android"},
		}, marketing.DownloadOptions{Title: "Get Glaze on every device", Opened: opened})
		if r.Changed() {
			log("download " + *opened)
		}
	})
}

func demoMarketingPromoBanner(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		_, copy := marketing.PromoBanner(c, "Spring sale: 20% off any yearly plan until May 31.", "SPRING20", marketing.PromoBannerOptions{})
		if copy.Clicked() {
			log("copied code")
		}
	})
}

func demoMarketingLaunchCountdown(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(480)
	launch := ui.Local(root, "launch", func() time.Time { return c.Now().Add(3*24*time.Hour + 7*time.Hour + 12*time.Minute + 5*time.Second) })
	root.Children(func() {
		marketing.LaunchCountdown(c, *launch, marketing.LaunchCountdownOptions{Title: "Glaze for co-ops launches soon",
			Text: "Shared billing, member splits and one kiln calendar for every co-op studio."})
	})
}

func demoMarketingAwards(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.Awards(c, []marketing.Award{
			{Title: "Best studio software", Org: "Craft Business Awards", Year: "2025", Icon: "trophy"},
			{Title: "Top rated", Org: "Maker Tools Review", Year: "2025", Icon: "star"},
			{Title: "Small business pick", Org: "Indie Software Guide", Year: "2024"},
			{Title: "Design finalist", Org: "Product Design Week", Year: "2024", Icon: "medal"},
		}, marketing.AwardsOptions{Title: "Recognised by makers"})
	})
}

func demoMarketingPartners(c *ui.Context, log func(string)) {
	ui.Column(c).Width(520).Children(func() {
		marketing.Partners(c, []marketing.PartnerTier{
			{Name: "Gold partners", Partners: []marketing.Partner{
				{Name: "Northfire Kilns", Icon: "flame", Text: "Electric kilns with Glaze sync built in."},
				{Name: "Claymaker Supply", Icon: "package", Text: "Clay and glaze delivered to member shelves."}}},
			{Name: "Silver partners", Partners: []marketing.Partner{
				{Name: "Wheelhouse Tools", Icon: "wrench"}, {Name: "Studio Insure", Icon: "shield-check"}, {Name: "Potters Guild", Icon: "users"}}},
			{Name: "Community", Partners: []marketing.Partner{
				{Name: "Open Clay Forum", Icon: "message-square"}, {Name: "Kiln Safety Lab", Icon: "flask-conical"}}},
		}, marketing.PartnersOptions{Title: "Partners", Text: "The suppliers and groups studios already trust."})
	})
}
