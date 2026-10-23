package chat

import (
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestCallBanner(t *testing.T) {
	joined, changed := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CallBanner(c, "Call in progress", &joined, CallBannerOptions{}).Changed() {
			changed++
		}
	}), 500, 200)
	tt.Frame()
	if !moves(tt, 200*time.Millisecond) {
		t.Error("beacon does not pulse")
	}
	before := baseline()
	click(t, tt, "Join")
	if !joined || changed != 1 || !tt.HasText("Leave") {
		t.Fatalf("joined %v changed %d texts %q", joined, changed, tt.Texts())
	}
	cues(t, before, sound.Tap)
	click(t, tt, "Leave")
	if joined || changed != 2 {
		t.Errorf("joined %v changed %d", joined, changed)
	}
	panics(t, "no title", func() { render(func(c *ui.Context) { CallBanner(c, "", &joined, CallBannerOptions{}) }) })
}

func TestJumpToLatest(t *testing.T) {
	count, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if JumpToLatest(c, count, JumpToLatestOptions{}).Changed() {
			changed++
		}
	}), 400, 200)
	tt.Frame()
	if tt.HasText("new message") {
		t.Errorf("pill shown at zero: %q", tt.Texts())
	}
	count = 3
	tt.Frame()
	if !tt.FramePending() {
		t.Error("pill does not pop in")
	}
	settles(t, tt, "pill")
	before := baseline()
	click(t, tt, "3 new messages")
	if changed != 1 {
		t.Errorf("changed %d", changed)
	}
	cues(t, before, sound.Tap)
	count = 1
	tt.Frame()
	if !tt.HasText("1 new message") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "negative count", func() { render(func(c *ui.Context) { JumpToLatest(c, -1, JumpToLatestOptions{}) }) })
}

func TestMessageStatus(t *testing.T) {
	state := Sending
	retried, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if MessageStatus(c, state, func() { retried++ }, MessageStatusOptions{}).Changed() {
			changed++
		}
	}), 400, 200)
	tt.Frame()
	if !moves(tt, 200*time.Millisecond) {
		t.Error("clock does not spin")
	}
	quiet := baseline()
	state = Failed
	tt.Frame()
	settles(t, tt, "status")
	cues(t, quiet)
	before := baseline()
	click(t, tt, "Retry")
	if retried != 1 || changed != 1 {
		t.Errorf("retried %d changed %d", retried, changed)
	}
	cues(t, before, sound.Send)
	panics(t, "nil retry", func() { render(func(c *ui.Context) { MessageStatus(c, Sending, nil, MessageStatusOptions{}) }) })
	panics(t, "unknown state", func() {
		render(func(c *ui.Context) { MessageStatus(c, 9, func() {}, MessageStatusOptions{}) })
	})
}

func TestTranslate(t *testing.T) {
	shown, changed := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Translate(c, "See you at noon", "On se voit à midi", &shown, TranslateOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt, "translate")
	before := baseline()
	click(t, tt, "Show translation")
	if !shown || changed != 1 {
		t.Fatalf("shown %v changed %d", shown, changed)
	}
	settles(t, tt, "translate")
	if !tt.HasText("On se voit à midi") || tt.HasText("See you at noon") {
		t.Errorf("texts %q", tt.Texts())
	}
	cues(t, before, sound.Toggle)
	click(t, tt, "Show original")
	settles(t, tt, "translate")
	if !tt.HasText("See you at noon") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "empty original", func() {
		render(func(c *ui.Context) { Translate(c, "", "x", &shown, TranslateOptions{}) })
	})
}

func TestCodeMessage(t *testing.T) {
	code := "fmt.Println(\"hi\")"
	copied := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, cp := CodeMessage(c, "go", code, CodeMessageOptions{})
		if cp.Clicked() {
			copied++
		}
	}), 500, 300)
	settles(t, tt, "code")
	before := baseline()
	click(t, tt, "Copy")
	if copied != 1 || tt.Clipboard() != code {
		t.Errorf("copied %d clipboard %q", copied, tt.Clipboard())
	}
	cues(t, before, sound.Pop)
	panics(t, "no code", func() { render(func(c *ui.Context) { CodeMessage(c, "go", "", CodeMessageOptions{}) }) })
}

