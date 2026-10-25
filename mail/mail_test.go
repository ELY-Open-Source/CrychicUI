package mail

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestInbox(t *testing.T) {
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Inbox(c, []Mail{testMail()}, &selected, InboxOptions{}) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Release")
	if selected != "m1" {
		t.Errorf("selected %q", selected)
	}
	cues(t, before, sound.Select)
}

func TestMailRowSwipe(t *testing.T) {
	archived := false
	tt := ui.NewTester(padded(func(c *ui.Context) { MailRow(c, testMail(), &archived) }), 500, 220)
	settles(t, tt)
	r, _ := tt.Find("Release")
	before := baseline()
	tt.Press(r.X+5, r.Y+5)
	tt.Move(r.X+500, r.Y+5)
	tt.Release(r.X+500, r.Y+5)
	tt.Advance(time.Second)
	if !archived {
		t.Error("message not archived")
	}
	cues(t, before, sound.Swipe)
}

func TestMailRowSwipeWide(t *testing.T) {
	archived := false
	tt := ui.NewTester(padded(func(c *ui.Context) { MailRow(c, testMail(), &archived) }), 900, 700)
	settles(t, tt)
	r, _ := tt.Find("Release")
	tt.Press(r.X+5, r.Y+5)
	tt.Move(r.X+600, r.Y+5)
	tt.Release(r.X+600, r.Y+5)
	tt.Advance(3 * time.Second)
	if !archived {
		t.Errorf("message not archived at rect %v", r)
	}
}

func TestComposeWindow(t *testing.T) {
	value := Compose{To: "ana@example.com", Subject: "Release", Body: "Ready"}
	sent := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, send := ComposeWindow(c, &value, ComposeWindowOptions{})
		sent = sent || send.Changed()
	}), 500, 400)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Send")
	if !sent {
		t.Error("send not reported")
	}
	cues(t, before, sound.Send)
}

func TestComposeWindowValidation(t *testing.T) {
	value := Compose{To: "bad", Subject: "Release", Body: "Ready"}
	sent := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, send := ComposeWindow(c, &value, ComposeWindowOptions{})
		sent = sent || send.Changed()
	}), 500, 400)
	settles(t, tt)
	if err := tt.Click("Send"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if sent {
		t.Error("invalid message sent")
	}
}

func TestRecipients(t *testing.T) {
	to := []string{"ana@example.com", "ben@example.com"}
	cc, bcc := []string{}, []string{}
	tt := ui.NewTester(padded(func(c *ui.Context) { Recipients(c, &to, &cc, &bcc, RecipientsOptions{}) }), 600, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "ana@example.com")
	if len(to) != 1 || to[0] != "ben@example.com" {
		t.Errorf("recipients %v", to)
	}
	cues(t, before, sound.Pop)
}

func TestFolders(t *testing.T) {
	selected := "inbox"
	folders := []Folder{{ID: "inbox", Name: "Inbox", Count: 2}, {ID: "sent", Name: "Sent"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Folders(c, folders, &selected) }), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Sent")
	if selected != "sent" || !tt.HasText("2") {
		t.Errorf("selected %q texts %q", selected, tt.Texts())
	}
	cues(t, before, sound.Select)
}

func TestMailLabels(t *testing.T) {
	selected := map[string]bool{}
	labels := []MailLabel{{ID: "work", Name: "Work", Color: ui.Color{R: 1, A: 1}}}
	tt := ui.NewTester(padded(func(c *ui.Context) { MailLabels(c, labels, &selected) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Work")
	if !selected["work"] {
		t.Errorf("labels %v", selected)
	}
	cues(t, before, sound.Toggle)
}

func TestMailThread(t *testing.T) {
	open := map[string]bool{"m2": true}
	messages := []ThreadMessage{{ID: "m1", Sender: "Ana", Body: "First", Time: "09:00"}, {ID: "m2", Sender: "Ben", Body: "Last", Time: "10:00"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { MailThread(c, messages, &open) }), 500, 300)
	settles(t, tt)
	if tt.HasText("First") || !tt.HasText("Last") {
		t.Errorf("initial texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Ana")
	if !open["m1"] || !tt.HasText("First") {
		t.Errorf("open %v texts %q", open, tt.Texts())
	}
	cues(t, before, sound.Toggle)
}

func TestMailAttachments(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		MailAttachments(c, []Attachment{{Name: "notes.pdf", Size: 2048}}, MailAttachmentsOptions{})
	}), 500, 180)
	settles(t, tt)
	if !tt.HasText("notes.pdf") || !tt.HasText("2.0 KB") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSnooze(t *testing.T) {
	var at time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if at.IsZero() {
			at = c.Now()
		}
		Snooze(c, &at, SnoozeOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Tomorrow")
	want := tt.Now().AddDate(0, 0, 1)
	if at.YearDay() != want.YearDay() || at.Hour() != 9 {
		t.Errorf("at %v want %v", at, want)
	}
	cues(t, before, sound.Select)
}

func TestSendLater(t *testing.T) {
	var at time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if at.IsZero() {
			at = c.Now()
		}
		SendLater(c, &at, SendLaterOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Tomorrow morning")
	if at.Hour() != 9 || at.YearDay() != tt.Now().AddDate(0, 0, 1).YearDay() {
		t.Errorf("at %v", at)
	}
	cues(t, before, sound.Select)
}

func TestNotification(t *testing.T) {
	opened := false
	notice := Notice{ID: "n1", Text: "New reply", Time: "now"}
	tt := ui.NewTester(padded(func(c *ui.Context) { Notification(c, notice, &opened) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "New reply")
	if !opened {
		t.Error("notification not opened")
	}
	cues(t, before, sound.Tap)
}

func TestNotificationSettings(t *testing.T) {
	settings := []NotificationSetting{{Event: "Replies", Email: true}}
	tt := ui.NewTester(padded(func(c *ui.Context) { NotificationSettings(c, &settings, NotificationSettingsOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Replies Push")
	if !settings[0].Push || !settings[0].Email {
		t.Errorf("settings %+v", settings[0])
	}
	cues(t, before, sound.Toggle)
}

func TestSignature(t *testing.T) {
	text := "Ana"
	tt := ui.NewTester(padded(func(c *ui.Context) { Signature(c, &text, SignatureOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	if err := tt.Click("Signature"); err != nil {
		t.Fatal(err)
	}
	tt.Type("!")
	tt.Frame()
	if text != "Ana!" || !tt.HasText("Ana!") {
		t.Errorf("text %q texts %q", text, tt.Texts())
	}
	cues(t, before, sound.Tick)
}

func TestOutOfOffice(t *testing.T) {
	var value AutoReply
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if value.Start.IsZero() {
			value = AutoReply{Start: c.Now(), End: c.Now().AddDate(0, 0, 1), Message: "Away"}
		}
		OutOfOffice(c, &value, OutOfOfficeOptions{})
	}), 500, 350)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Out of office")
	if !value.Enabled {
		t.Error("automatic reply disabled")
	}
	cues(t, before, sound.Toggle)
}
