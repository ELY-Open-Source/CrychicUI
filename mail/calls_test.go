package mail

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestContacts(t *testing.T) {
	selected := ""
	people := []Contact{{ID: "ana", Name: "Ana", Email: "ana@example.com"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Contacts(c, people, &selected) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Ana")
	if selected != "ana" || !tt.HasText("A") {
		t.Errorf("selected %q texts %q", selected, tt.Texts())
	}
	cues(t, before, sound.Select)
}

func TestContactCard(t *testing.T) {
	action := ""
	person := Contact{ID: "ana", Name: "Ana", Email: "ana@example.com", Phone: "123"}
	tt := ui.NewTester(padded(func(c *ui.Context) { ContactCard(c, person, &action, ContactCardOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Call")
	if action != "call" {
		t.Errorf("action %q", action)
	}
	cues(t, before, sound.Tap)
}

func TestCallScreen(t *testing.T) {
	ended := false
	started := time.Now().Add(-65 * time.Second)
	tt := ui.NewTester(padded(func(c *ui.Context) { CallScreen(c, "Ana", started, &ended, CallScreenOptions{}) }), 500, 250)
	tt.Advance(100 * time.Millisecond)
	if !tt.HasText("01:05") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "End")
	if !ended {
		t.Error("call not ended")
	}
	cues(t, before, sound.Close)
}

func TestIncomingCall(t *testing.T) {
	accepted, declined := false, false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		IncomingCall(c, "Ana", &accepted, &declined, IncomingCallOptions{})
	}), 500, 300)
	tt.Advance(100 * time.Millisecond)
	before := baseline()
	click(t, tt, "Accept")
	if !accepted || declined {
		t.Errorf("accepted %v declined %v", accepted, declined)
	}
	cues(t, before, sound.Notify)
}

func TestCallControls(t *testing.T) {
	state := CallState{Video: true}
	tt := ui.NewTester(padded(func(c *ui.Context) { CallControls(c, &state, CallControlsOptions{}) }), 600, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Mute")
	click(t, tt, "End")
	if !state.Muted || !state.Video || !state.Ended {
		t.Errorf("state %+v", state)
	}
	cues(t, before, sound.Toggle, sound.Toggle)
}

func TestVideoGrid(t *testing.T) {
	people := []Participant{{ID: "ana", Name: "Ana", Speaking: true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { VideoGrid(c, people) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("Ana") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestParticipants(t *testing.T) {
	people := []Participant{{ID: "ana", Name: "Ana", Muted: true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Participants(c, people) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("Ana") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestScreenShare(t *testing.T) {
	sharing := true
	tt := ui.NewTester(padded(func(c *ui.Context) { ScreenShare(c, &sharing, ScreenShareOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Stop")
	if sharing {
		t.Error("sharing still active")
	}
	cues(t, before, sound.Close)
}

func TestVoicemail(t *testing.T) {
	playing := -1
	messages := []VoiceMessage{{ID: "v1", Sender: "Ana", Time: "1m", Transcript: "Call me"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Voicemail(c, messages, &playing, VoicemailOptions{}) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Play")
	if playing != 0 || !tt.HasText("Call me") {
		t.Errorf("playing %d texts %q", playing, tt.Texts())
	}
	cues(t, before, sound.Toggle)
}

func TestMessageTemplates(t *testing.T) {
	selected := -1
	replies := []SavedReply{{ID: "r1", Title: "Thanks", Body: "Thank you"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { MessageTemplates(c, replies, &selected) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Thanks")
	if selected != 0 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}

func TestUnsubscribe(t *testing.T) {
	reason := ""
	confirmed, reported := false, false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, confirm := Unsubscribe(c, []string{"Too many"}, &reason, &confirmed, UnsubscribeOptions{})
		reported = reported || confirm.Changed()
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Too many")
	click(t, tt, "Unsubscribe")
	if !confirmed || !reported || reason != "Too many" {
		t.Errorf("reason %q confirmed %v reported %v", reason, confirmed, reported)
	}
	cues(t, before, sound.Select, sound.Success)
}

func TestSpam(t *testing.T) {
	spam := true
	tt := ui.NewTester(padded(func(c *ui.Context) { Spam(c, &spam, SpamOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Not spam")
	if spam {
		t.Error("still spam")
	}
	cues(t, before, sound.Tap)
}

func TestPriorityInbox(t *testing.T) {
	other := testMail()
	other.ID, other.Subject = "m2", "Later"
	tt := ui.NewTester(padded(func(c *ui.Context) {
		PriorityInbox(c, []Mail{testMail()}, []Mail{other}, PriorityInboxOptions{})
	}), 500, 350)
	settles(t, tt)
	if !tt.HasText("Important") || !tt.HasText("Other") || !tt.HasText("Later") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestMailSearch(t *testing.T) {
	query := MailQuery{}
	tt := ui.NewTester(padded(func(c *ui.Context) { MailSearch(c, &query, MailSearchOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Has attachment")
	if !query.Attachments || query.Links {
		t.Errorf("query %+v", query)
	}
	cues(t, before, sound.Toggle)
}

func TestUndoSend(t *testing.T) {
	visible := true
	tt := ui.NewTester(padded(func(c *ui.Context) { UndoSend(c, &visible, UndoSendOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Undo")
	if visible {
		t.Error("toast still visible")
	}
	cues(t, before, sound.Close)
}

func TestReadingPane(t *testing.T) {
	selected := "m1"
	second := testMail()
	second.ID, second.Subject, second.Body = "m2", "Second", "Second body"
	mails := []Mail{testMail(), second}
	tt := ui.NewTester(padded(func(c *ui.Context) { ReadingPane(c, mails, &selected) }), 700, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Second")
	if selected != "m2" || !tt.HasText("Second body") {
		t.Errorf("selected %q texts %q", selected, tt.Texts())
	}
	cues(t, before, sound.Select)
}
