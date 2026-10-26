package marketing

import (
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCookieBannerChoices(t *testing.T) {
	var consent Consent
	tt := ui.NewTester(padded(func(c *ui.Context) { CookieBanner(c, &consent, CookieBannerOptions{}) }), 480, 600)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Accept all")
	cues(t, before, sound.Tap)
	if consent != (Consent{Decided: true, Analytics: true, Marketing: true}) {
		t.Fatalf("accept all gave %+v", consent)
	}
	settles(t, tt)
	click(t, tt, "Change")
	settles(t, tt)
	if consent.Decided || !exact(tt, "Always on") {
		t.Fatalf("Change did not reopen the settings: %+v %q", consent, tt.Texts())
	}
	before = baseline()
	click(t, tt, "Analytics")
	cues(t, before, sound.Toggle)
	before = baseline()
	click(t, tt, "Necessary")
	cues(t, before)
	click(t, tt, "Save choices")
	if consent != (Consent{Decided: true, Marketing: true}) {
		t.Errorf("saved %+v", consent)
	}
	panics(t, "consent", func(c *ui.Context) { CookieBanner(c, nil, CookieBannerOptions{}) })
}

func TestContactSectionSends(t *testing.T) {
	var form ContactForm
	sent := false
	info := []ContactInfo{{Icon: "mail", Label: "Email", Value: "hello@glaze.example"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ContactSection(c, info, &form, &sent, ContactSectionOptions{}) }), 520, 700)
	settles(t, tt)
	for _, f := range [][2]string{{"Your name", "Maya"}, {"How can we help?", "Hi"}, {"Email address", "maya"}} {
		click(t, tt, f[0])
		tt.Type(f[1])
		tt.Frame()
	}
	before := baseline()
	click(t, tt, "Send message")
	cues(t, before)
	click(t, tt, "Email address")
	tt.Type("@kiln.co")
	tt.Frame()
	before = baseline()
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	cues(t, before, sound.Success)
	if !sent || !exact(tt, "Thanks, we will reply within a day") {
		t.Fatalf("Enter did not send: %v %q", sent, tt.Texts())
	}
	tt.Type("x")
	tt.Frame()
	if sent {
		t.Error("an edit kept the sent state")
	}
	panics(t, "Write to us", func(c *ui.Context) { ContactSection(c, nil, nil, nil, ContactSectionOptions{Title: "Write to us"}) })
}

func TestProofAndRating(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		SocialProof(c, []Face{{Name: "Maya Chen"}, {Name: "Priya Nair"}}, 1240, SocialProofOptions{Rating: 4.9})
		AppRating(c, 4.8, [5]int{1820, 240, 61, 18, 27}, AppRatingOptions{})
		AppRating(c, 5, [5]int{1, 0, 0, 0, 0}, AppRatingOptions{Store: "Play Store"})
	}), 480, 900)
	settles(t, tt)
	for _, s := range []string{"Loved by 1,240 studios", "4.9", "4.8", "2,166 ratings", "1 rating"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	panics(t, "count of at least 2", func(c *ui.Context) {
		SocialProof(c, []Face{{Name: "A"}, {Name: "B"}}, 1, SocialProofOptions{})
	})
	panics(t, "score of 0 to 5", func(c *ui.Context) { AppRating(c, 5.1, [5]int{}, AppRatingOptions{}) })
	panics(t, "without overflow", func(c *ui.Context) { AppRating(c, 4, [5]int{math.MaxInt, 1}, AppRatingOptions{}) })
}

func TestDownloadAndPromo(t *testing.T) {
	opened := ""
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Download(c, []Platform{{Name: "macOS", Icon: "laptop", URL: "https://x.example/mac"},
			{Name: "Windows", Icon: "monitor", URL: "https://x.example/win"}}, DownloadOptions{Opened: &opened})
		PromoBanner(c, "Spring sale", "SPRING20", PromoBannerOptions{})
	}), 520, 600)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Download for Windows")
	cues(t, before, sound.Tap)
	if urls := tt.OpenedURLs(); opened != "https://x.example/win" || len(urls) != 1 || urls[0] != opened {
		t.Errorf("opened %q, urls %q", opened, urls)
	}
	before = baseline()
	click(t, tt, "Copy code SPRING20")
	cues(t, before, sound.Pop)
	if tt.Clipboard() != "SPRING20" || !exact(tt, "Copied") {
		t.Errorf("clipboard %q, texts %q", tt.Clipboard(), tt.Texts())
	}
	tt.Advance(2100 * time.Millisecond)
	settles(t, tt)
	if exact(tt, "Copied") || !exact(tt, "SPRING20") {
		t.Errorf("the copied note stayed: %q", tt.Texts())
	}
	panics(t, "Apps", func(c *ui.Context) {
		Download(c, []Platform{{Name: "macOS", Icon: "laptop"}}, DownloadOptions{Title: "Apps"})
	})
}

