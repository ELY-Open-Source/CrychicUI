package gallery

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestSocialDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Post":          {press("Like"), "liked true"},
		"Feed":          {press("New posts"), "fresh false"},
		"PostComposer":  {press("Post"), "sent"},
		"PostActions":   {press("Like"), "liked true"},
		"CommentThread": {press("Reply"), "replied"},
		"ProfileHeader": {press("Follow"), "following true"},
		"WhoToFollow":   {press("ana Follow"), "following true"},
		"StoryRing":     {press("Ana"), "opened s1"},
		"Hashtag":       {press("#mygo"), "opened true"},
		"PollPost":      {press("Yes"), "voted 0"},
		"Mentions":      {press("@ana"), "selected ana"},
		"UserCard":      {press("Ana"), "open true"},
		"Followers":     {press("ana Follow back"), "following true"},
		"ReactionBar":   {press("Heart"), "count 3"},
	}
	quiet := map[string]bool{"ProfileStats": true, "Trending": true, "MediaPost": true, "ProfileBadges": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "social" || e.ID >= 1209 {
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
	if seen != 18 {
		t.Errorf("checked %d social demos", seen)
	}
}
