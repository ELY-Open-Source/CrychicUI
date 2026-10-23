package gallery

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

var errNoText = errors.New("text not shown")

func TestChatDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Bubble":      {press("Flip side"), "mine true"},
		"MessageList": {press("Receive"), "5 messages"},
		"Composer": {func(tt *ui.Tester) error {
			if err := tt.Click("Message"); err != nil {
				return err
			}
			tt.Type("hi")
			tt.Key(0, ui.KeyEnter)
			return nil
		}, "sent 2 runes"},
		"Receipt":          {press("Advance"), "state 1"},
		"Reactions":        {press("👍"), "counts 4 1 0"},
		"ReactionPicker":   {press("❤️"), "picked ❤️"},
		"ReplyQuote":       {press("Quoted message"), "jumped"},
		"Thread":           {press("3 replies"), "open true"},
		"Attachment":       {press("report.pdf"), "opened report.pdf"},
		"ImageMessage":     {press("Harbor photo"), "open true"},
		"VoiceMessage":     {press("Play"), "playing true"},
		"LinkPreview":      {press("Toggle image"), "image false"},
		"ConversationList": {press("Kai Mori"), "selected Kai Mori"},
		"ChatSearch": {func(tt *ui.Tester) error {
			tt.Key(0, ui.KeyTab)
			tt.Type("lunch")
			tt.Frame()
			return tt.Click("Lunch at noon?")
		}, "jump 1"},
		"Pinned": {press("Pinned message"), "jumped"},
		"MessageMenu": {func(tt *ui.Tester) error {
			r, ok := tt.Find("Can you review the deck before then?")
			if !ok {
				return errNoText
			}
			tt.Move(r.X+r.W/2, r.Y+r.H/2)
			tt.Frame()
			return tt.Click("Copy")
		}, "action copy"},
		"Mention":        {press("@ana"), "opened @ana"},
		"QuickReplies":   {press("Call me"), "reply: Call me"},
		"PollMessage":    {press("Salad"), "votes [3 2 2]"},
		"ContactMessage": {press("Add"), "added true"},
		"ChannelList":    {press("random"), "selected random"},
		"Members":        {press("Bo Tan"), "picked Bo Tan"},
		"CallBanner":     {press("Join"), "joined true"},
		"JumpToLatest":   {press("5 new messages"), "jumped"},
		"MessageStatus":  {press("Retry"), "retried"},
		"Translate":      {press("Show translation"), "translated true"},
		"CodeMessage":    {press("Copy"), "copied 21 runes"},
		"ChatEmpty":      {press("Say hi"), "waved hello"},
		"OnlineBar":      {press("Kai joins"), "2 online"},
		"HoldToRecord": {func(tt *ui.Tester) error {
			r, ok := tt.Find("Voice note")
			if !ok {
				return errNoText
			}
			tt.Press(r.X+r.W/2, r.Y+r.H/2)
			tt.Frame()
			tt.Release(r.X+r.W/2, r.Y+r.H/2)
			return nil
		}, "recording true"},
		"ChatHeader": {press("Voice call"), "voice call"},
		"ChatWindow": {func(tt *ui.Tester) error {
			if err := tt.Click("Message"); err != nil {
				return err
			}
			tt.Type("hi")
			tt.Key(0, ui.KeyEnter)
			return nil
		}, "sent 2 runes"},
	}
	// Quiet demos only show; they have no events.
	quiet := map[string]bool{"Typing": true, "SystemMessage": true, "DayDivider": true, "EditedMessage": true, "DeletedMessage": true,
		"LocationMessage": true, "Sticker": true, "EmojiBurst": true, "MarkdownMessage": true, "Forwarded": true,
		"MessageTime": true, "Handoff": true, "UnreadDivider": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "chat" {
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
		tt.Advance(time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		found := false
		for _, l := range logs {
			found = found || strings.HasPrefix(l, act.want)
		}
		if !found {
			t.Errorf("%s logged %q, want a line starting %q", e.Name, logs, act.want)
		}
	}
	if seen != len(acts)+len(quiet) {
		t.Errorf("checked %d chat demos, want %d", seen, len(acts)+len(quiet))
	}
}