func TestLaunchCountdown(t *testing.T) {
	var launch time.Time
	add := 25*time.Hour + time.Minute + time.Second
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if launch.IsZero() {
			launch = c.Now().Add(add)
		}
		LaunchCountdown(c, launch, LaunchCountdownOptions{})
	}), 480, 500)
	tt.Frame()
	if !exact(tt, "1") || count(tt, "01") != 3 {
		t.Fatalf("1d 1h 1m 1s shows %q", tt.Texts())
	}
	tt.Advance(1500 * time.Millisecond)
	if count(tt, "01") != 2 || !exact(tt, "00") {
		t.Errorf("a second later shows %q", tt.Texts())
	}
	past := ui.NewTester(padded(func(c *ui.Context) { LaunchCountdown(c, c.Now().Add(-time.Second), LaunchCountdownOptions{}) }), 480, 500)
	past.Frame()
	if !exact(past, "We're live") {
		t.Errorf("past launch shows %q", past.Texts())
	}
	launch, add = time.Time{}, time.Second+time.Nanosecond
	rt := ui.NewTester(padded(func(c *ui.Context) {
		if launch.IsZero() {
			launch = c.Now().Add(add)
		}
		LaunchCountdown(c, launch, LaunchCountdownOptions{})
	}), 480, 500)
	rt.Frame()
	if !exact(rt, "02") {
		t.Errorf("1 s and 1 ns left shows %q, want whole seconds rounded up", rt.Texts())
	}
	panics(t, "launch time", func(c *ui.Context) { LaunchCountdown(c, time.Time{}, LaunchCountdownOptions{}) })
}

func TestAwardsAndPartners(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Awards(c, []Award{{Title: "Best studio software", Org: "Craft Awards", Year: "2025"}, {Title: "Top rated"}}, AwardsOptions{})
		Partners(c, []PartnerTier{{Name: "Gold", Partners: []Partner{{Name: "Northfire", Icon: "flame", Text: "Kilns"}}},
			{Name: "Silver", Partners: []Partner{{Name: "Wheelhouse", Icon: "wrench"}, {Name: "Guild", Icon: "users"}}}}, PartnersOptions{})
	}), 520, 900)
	settles(t, tt)
	for _, s := range []string{"Best studio software", "Craft Awards", "2025", "Kilns", "2"} {
		if !exact(tt, s) {
			t.Errorf("missing %q: %q", s, tt.Texts())
		}
	}
	panics(t, "repeated", func(c *ui.Context) { Awards(c, []Award{{Title: "x"}, {Title: "x"}}, AwardsOptions{}) })
	panics(t, "repeated", func(c *ui.Context) {
		Partners(c, []PartnerTier{{Name: "A", Partners: []Partner{{Name: "p", Icon: "flame"}}},
			{Name: "B", Partners: []Partner{{Name: "p", Icon: "flame"}}}}, PartnersOptions{})
	})
	panics(t, "needs partners", func(c *ui.Context) { Partners(c, []PartnerTier{{Name: "A"}}, PartnersOptions{}) })
}

func testRoles() []Opening {
	return []Opening{{Title: "Go engineer", Team: "Engineering", Place: "Remote", Kind: "Full-time"},
		{Title: "Product designer", Team: "Design", Place: "Portland", Kind: "Full-time"},
		{Title: "Data engineer", Team: "Engineering", Place: "Remote", Kind: "Contract"}}
}

func TestCareersFilterAndOpen(t *testing.T) {
	roles, opened := testRoles(), ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Careers(c, roles, &opened, CareersOptions{}) }), 480, 700)
	settles(t, tt)
	if !exact(tt, "3 open roles") {
		t.Errorf("count: %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Design")
	cues(t, before, sound.Select)
	settles(t, tt)
	if !exact(tt, "1 open role") || exact(tt, "Go engineer") {
		t.Errorf("Design filter shows %q", tt.Texts())
	}
	before = baseline()
	click(t, tt, "Design")
	cues(t, before)
	before = baseline()
	click(t, tt, "Product designer")
	cues(t, before, sound.Tap)
	if opened != "Product designer" {
		t.Errorf("opened %q", opened)
	}
	opened, roles = "", roles[:1]
	tt.Frame()
	settles(t, tt)
	if !exact(tt, "1 open role") || !exact(tt, "Go engineer") {
		t.Errorf("a vanished team kept its filter: %q", tt.Texts())
	}
	roles, opened = testRoles(), ""
	click(t, tt, "All teams")
	settles(t, tt)
	swapUnderPress(t, tt, "Product designer", func() { roles[0], roles[1] = roles[1], roles[0] })
	if opened == "Go engineer" {
		t.Error("a reorder under the press opened the role that slid in")
	}
	panics(t, "Hiring", func(c *ui.Context) {
		o := ""
		Careers(c, []Opening{{Title: "x", Team: "t", Place: "p"}}, &o, CareersOptions{Title: "Hiring"})
	})
}

func TestOptionalCopyKeepsFocus(t *testing.T) {
	more, hits := false, 0
	word := func(s string) string {
		if more {
			return s
		}
		return ""
	}
	opened := ""
	focusHolds(t, "Integrations", "Daybook", &more, &hits, func(c *ui.Context) {
		before := opened
		Integrations(c, testApps(), &opened, IntegrationsOptions{Title: word("Works with your tools")})
		if opened != before {
			hits++
			opened = ""
		}
	})
	more, hits = false, 0
	focusHolds(t, "Download", "Download for macOS", &more, &hits, func(c *ui.Context) {
		if Download(c, []Platform{{Name: "macOS", Icon: "laptop", URL: "https://x.example/mac"}},
			DownloadOptions{Title: word("Get Glaze")}).Changed() {
			hits++
		}
	})
	more, hits = false, 0
	focusHolds(t, "CaseStudy", "Read the story", &more, &hits, func(c *ui.Context) {
		var cover *ui.Bitmap
		if more {
			cover = photo()
		}
		_, r := CaseStudy(c, Case{Company: "Riverbend", Icon: "flame", Quote: "q", Name: "Priya Nair", Cover: cover,
			Metrics: []Metric{{Value: "1", Label: "x"}}}, CaseStudyOptions{})
		if r.Clicked() {
			hits++
		}
	})
}
