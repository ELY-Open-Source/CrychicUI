package chat

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func sample() []Message {
	return []Message{
		{Sender: "Ana", Text: "Morning!", At: now},
		{Sender: "Ana", Text: "Lunch at noon?", At: now.Add(time.Minute)},
		{Sender: "me", Text: "Sure", Mine: true, At: now.Add(2 * time.Minute)},
	}
}

func TestMessageList(t *testing.T) {
	list := sample()
	tt := ui.NewTester(padded(func(c *ui.Context) { MessageList(c, "Chat", list, MessageListOptions{}) }), 500, 400)
	settles(t, tt, "list")
	if n := count(tt, "Ana"); n != 1 {
		t.Errorf("sender shown %d times, want once per group", n)
	}
	a, b := find(t, tt, "Morning!"), find(t, tt, "Sure")
	if a.X >= b.X {
		t.Errorf("theirs at %v, mine at %v: mine should sit right", a.X, b.X)
	}
	before := baseline()
	list = append(list, Message{Sender: "Ana", Text: "Great", At: now.Add(20 * time.Minute)})
	tt.Frame()
	if !tt.FramePending() {
		t.Error("new row does not slide up")
	}
	settles(t, tt, "list")
	if n := count(tt, "Ana"); n != 2 {
		t.Errorf("sender shown %d times after a new group, want 2", n)
	}
	cues(t, before)
	quiet := baseline()
	list = append(list, Message{Sender: "me", Text: "Bye", Mine: true, At: now.Add(21 * time.Minute)})
	tt.Frame()
	settles(t, tt, "list")
	cues(t, quiet)
	panics(t, "no label", func() { render(func(c *ui.Context) { MessageList(c, "", list, MessageListOptions{}) }) })
	panics(t, "bad message", func() {
		render(func(c *ui.Context) { MessageList(c, "x", []Message{{Sender: "a"}}, MessageListOptions{}) })
	})
	panics(t, "negative height", func() {
		render(func(c *ui.Context) { MessageList(c, "x", list, MessageListOptions{Height: -1}) })
	})
}

func TestConversationList(t *testing.T) {
	items := []Conversation{
		{Name: "Ana Reyes", Last: "See you", At: now, Unread: 2},
		{Name: "Kai Mori", Last: "ok", At: now.Add(-time.Hour)},
		{Name: "Team", Last: "Standup moved", At: now.Add(-2 * time.Hour), Avatar: red},
	}
	selected, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ConversationList(c, "Chats", items, &selected).Changed() {
			changed++
		}
	}), 500, 400)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("rows do not arrive")
	}
	settles(t, tt, "rows")
	if !tt.HasText("2") {
		t.Errorf("no unread count in %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Kai Mori")
	if selected != 1 || changed != 1 {
		t.Fatalf("selected %d changed %d", selected, changed)
	}
	click(t, tt, "Kai Mori")
	if changed != 1 {
		t.Error("reselecting reported a change")
	}
	cues(t, before, sound.Select)
	panics(t, "no items", func() { render(func(c *ui.Context) { ConversationList(c, "x", nil, &selected) }) })
	panics(t, "no time", func() {
		render(func(c *ui.Context) { ConversationList(c, "x", []Conversation{{Name: "a"}}, &selected) })
	})
}

func TestChatSearch(t *testing.T) {
	list := sample()
	query, jump, changed := "", -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, _ := ChatSearch(c, "Search", list, &query, &jump, ChatSearchOptions{})
		if root.Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "search")
	tt.Key(0, ui.KeyTab)
	tt.Type("lunch")
	tt.Frame()
	if query != "lunch" {
		t.Fatalf("query %q", query)
	}
	settles(t, tt, "hits")
	before := baseline()
	click(t, tt, "Lunch at noon?")
	if jump != 1 || changed != 1 {
		t.Errorf("jump %d changed %d", jump, changed)
	}
	cues(t, before, sound.Select)
	query = "zzz"
	tt.Frame()
	settles(t, tt, "empty")
	if !tt.HasText("No messages match") || tt.HasText("Lunch at noon?") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "nil jump", func() {
		render(func(c *ui.Context) { ChatSearch(c, "x", list, &query, nil, ChatSearchOptions{}) })
	})
}

func TestPinned(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Pinned(c, "Ana", "Standup at 10", PinnedOptions{}).Clicked() {
			clicked++
		}
	}), 500, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bar does not drop in")
	}
	settles(t, tt, "bar")
	before := baseline()
	click(t, tt, "Pinned message")
	if clicked != 1 {
		t.Errorf("clicked %d", clicked)
	}
	cues(t, before, sound.Tap)
	panics(t, "no text", func() { render(func(c *ui.Context) { Pinned(c, "a", "", PinnedOptions{}) }) })
}

func TestMessageMenu(t *testing.T) {
	chosen, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root := MessageMenu(c, "Message", &chosen, MessageMenuOptions{}, func() {
			shell(c, false).Children(func() { ui.Text(c, "Can you review the deck?") })
		})
		if root.Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "message")
	if _, ok := tt.Find("Reply"); ok {
		t.Error("menu shown before hover")
	}
	r := find(t, tt, "Can you review the deck?")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	settles(t, tt, "menu")
	before := baseline()
	click(t, tt, "Reply")
	if chosen != 0 || changed != 1 {
		t.Fatalf("chosen %d changed %d", chosen, changed)
	}
	cues(t, before, sound.Tap)
	tt.Move(2, 2)
	tt.Frame()
	settles(t, tt, "menu")
	if _, ok := tt.Find("Reply"); ok {
		t.Error("menu stays after the pointer leaves")
	}
	panics(t, "nil content", func() {
		render(func(c *ui.Context) { MessageMenu(c, "x", &chosen, MessageMenuOptions{}, nil) })
	})
	chosen = -1
	short := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(150).Children(func() {
			MessageMenu(c, "Short", &chosen, MessageMenuOptions{}, func() { ui.Text(c, "Hi") })
		})
	}, 500, 400)
	short.Advance(2 * time.Second)
	r = find(t, short, "Hi")
	short.Move(r.X+r.W/2, r.Y+r.H/2)
	short.Frame()
	short.Advance(2 * time.Second)
	click(t, short, "Reply")
	if chosen != 0 {
		t.Errorf("short message: chosen %d, want 0", chosen)
	}
}
