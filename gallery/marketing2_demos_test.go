package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestMarketing2DemosLog(t *testing.T) {
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
	contact := func(tt *ui.Tester) error {
		for _, f := range [][2]string{{"Your name", "Maya Chen"}, {"Email address", "maya@kiln.co"}, {"How can we help?", "Do you support co-ops?"}} {
			if err := tt.Click(f[0]); err != nil {
				return err
			}
			tt.Type(f[1])
			tt.Frame()
		}
		return press("Send message")(tt)
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Integrations":   {press("Daybook"), "opened Daybook"},
		"CaseStudy":      {press("Read the story"), "read story"},
		"BlogCard":       {press("How to plan a glaze firing for forty members"), "opened post"},
		"BlogList":       {press("Pricing a membership that covers the kiln"), "opened Pricing a membership that covers the kiln"},
		"CookieBanner":   {press("Accept all"), "decided true, analytics true, marketing true"},
		"ContactSection": {contact, "message sent"},
		"Download":       {press("Download for Windows"), "download https://glaze.example/download/windows"},
		"PromoBanner":    {press("Copy code SPRING20"), "copied code"},
		"Careers":        {press("Product designer"), "opened Product designer"},
	}
	quiet := map[string]bool{"History": true, "HowItWorks": true, "VersusTable": true, "Byline": true, "Changelog": true,
		"PublicRoadmap": true, "SocialProof": true, "AppRating": true, "LaunchCountdown": true, "Awards": true, "Partners": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "marketing" || e.ID < 1411 || e.ID > 1430 {
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
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 900)
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

func TestMarketing2SecondaryLogs(t *testing.T) {
	for _, c := range []struct {
		name  string
		press []string
		want  []string
	}{
		{"CookieBanner", []string{"Settings", "Analytics"}, []string{"settings open", "decided false, analytics true, marketing false"}},
		{"Careers", []string{"Design"}, []string{`team "Design"`}},
	} {
		for _, e := range Entries() {
			if e.Category != "marketing" || e.Name != c.name {
				continue
			}
			var logs []string
			tt := ui.NewTester(func(ctx *ui.Context) { e.Demo(ctx, func(s string) { logs = append(logs, s) }) }, 900, 900)
			tt.Advance(2 * time.Second)
			for _, p := range c.press {
				if err := tt.Click(p); err != nil {
					t.Fatalf("%s: %v", c.name, err)
				}
				tt.Advance(time.Second)
			}
			for _, w := range c.want {
				if !strings.Contains(strings.Join(logs, "\n"), w) {
					t.Errorf("%s logs %q, want %q", c.name, logs, w)
				}
			}
		}
	}
}
