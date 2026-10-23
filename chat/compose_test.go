package chat

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestComposer(t *testing.T) {
	draft, changed := "", 0
	var sent []string
	var attached int
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, send, attach := Composer(c, &draft, func(s string) { sent = append(sent, s) }, ComposerOptions{})
		if root.Changed() {
			changed++
		}
		if send.Clicked() && draft != "" {
			t.Error("send clicked with a draft left")
		}
		if attach.Clicked() {
			attached++
		}
	}), 500, 300)
	settles(t, tt, "composer")
	before := baseline()
	click(t, tt, "Message")
	tt.Type("hello")
	tt.Frame()
	if draft != "hello" || changed == 0 {
		t.Fatalf("draft %q changed %d after typing", draft, changed)
	}
	edits := changed
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(sent) != 1 || sent[0] != "hello" || draft != "" || changed <= edits {
		t.Fatalf("sent %q draft %q changed %d", sent, draft, changed)
	}
	cues(t, before, sound.Send)
	tt.Type("   ")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(sent) != 1 {
		t.Errorf("blank text sent: %q", sent)
	}
	tt.Command("selectAll")
	tt.Type("yo")
	tt.Frame()
	edits = changed
	again := baseline()
	click(t, tt, "Send")
	if len(sent) != 2 || sent[1] != "yo" || draft != "" || changed <= edits {
		t.Errorf("sent %q draft %q changed %d", sent, draft, changed)
	}
	cues(t, again, sound.Send)
	third := baseline()
	click(t, tt, "Attach")
	if attached != 1 {
		t.Errorf("attached %d", attached)
	}
	cues(t, third, sound.Tap)
	panics(t, "nil send", func() { render(func(c *ui.Context) { Composer(c, &draft, nil, ComposerOptions{}) }) })
}

func TestQuickReplies(t *testing.T) {
	chosen := -1
	tt := ui.NewTester(padded(func(c *ui.Context) { QuickReplies(c, "Quick replies", []string{"Yes", "No", "Later"}, &chosen) }), 500, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("chips do not arrive")
	}
	settles(t, tt, "chips")
	before := baseline()
	click(t, tt, "No")
	if chosen != 1 {
		t.Errorf("chosen %d", chosen)
	}
	cues(t, before, sound.Send)
	panics(t, "no replies", func() { render(func(c *ui.Context) { QuickReplies(c, "x", nil, &chosen) }) })
	panics(t, "choice out of range", func() {
		n := 3
		render(func(c *ui.Context) { QuickReplies(c, "x", []string{"a"}, &n) })
	})
}

func TestAttachment(t *testing.T) {
	clicked := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Attachment(c, "report.pdf", 1234567).Clicked() {
			clicked++
		}
	}), 400, 200)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("chip does not rise in")
	}
	settles(t, tt, "chip")
	if !tt.HasText("1.2 MB") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "report.pdf")
	if clicked != 1 {
		t.Errorf("clicked %d", clicked)
	}
	cues(t, before, sound.Tap)
	panics(t, "negative size", func() { render(func(c *ui.Context) { Attachment(c, "x", -1) }) })
}

func TestImageMessage(t *testing.T) {
	open, changed := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ImageMessage(c, "Photo", red, &open, ImageMessageOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	settles(t, tt, "thumbnail")
	r := find(t, tt, "Photo")
	if !at(tt, r.X+r.W/2, r.Y+r.H/2, redInk) {
		t.Error("thumbnail not drawn")
	}
	before := baseline()
	click(t, tt, "Photo")
	if !open || changed != 1 {
		t.Fatalf("open %v changed %d", open, changed)
	}
	settles(t, tt, "lightbox")
	if !at(tt, 250, 200, redInk) {
		t.Error("lightbox does not show the photo")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || changed != 2 {
		t.Errorf("after Escape: open %v changed %d", open, changed)
	}
	cues(t, before, sound.Open, sound.Close)
	panics(t, "nil image", func() { render(func(c *ui.Context) { ImageMessage(c, "x", nil, &open, ImageMessageOptions{}) }) })
	panics(t, "negative width", func() {
		render(func(c *ui.Context) { ImageMessage(c, "x", red, &open, ImageMessageOptions{Width: -1}) })
	})
}

func TestVoiceMessage(t *testing.T) {
	playing, changed := false, 0
	levels := []float64{0.2, 0.8, 0.5, 1, 0.3}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if VoiceMessage(c, "Voice note", levels, 3*time.Second, &playing, VoiceMessageOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	settles(t, tt, "voice")
	if !tt.HasText("0:03") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Play")
	if !playing || changed != 1 {
		t.Fatalf("playing %v changed %d", playing, changed)
	}
	cues(t, before, sound.Toggle)
	tt.Advance(1100 * time.Millisecond)
	if !tt.HasText("0:02") || !tt.HasText("Pause") {
		t.Errorf("texts %q after a second", tt.Texts())
	}
	tt.Advance(2500 * time.Millisecond)
	if playing || changed != 2 || !tt.HasText("0:00") {
		t.Errorf("at the end: playing %v changed %d texts %q", playing, changed, tt.Texts())
	}
	settles(t, tt, "voice")
	panics(t, "zero length", func() {
		render(func(c *ui.Context) { VoiceMessage(c, "x", levels, 0, &playing, VoiceMessageOptions{}) })
	})
}
