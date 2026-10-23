package chat

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPollMessage(t *testing.T) {
	options := []string{"Pizza", "Sushi"}
	votes := []int{1, 2}
	chosen, changed := 1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if PollMessage(c, "Lunch?", options, &votes, &chosen, PollMessageOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bars do not grow")
	}
	settles(t, tt, "poll")
	if !tt.HasText("67%") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Pizza")
	if chosen != 0 || votes[0] != 2 || votes[1] != 1 || changed != 1 {
		t.Fatalf("chosen %d votes %v changed %d", chosen, votes, changed)
	}
	cues(t, before, sound.Select)
	settles(t, tt, "poll")
	if !tt.HasText("67%") || !tt.HasText("33%") || !tt.HasText("3 votes") {
		t.Errorf("texts %q", tt.Texts())
	}
	click(t, tt, "Pizza")
	if votes[0] != 2 || changed != 1 {
		t.Errorf("revote on the same option: votes %v changed %d", votes, changed)
	}
	panics(t, "votes out of step", func() {
		v := []int{1}
		render(func(c *ui.Context) { PollMessage(c, "q", options, &v, &chosen, PollMessageOptions{}) })
	})
	panics(t, "negative vote", func() {
		v := []int{-1, 0}
		render(func(c *ui.Context) { PollMessage(c, "q", options, &v, &chosen, PollMessageOptions{}) })
	})
	panics(t, "empty option", func() {
		v := []int{0, 0}
		render(func(c *ui.Context) { PollMessage(c, "q", []string{"a", ""}, &v, &chosen, PollMessageOptions{}) })
	})
	panics(t, "chosen without a vote", func() {
		v := []int{0, 1}
		n := 0
		render(func(c *ui.Context) { PollMessage(c, "q", options, &v, &n, PollMessageOptions{}) })
	})
}

func TestChannelList(t *testing.T) {
	channels := []Channel{{Name: "general", Unread: 2}, {Name: "random"}, {Name: "design", Unread: 1}}
	selected, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ChannelList(c, "Channels", channels, &selected).Changed() {
			changed++
		}
	}), 400, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("rows do not arrive")
	}
	settles(t, tt, "channels")
	before := baseline()
	click(t, tt, "random")
	if selected != 1 || changed != 1 {
		t.Fatalf("selected %d changed %d", selected, changed)
	}
	click(t, tt, "random")
	if changed != 1 {
		t.Error("reselecting reported a change")
	}
	cues(t, before, sound.Select)
	panics(t, "negative unread", func() {
		render(func(c *ui.Context) { ChannelList(c, "x", []Channel{{Name: "a", Unread: -1}}, &selected) })
	})
}

func TestMembers(t *testing.T) {
	online, offline := []string{"Ana", "Kai"}, []string{"Bo"}
	chosen, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Members(c, "Members", online, offline, &chosen, MembersOptions{}).Changed() {
			changed++
		}
	}), 400, 300)
	settles(t, tt, "members")
	if !tt.HasText("Online · 2") || !tt.HasText("Offline · 1") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Bo")
	if chosen != 2 || changed != 1 {
		t.Fatalf("chosen %d changed %d", chosen, changed)
	}
	cues(t, before, sound.Tap)
	panics(t, "no members", func() {
		render(func(c *ui.Context) { Members(c, "x", nil, nil, &chosen, MembersOptions{}) })
	})
	panics(t, "empty name", func() {
		render(func(c *ui.Context) { Members(c, "x", []string{""}, nil, &chosen, MembersOptions{}) })
	})
}

func TestOnlineBar(t *testing.T) {
	names := []string{"Ana Reyes"}
	tt := ui.NewTester(padded(func(c *ui.Context) { OnlineBar(c, "Online", names) }), 400, 200)
	tt.Frame()
	settles(t, tt, "bar")
	if !tt.HasText("AR") {
		t.Errorf("texts %q", tt.Texts())
	}
	names = append(names, "Kai Mori")
	tt.Frame()
	if !tt.FramePending() {
		t.Error("new avatar does not pop")
	}
	settles(t, tt, "bar")
	if !tt.HasText("KM") || !tt.HasText("AR") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "duplicate names", func() {
		render(func(c *ui.Context) { OnlineBar(c, "x", []string{"a", "a"}) })
	})
	panics(t, "no names", func() { render(func(c *ui.Context) { OnlineBar(c, "x", nil) }) })
}
