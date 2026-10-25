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
		"WhoToFollow":   {press("Follow Priya Nair"), "following 2"},
		"StoryRing":     {press("Daniel Okafor"), "opened s2"},
		"Hashtag":       {press("#ceramics"), "opened true"},
		"PollPost":      {press("Speckled oat"), "voted 0"},
		"Mentions":      {press("@dokafor"), "selected daniel"},
		"UserCard":      {press("Sofia Rossi"), "open false"},
		"Followers":     {press("Follow back Priya Nair"), "following 2"},
		"ReactionBar":   {press("Love"), "love 25 true"},
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

func TestSocialSecondaryLogs(t *testing.T) {
	cases := []struct{ entry, press, want string }{
		{"Post", "Comment", "action comment"},
		{"PostActions", "Save", "saved false"},
		{"ReactionBar", "Fire", "fire 10 false"},
	}
	demos := map[string]Entry{}
	for _, e := range Entries() {
		if e.Category == "social" {
			demos[e.Name] = e
		}
	}
	for _, tc := range cases {
		var logs []string
		e := demos[tc.entry]
		tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(s string) { logs = append(logs, s) }) }, 900, 700)
		tt.Advance(2 * time.Second)
		if err := tt.Click(tc.press); err != nil {
			t.Fatalf("%s: %v", tc.entry, err)
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), tc.want) {
			t.Errorf("%s %s logs %q, want %q", tc.entry, tc.press, logs, tc.want)
		}
	}
	e := demos["Feed"]
	tt := ui.NewTester(func(c *ui.Context) { e.Demo(c, func(string) {}) }, 900, 1400)
	tt.Advance(2 * time.Second)
	for _, name := range []string{"Maya Chen", "Priya Nair", "Lucas Martin"} {
		if !strings.Contains(strings.Join(tt.Texts(), "\n"), name) {
			t.Errorf("Feed lacks %s: %q", name, tt.Texts())
		}
	}
}
