package chat

import (
	"image/color"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestBubble(t *testing.T) {
	var acc color.RGBA
	m := Message{Sender: "me", Text: "Hi there", Mine: true, At: now}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		acc = rgba(c.Theme().Accent)
		Bubble(c, m, BubbleOptions{ShowTime: true})
	}), 400, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bubble does not rise in")
	}
	settles(t, tt, "bubble")
	r := find(t, tt, "Hi there")
	if !at(tt, r.X-4, r.Y+r.H/2, acc) {
		t.Error("my bubble is not accent")
	}
	if r.X < 200 {
		t.Errorf("my bubble at x %v, want the right side", r.X)
	}
	if !tt.HasText("9:30") {
		t.Errorf("no time in %q", tt.Texts())
	}
	panics(t, "no text", func() { render(func(c *ui.Context) { Bubble(c, Message{Sender: "a", At: now}, BubbleOptions{}) }) })
	panics(t, "no time", func() { render(func(c *ui.Context) { Bubble(c, Message{Sender: "a", Text: "x"}, BubbleOptions{}) }) })
}

func TestTyping(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Typing(c, "Ana is typing") }), 300, 200)
	tt.Frame()
	if !moves(tt, 150*time.Millisecond) {
		t.Error("dots do not bounce")
	}
	tt.Advance(2 * time.Second)
	if !tt.FramePending() {
		t.Error("dots stopped")
	}
	panics(t, "no label", func() { render(func(c *ui.Context) { Typing(c, "") }) })
}

func TestReceipt(t *testing.T) {
	state := Sent
	tt := ui.NewTester(padded(func(c *ui.Context) { Receipt(c, state, ReceiptOptions{}) }), 300, 200)
	settles(t, tt, "receipt")
	if !tt.HasText("Sent") {
		t.Errorf("texts %q", tt.Texts())
	}
	state = Read
	tt.Frame()
	if !tt.FramePending() {
		t.Error("check does not cross-fade")
	}
	settles(t, tt, "receipt")
	if !tt.HasText("Read") || tt.HasText("Sent") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "unknown state", func() { render(func(c *ui.Context) { Receipt(c, 9, ReceiptOptions{}) }) })
}

func TestReactions(t *testing.T) {
	items := []Reaction{{Emoji: "👍", Count: 2}, {Emoji: "🎉", Count: 1, Mine: true}}
	changed := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Reactions(c, "Reactions", &items).Changed() {
			changed++
		}
	}), 400, 200)
	settles(t, tt, "reactions")
	before := baseline()
	click(t, tt, "👍")
	if items[0].Count != 3 || !items[0].Mine || changed != 1 {
		t.Fatalf("after a tap: %+v changed %d", items[0], changed)
	}
	if !tt.FramePending() {
		t.Error("emoji does not pop")
	}
	settles(t, tt, "reactions")
	click(t, tt, "👍")
	if items[0].Count != 2 || items[0].Mine {
		t.Errorf("after a second tap: %+v", items[0])
	}
	cues(t, before, sound.Pop, sound.Pop)
	panics(t, "no items", func() { render(func(c *ui.Context) { Reactions(c, "x", nil) }) })
	panics(t, "negative count", func() {
		render(func(c *ui.Context) { Reactions(c, "x", &[]Reaction{{Emoji: "a", Count: -1}}) })
	})
	panics(t, "mine without a count", func() {
		render(func(c *ui.Context) { Reactions(c, "x", &[]Reaction{{Emoji: "a", Mine: true}}) })
	})
}

func TestReactionPicker(t *testing.T) {
	chosen := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { ReactionPicker(c, "React", []string{"👍", "❤️", "😂"}, &chosen) }), 400, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("picker does not scale in")
	}
	settles(t, tt, "picker")
	before := baseline()
	click(t, tt, "❤️")
	if chosen != "❤️" {
		t.Errorf("chosen %q", chosen)
	}
	cues(t, before, sound.Select)
	panics(t, "no emoji", func() { render(func(c *ui.Context) { ReactionPicker(c, "x", nil, &chosen) }) })
	panics(t, "empty emoji", func() { render(func(c *ui.Context) { ReactionPicker(c, "x", []string{""}, &chosen) }) })
}

func TestReplyQuote(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ReplyQuote(c, "Ana", "See you at noon", ReplyQuoteOptions{}).Clicked() {
			clicked++
		}
	}), 400, 200)
	settles(t, tt, "quote")
	before := baseline()
	click(t, tt, "Quoted message")
	if clicked != 1 {
		t.Errorf("clicked %d", clicked)
	}
	cues(t, before, sound.Tap)
	panics(t, "no sender", func() { render(func(c *ui.Context) { ReplyQuote(c, "", "x", ReplyQuoteOptions{}) }) })
}

