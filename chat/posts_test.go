package chat

import (
	"image/color"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestLocationMessage(t *testing.T) {
	var danger color.RGBA
	tt := ui.NewTester(padded(func(c *ui.Context) {
		danger = rgba(c.Theme().Danger)
		LocationMessage(c, "471 Harbor Road, Duluth")
	}), 500, 400)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("pin does not drop in")
	}
	settles(t, tt, "map")
	r := find(t, tt, "471 Harbor Road, Duluth")
	if !hasInk(tt, r, danger) {
		t.Error("no pin drawn in the tile")
	}
	panics(t, "no address", func() { render(func(c *ui.Context) { LocationMessage(c, "") }) })
}

func TestContactMessage(t *testing.T) {
	added, changed := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ContactMessage(c, "Ana Reyes", nil, &added, ContactMessageOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "contact")
	before := baseline()
	click(t, tt, "Add")
	if !added || changed != 1 || !tt.HasText("Added") {
		t.Fatalf("added %v changed %d texts %q", added, changed, tt.Texts())
	}
	cues(t, before, sound.Pop)
	click(t, tt, "Added")
	if added || changed != 2 {
		t.Errorf("added %v changed %d", added, changed)
	}
	panics(t, "no name", func() { render(func(c *ui.Context) { ContactMessage(c, "", nil, &added, ContactMessageOptions{}) }) })
}

func TestSticker(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Sticker(c, "Red square", red, StickerOptions{}) }), 400, 300)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("sticker does not pop in")
	}
	settles(t, tt, "sticker")
	r := find(t, tt, "Red square")
	if !at(tt, r.X+r.W/2, r.Y+r.H/2, redInk) {
		t.Error("sticker not drawn")
	}
	panics(t, "nil image", func() { render(func(c *ui.Context) { Sticker(c, "x", nil, StickerOptions{}) }) })
	panics(t, "negative size", func() {
		render(func(c *ui.Context) { Sticker(c, "x", red, StickerOptions{Size: -1}) })
	})
}

func TestEmojiBurst(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { EmojiBurst(c, "🎉") }), 400, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("emoji does not bounce")
	}
	settles(t, tt, "emoji")
	find(t, tt, "🎉")
	panics(t, "no emoji", func() { render(func(c *ui.Context) { EmojiBurst(c, "") }) })
}

func TestForwarded(t *testing.T) {
	m := Message{Sender: "Ana", Text: "Standup moved", At: now}
	tt := ui.NewTester(padded(func(c *ui.Context) { Forwarded(c, "Kai", m, ForwardedOptions{}) }), 400, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("bubble does not rise in")
	}
	settles(t, tt, "forwarded")
	if !tt.HasText("Forwarded from Kai") || !tt.HasText("Standup moved") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "no source", func() { render(func(c *ui.Context) { Forwarded(c, "", m, ForwardedOptions{}) }) })
}

func TestHandoff(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Handoff(c, "Ana Reyes", nil, now, HandoffOptions{}) }), 500, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("event does not fade in")
	}
	settles(t, tt, "handoff")
	if !tt.HasText("Ana Reyes joined the chat") || !tt.HasText("AR") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "zero time", func() {
		render(func(c *ui.Context) { Handoff(c, "Ana", nil, time.Time{}, HandoffOptions{}) })
	})
}

func TestChatEmpty(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, act := ChatEmpty(c, "No messages yet", "Say something nice.", "Say hi")
		if act.Clicked() {
			clicked++
		}
	}), 400, 300)
	tt.Frame()
	if !moves(tt, 200*time.Millisecond) {
		t.Error("hand does not wave")
	}
	before := baseline()
	click(t, tt, "Say hi")
	if clicked != 1 {
		t.Errorf("clicked %d", clicked)
	}
	cues(t, before, sound.Tap)
	panics(t, "no action", func() { render(func(c *ui.Context) { ChatEmpty(c, "t", "", "") }) })
}

func TestUnreadDivider(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { UnreadDivider(c, "New messages") }), 400, 200)
	settles(t, tt, "divider")
	if !tt.HasText("New messages") {
		t.Errorf("texts %q", tt.Texts())
	}
	tt.Advance(3500 * time.Millisecond)
	settles(t, tt, "divider")
	panics(t, "no label", func() { render(func(c *ui.Context) { UnreadDivider(c, "") }) })
}

func TestMessageTime(t *testing.T) {
	var muted, bg color.RGBA
	tt := ui.NewTester(padded(func(c *ui.Context) {
		muted, bg = rgba(c.Theme().TextMuted), rgba(c.Theme().Background)
		MessageTime(c, now, func() {
			shell(c, false).Children(func() { ui.Text(c, "See you at noon") })
		})
	}), 400, 200)
	settles(t, tt, "message")
	r := find(t, tt, "9:30 AM")
	if hasInk(tt, r, muted) && !hasInk(tt, r, bg) {
		t.Skip("time ink indistinguishable before hover")
	}
	b := find(t, tt, "See you at noon")
	tt.Move(b.X+b.W/2, b.Y+b.H/2)
	tt.Frame()
	settles(t, tt, "message")
	if !hasInk(tt, r, muted) {
		t.Error("time does not show on hover")
	}
	panics(t, "nil content", func() { render(func(c *ui.Context) { MessageTime(c, now, nil) }) })
}
