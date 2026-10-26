package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestMarketingDemosLog(t *testing.T) {
	press := func(names ...string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			for _, name := range names {
				if err := tt.Click(name); err != nil {
					return err
				}
				tt.Frame()
			}
			return nil
		}
	}
	signup := func(button string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error {
			if err := tt.Click("Email address"); err != nil {
				return err
			}
			tt.Type("maya@kiln.co")
			tt.Frame()
			return press(button)(tt)
		}
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Hero":          {press("Start free trial"), "start trial"},
		"HeroSplit":     {press("See how it works"), "see how it works"},
		"Testimonials":  {press("Next quote"), "quote 2"},
		"PricingTable":  {press("Choose Studio"), "plan Studio"},
		"PricingCard":   {press("Choose Studio"), "chose Studio"},
		"BillingSwitch": {press("Yearly"), "yearly true"},
		"FAQ":           {press("Do you take a cut of class fees?"), `open "Do you take a cut of class fees?"`},
		"CTA":           {press("Start free trial"), "start trial"},
		"Newsletter":    {signup("Subscribe"), "email 12 runes\nsubscribed"},
		"SiteFooter":    {press("Pricing"), "opened https://glaze.example/pricing"},
		"SiteNav":       {press("Pricing"), "page Pricing"},
		"AnnouncePill":  {press("Kiln scheduling is here"), "opened announcement"},
		"Waitlist":      {signup("Join"), "place 1205"},
	}
	quiet := map[string]bool{"FeatureGrid": true, "FeatureRow": true, "BentoFeatures": true, "LogoCloud": true, "StatsBand": true,
		"PlanCompare": true, "Team": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "marketing" || e.ID > 1410 {
			continue
		}
		seen++
		if quiet[e.Name] {
			continue
		}
		act, ok := acts[e.Name]
		if !ok {
			t.Errorf("%s has no logged action", e.Name)
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(2 * time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 20 {
		t.Errorf("checked %d marketing demos", seen)
	}
}

func TestNewsletterDemoEnter(t *testing.T) {
	for _, e := range Entries() {
		if e.Category != "marketing" || e.Name != "Newsletter" {
			continue
		}
		var logs []string
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(2 * time.Second)
		if err := tt.Click("Email address"); err != nil {
			t.Fatal(err)
		}
		tt.Type("maya@kiln.co")
		tt.Frame()
		tt.Key(0, ui.KeyEnter)
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), "subscribed") {
			t.Errorf("Enter signup logs %q", logs)
		}
		return
	}
	t.Error("no Newsletter demo")
}