func TestThread(t *testing.T) {
	open := false
	tt := ui.NewTester(padded(func(c *ui.Context) { Thread(c, 3, &open, ThreadOptions{}) }), 400, 200)
	settles(t, tt, "thread")
	before := baseline()
	click(t, tt, "3 replies")
	if !open {
		t.Fatal("thread not open")
	}
	if !tt.FramePending() {
		t.Error("chevron does not turn")
	}
	settles(t, tt, "thread")
	click(t, tt, "3 replies")
	if open {
		t.Error("thread not closed")
	}
	cues(t, before, sound.Open, sound.Close)
	one := ui.NewTester(padded(func(c *ui.Context) { Thread(c, 1, &open, ThreadOptions{}) }), 400, 200)
	one.Frame()
	if !one.HasText("1 reply") {
		t.Errorf("texts %q", one.Texts())
	}
	own := ui.NewTester(padded(func(c *ui.Context) { Thread(c, 1, &open, ThreadOptions{One: "%d answer", Many: "%d answers"}) }), 400, 200)
	own.Frame()
	if !own.HasText("1 answer") {
		t.Errorf("texts %q", own.Texts())
	}
	panics(t, "negative count", func() { render(func(c *ui.Context) { Thread(c, -1, &open, ThreadOptions{}) }) })
}

func TestEditedMessage(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { EditedMessage(c, now, EditedMessageOptions{}) }), 400, 200)
	settles(t, tt, "mark")
	r := find(t, tt, "(edited)")
	shown := func() bool {
		for _, s := range tt.Texts() {
			if strings.Contains(s, "2026") {
				return true
			}
		}
		return false
	}
	if shown() {
		t.Error("time shown before hover")
	}
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(800 * time.Millisecond)
	if !shown() {
		t.Errorf("no time in tip: %q", tt.Texts())
	}
	panics(t, "zero time", func() { render(func(c *ui.Context) { EditedMessage(c, time.Time{}, EditedMessageOptions{}) }) })
}

func TestDeletedMessage(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { DeletedMessage(c, DeletedMessageOptions{}) }), 400, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("placeholder does not fade in")
	}
	settles(t, tt, "placeholder")
	find(t, tt, "Message deleted")
	own := ui.NewTester(padded(func(c *ui.Context) { DeletedMessage(c, DeletedMessageOptions{Text: "Gone"}) }), 400, 200)
	own.Frame()
	if !own.HasText("Gone") {
		t.Errorf("texts %q", own.Texts())
	}
}

func TestMention(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Mention(c, "ana").Clicked() {
			clicked++
		}
	}), 400, 200)
	settles(t, tt, "mention")
	before := baseline()
	click(t, tt, "@ana")
	if clicked != 1 {
		t.Errorf("clicked %d", clicked)
	}
	cues(t, before, sound.Tap)
	panics(t, "no name", func() { render(func(c *ui.Context) { Mention(c, "") }) })
}

func TestSystemMessage(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { SystemMessage(c, "Ana joined") }), 400, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("line does not fade in")
	}
	settles(t, tt, "system line")
	r := find(t, tt, "Ana joined")
	if mid := r.X + r.W/2; mid < 150 || mid > 250 {
		t.Errorf("line centered at %v", mid)
	}
	panics(t, "no text", func() { render(func(c *ui.Context) { SystemMessage(c, "") }) })
}

func TestDayDivider(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { DayDivider(c, "Today") }), 400, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("divider does not fade in")
	}
	settles(t, tt, "divider")
	r := find(t, tt, "Today")
	if mid := r.X + r.W/2; mid < 150 || mid > 250 {
		t.Errorf("label centered at %v", mid)
	}
	panics(t, "no label", func() { render(func(c *ui.Context) { DayDivider(c, "") }) })
}

func TestLinkPreview(t *testing.T) {
	var img *ui.Bitmap
	tt := ui.NewTester(padded(func(c *ui.Context) { LinkPreview(c, "MyGo docs", "mygo.egoist.dev", img) }), 400, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("card does not fade in")
	}
	settles(t, tt, "preview")
	find(t, tt, "MyGo docs")
	host := find(t, tt, "mygo.egoist.dev")
	img = red
	tt.Frame()
	settles(t, tt, "preview")
	if r := find(t, tt, "mygo.egoist.dev"); r.Y <= host.Y {
		t.Error("image did not push the text down")
	}
	panics(t, "no host", func() { render(func(c *ui.Context) { LinkPreview(c, "x", "", nil) }) })
}
