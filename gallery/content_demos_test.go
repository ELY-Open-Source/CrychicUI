package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestContentDemosLog(t *testing.T) {
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
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"DocTOC":        {press("Book a kiln"), "active 3"},
		"DocSidebar":    {press("Studio setup"), "page Studio setup"},
		"DocCodeTabs":   {press("Go"), "tab 1"},
		"Footnotes":     {press("Back to reference 2"), "back 1"},
		"DocSearch":     {press("Webhooks"), "picked Webhooks"},
		"PrevNext":      {press("Book a kiln"), "page Book a kiln"},
		"EditLink":      {press("Edit this page"), "edit opened"},
		"WasHelpful":    {press("Yes"), "vote 1"},
		"VersionSelect": {press("Docs version", "v3.1"), "version v3.1"},
		"LinkCard":      {press("Community forum"), "opened forum"},
		"PrintBar":      {press("Copy link"), "link copied"},
	}
	quiet := map[string]bool{"Article": true, "Admonition": true, "DocSteps": true, "APIRow": true, "PropTable": true, "Markdown": true,
		"Glossary": true, "DocCrumbs": true, "CaptionedImage": true, "Embed": true, "ReadingTime": true, "StickyNote": true, "Quote": true,
		"StoryTimeline": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "content" {
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
	if seen != 25 {
		t.Errorf("checked %d content demos", seen)
	}
}

func TestContentSecondaryLogs(t *testing.T) {
	for _, c := range []struct {
		name  string
		press []string
		want  string
	}{
		{"DocCodeTabs", []string{"Copy code"}, "copied"},
		{"PrintBar", []string{"Print"}, "print"},
		{"DocSidebar", []string{"Kilns"}, "section toggled"},
		{"VersionSelect", []string{"Docs version"}, "menu open"},
	} {
		for _, e := range Entries() {
			if e.Category != "content" || e.Name != c.name {
				continue
			}
			var logs []string
			tt := ui.NewTester(func(ctx *ui.Context) { e.Demo(ctx, func(s string) { logs = append(logs, s) }) }, 900, 900)
			tt.Advance(2 * time.Second)
			for _, p := range c.press {
				if err := tt.Click(p); err != nil {
					t.Fatalf("%s: %v", c.name, err)
				}
				tt.Frame()
			}
			if !strings.Contains(strings.Join(logs, "\n"), c.want) {
				t.Errorf("%s logs %q, want %q", c.name, logs, c.want)
			}
		}
	}
}
