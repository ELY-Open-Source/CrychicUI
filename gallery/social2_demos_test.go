package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestSocial2DemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Repost": {func(tt *ui.Tester) error {
			if err := tt.Click("Repost"); err != nil {
				return err
			}
			tt.Frame()
			return tt.Click("Quote")
		}, "action quote"},
		"Bookmark":    {press("Save"), "saved true"},
		"Community":   {press("Join"), "joined true"},
		"Invite":      {press("ana Invite"), "invited true"},
		"Question":    {press("Upvote"), "votes 1"},
		"VoteArrows":  {press("Upvote"), "score 1"},
		"SocialLinks": {press("GitHub"), "opened https://github.com"},
		"Gift":        {press("Flower"), "count 1"},
	}
	quiet := map[string]bool{"GroupCard": true, "Ranking": true, "Review": true, "Testimonial": true,
		"NotificationsFeed": true, "ProfileComplete": true, "Verified": true, "Live": true, "ShareCount": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "social" || e.ID < 1209 {
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
	if seen != 17 {
		t.Errorf("checked %d second social demos", seen)
	}
}