func TestMarkdownMessage(t *testing.T) {
	m := Message{Sender: "Ana", Text: "**Bold** and `code`", At: now}
	tt := ui.NewTester(padded(func(c *ui.Context) { MarkdownMessage(c, m) }), 500, 300)
	settles(t, tt, "markdown")
	joined := strings.Join(tt.Texts(), " ")
	if !strings.Contains(joined, "Bold") || !strings.Contains(joined, "code") {
		t.Errorf("texts %q", joined)
	}
	fenced := ui.NewTester(padded(func(c *ui.Context) {
		MarkdownMessage(c, Message{Sender: "Ana", Text: "```go\nfmt.Println(1)\n```", At: now})
	}), 500, 300)
	settles(t, fenced, "fenced")
	recent := baseline()
	click(t, fenced, "Copy")
	if fenced.Clipboard() != "fmt.Println(1)\n" {
		t.Errorf("clipboard %q", fenced.Clipboard())
	}
	cues(t, recent, sound.Pop)
	linked := ui.NewTester(padded(func(c *ui.Context) {
		MarkdownMessage(c, Message{Sender: "Ana", Text: "[Read docs](https://example.com/docs)", At: now})
	}), 500, 300)
	settles(t, linked, "linked")
	again := baseline()
	click(t, linked, "Read docs")
	if got := linked.OpenedURLs(); len(got) != 1 || got[0] != "https://example.com/docs" {
		t.Errorf("opened %v", got)
	}
	cues(t, again, sound.Open)
	src := "[Jump to target](#target)\n\n" + strings.Repeat("A filler paragraph.\n\n", 20) + "# Target\n"
	st := ui.ScrollState{}
	scrolled := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() {
			ui.Scroll(c).Height(150).TrackScroll(&st).Children(func() {
				MarkdownMessage(c, Message{Sender: "Ana", Text: src, At: now})
			})
		})
	}, 400, 300)
	scrolled.Advance(2 * time.Second)
	third := baseline()
	click(t, scrolled, "Jump to target")
	scrolled.Advance(time.Second)
	if st.Y <= 0 {
		t.Errorf("anchor did not scroll: Y %v of %v", st.Y, st.MaxY)
	}
	cues(t, third, sound.Open)
	flat := "[Jump to target](#target)\n\nBeginning.\n\n" + strings.Repeat("A filler paragraph.\n\n", 20) + "# !!!\n\nDestination.\n"
	st2 := ui.ScrollState{}
	quiet := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() {
			ui.Scroll(c).Height(150).TrackScroll(&st2).Children(func() {
				MarkdownMessage(c, Message{Sender: "Ana", Text: flat, At: now})
			})
		})
	}, 400, 300)
	quiet.Advance(2 * time.Second)
	if st2.Y != 0 {
		t.Errorf("empty slug scrolled without a click: Y %v", st2.Y)
	}
	hd := ui.NewTester(padded(func(c *ui.Context) {
		MarkdownMessage(c, Message{Sender: "Ana", Text: "# See [Read docs](https://example.com/docs)", At: now})
	}), 500, 300)
	settles(t, hd, "heading")
	fourth := baseline()
	click(t, hd, "Read docs")
	if got := hd.OpenedURLs(); len(got) != 1 {
		t.Errorf("heading link opened %v", got)
	}
	cues(t, fourth, sound.Open)
	dup := "Go [Jump](#target) now.\n\n# Target\n\n" + strings.Repeat("A filler paragraph.\n\n", 30) + "# Target\n"
	st3 := ui.ScrollState{}
	dups := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() {
			ui.Scroll(c).Height(150).TrackScroll(&st3).Children(func() {
				MarkdownMessage(c, Message{Sender: "Ana", Text: dup, At: now})
			})
		})
	}, 400, 300)
	dups.Advance(2 * time.Second)
	click(t, dups, "Jump")
	dups.Advance(time.Second)
	if st3.Y != 0 {
		t.Errorf("duplicate headings jumped past the first: Y %v of %v", st3.Y, st3.MaxY)
	}
	panics(t, "bad message", func() { render(func(c *ui.Context) { MarkdownMessage(c, Message{Sender: "a"}) }) })
}

func TestHoldToRecord(t *testing.T) {
	recording, changed, cancelled := false, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root := HoldToRecord(c, "Voice note", &recording, func() { cancelled++ }, HoldToRecordOptions{})
		if root.Changed() {
			changed++
		}
	}), 400, 200)
	settles(t, tt, "recorder")
	before := baseline()
	r := find(t, tt, "Voice note")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	if !recording || changed != 1 {
		t.Fatalf("recording %v changed %d", recording, changed)
	}
	tt.Release(x, y)
	tt.Frame()
	if recording || changed != 2 {
		t.Fatalf("recording %v changed %d", recording, changed)
	}
	tt.Press(x, y)
	tt.Frame()
	tt.Move(x, y-100)
	tt.Frame()
	if recording || cancelled != 1 {
		t.Fatalf("teleport slide: recording %v cancelled %d", recording, cancelled)
	}
	tt.Release(x, y-100)
	tt.Frame()
	tt.Press(x, y)
	tt.Frame()
	for range 10 {
		tt.Move(x, y-10)
		tt.Frame()
		tt.Advance(30 * time.Millisecond)
		y -= 10
	}
	if recording || cancelled != 2 {
		t.Fatalf("gradual slide: recording %v cancelled %d", recording, cancelled)
	}
	tt.Release(x, y)
	tt.Frame()
	if changed != 6 {
		t.Errorf("changed %d, want 6", changed)
	}
	cues(t, before, sound.Toggle, sound.Toggle, sound.Toggle, sound.Close, sound.Toggle, sound.Close)
	panics(t, "nil cancel", func() {
		render(func(c *ui.Context) { HoldToRecord(c, "x", &recording, nil, HoldToRecordOptions{}) })
	})
}

func TestChatHeader(t *testing.T) {
	voice, video, changed := 0, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, v, d := ChatHeader(c, "Ana Reyes", "Online", nil, ChatHeaderOptions{})
		if root.Changed() {
			changed++
		}
		if v.Clicked() {
			voice++
		}
		if d.Clicked() {
			video++
		}
	}), 500, 200)
	settles(t, tt, "header")
	before := baseline()
	click(t, tt, "Voice call")
	if voice != 1 || changed != 1 {
		t.Errorf("voice %d changed %d", voice, changed)
	}
	click(t, tt, "Video call")
	if video != 1 || changed != 2 {
		t.Errorf("video %d changed %d", video, changed)
	}
	cues(t, before, sound.Tap, sound.Tap)
	panics(t, "no status", func() {
		render(func(c *ui.Context) { ChatHeader(c, "Ana", "", nil, ChatHeaderOptions{}) })
	})
}

func TestChatWindow(t *testing.T) {
	list := sample()
	var sent []string
	changed := 0
	voice, attach := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		root, _, v, _, at := ChatWindow(c, "Ana Reyes", "Online", list, func(s string) { sent = append(sent, s) }, ChatWindowOptions{Height: 300})
		if root.Changed() {
			changed++
		}
		if v.Clicked() {
			voice++
		}
		if at.Clicked() {
			attach++
		}
	}, 500, 600)
	settles(t, tt, "window")
	if !tt.HasText("Online") || !tt.HasText("Morning!") {
		t.Fatalf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Message")
	tt.Type("on my way")
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if len(sent) != 1 || sent[0] != "on my way" || changed == 0 {
		t.Fatalf("sent %q changed %d", sent, changed)
	}
	cues(t, before, sound.Send)
	again := baseline()
	click(t, tt, "Voice call")
	if voice != 1 {
		t.Errorf("voice %d", voice)
	}
	click(t, tt, "Attach")
	if attach != 1 {
		t.Errorf("attach %d", attach)
	}
	cues(t, again, sound.Tap, sound.Tap)
	panics(t, "nil send", func() {
		render(func(c *ui.Context) { ChatWindow(c, "Ana", "Online", list, nil, ChatWindowOptions{}) })
	})
}
