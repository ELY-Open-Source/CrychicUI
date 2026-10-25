package mail

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// exact reports a text shown verbatim, not as a substring.
func exact(tt *ui.Tester, s string) bool { return slices.Contains(tt.Texts(), s) }

func panics(t *testing.T, want string, view func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), want) {
			t.Errorf("panic %v, want %q", r, want)
		}
	}()
	ui.Render(view, 600, 600, 1)
}

// swapUnderPress presses name, reorders the data, then releases there.
func swapUnderPress(t *testing.T, tt *ui.Tester, name string, reorder func()) {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok {
		t.Fatalf("%q missing", name)
	}
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	reorder()
	tt.Frame()
	tt.Release(x, y)
	tt.Frame()
}

func twoMails() []Mail {
	a := Mail{ID: "a", Sender: "Maya Chen", Subject: "Launch", Time: "9:42", Unread: true}
	b := Mail{ID: "b", Sender: "Daniel Okafor", Subject: "Crash fix", Time: "9:10"}
	return []Mail{a, b}
}

func TestInboxIdentityAndEmpty(t *testing.T) {
	mails, selected := twoMails(), ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Inbox(c, mails, &selected, InboxOptions{}) }), 500, 400)
	settles(t, tt)
	if !exact(tt, "1 unread") {
		t.Errorf("texts %q", tt.Texts())
	}
	swapUnderPress(t, tt, "Launch", func() { mails[0], mails[1] = mails[1], mails[0] })
	if selected == "b" {
		t.Error("release picked the mail that slid in")
	}
	mails = nil
	tt.Frame()
	if !exact(tt, "No mail here") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, `repeated "a"`, func(c *ui.Context) {
		Inbox(c, []Mail{twoMails()[0], twoMails()[0]}, &selected, InboxOptions{})
	})
}

func TestLabelsAndTemplatesKeepIdentity(t *testing.T) {
	labels := []MailLabel{{ID: "launch", Name: "Launch"}, {ID: "bugs", Name: "Bugs"}}
	on := map[string]bool{}
	tt := ui.NewTester(padded(func(c *ui.Context) { MailLabels(c, labels, &on) }), 500, 200)
	settles(t, tt)
	swapUnderPress(t, tt, "Launch", func() { labels[0], labels[1] = labels[1], labels[0] })
	if on["bugs"] {
		t.Error("toggled the label that slid in")
	}
	replies := []SavedReply{{ID: "a", Title: "Thanks", Body: "x"}, {ID: "b", Title: "Later", Body: "y"}}
	selected := -1
	tt = ui.NewTester(padded(func(c *ui.Context) { MessageTemplates(c, replies, &selected, MessageTemplatesOptions{}) }), 500, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Thanks", func() { replies[0], replies[1] = replies[1], replies[0] })
	if selected >= 0 && replies[selected].ID != "a" {
		t.Errorf("selected the reply that slid in: %d", selected)
	}
}

func TestAttachmentsSummary(t *testing.T) {
	files := []Attachment{{Name: "notes.pdf", Size: 2048}}
	tt := ui.NewTester(padded(func(c *ui.Context) { MailAttachments(c, files, MailAttachmentsOptions{}) }), 500, 200)
	tt.Frame()
	if !exact(tt, "1 attachment · 2.0 KB") {
		t.Errorf("texts %q", tt.Texts())
	}
	files = append(files, Attachment{Name: "map.png", Size: 1024})
	tt.Frame()
	if !exact(tt, "2 attachments · 3.0 KB") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, "invalid file", func(c *ui.Context) {
		MailAttachments(c, []Attachment{{Name: "a", Size: math.MaxInt64}, {Name: "b", Size: 1}}, MailAttachmentsOptions{})
	})
}

func TestLaterTodayHoldsWithinTheHour(t *testing.T) {
	zone := time.FixedZone("test", 0)
	a, okA := laterToday(time.Date(2026, 10, 23, 10, 5, 0, 0, zone))
	b, okB := laterToday(time.Date(2026, 10, 23, 10, 50, 0, 0, zone))
	if !okA || !okB || !a.Equal(b) || a.Hour() != 13 || a.Minute() != 0 {
		t.Errorf("later %v and %v differ within one hour", a, b)
	}
}

func TestSendLaterChoicesDiffer(t *testing.T) {
	var at time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if at.IsZero() {
			at = c.Now()
		}
		SendLater(c, &at, SendLaterOptions{})
	}), 500, 400)
	settles(t, tt)
	click(t, tt, "Monday morning")
	if !at.Equal(mondayAfter(tt.Now())) {
		t.Errorf("at %v", at)
	}
	sunday := time.Date(2026, 10, 25, 20, 0, 0, 0, time.UTC)
	for _, now := range []time.Time{sunday, sunday.AddDate(0, 0, 1), sunday.AddDate(0, 0, 6)} {
		m := mondayAfter(now)
		if m.Weekday() != time.Monday || m.Hour() != 9 || m.Sub(now) <= 24*time.Hour {
			t.Errorf("from %v, Monday choice %v is not past tomorrow", now.Weekday(), m)
		}
	}
}

func TestOutOfOfficeFoldsFields(t *testing.T) {
	var value AutoReply
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if value.Start.IsZero() {
			value = AutoReply{Start: c.Now(), End: c.Now().AddDate(0, 0, 1), Message: "Away until Monday"}
		}
		OutOfOffice(c, &value, OutOfOfficeOptions{})
	}), 500, 400)
	settles(t, tt)
	if exact(tt, "Start") || !exact(tt, "Automatic replies are off") {
		t.Errorf("texts %q", tt.Texts())
	}
	click(t, tt, "Out of office")
	settles(t, tt)
	if !exact(tt, "Start") || !tt.HasText("Replying") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestContactsLetters(t *testing.T) {
	people := []Contact{{ID: "s", Name: "Sofia Rossi", Email: "s@x.io"}, {ID: "d", Name: "Daniel Okafor", Email: "d@x.io"}}
	selected := ""
	tt := ui.NewTester(padded(func(c *ui.Context) { Contacts(c, people, &selected, ContactsOptions{}) }), 500, 400)
	settles(t, tt)
	d, _ := tt.Find("D")
	s, _ := tt.Find("S")
	if d.W == 0 || s.Y <= d.Y {
		t.Errorf("letters not sorted: D %v S %v", d, s)
	}
	people = nil
	tt.Frame()
	if !exact(tt, "No contacts") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCallLoopsStop(t *testing.T) {
	sharing := true
	tt := ui.NewTester(padded(func(c *ui.Context) { ScreenShare(c, &sharing, ScreenShareOptions{}) }), 500, 200)
	tt.Advance(time.Second)
	click(t, tt, "Stop")
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Error("stopped share still animates")
	}
	playing := 0
	messages := []VoiceMessage{{ID: "v", Sender: "Daniel Okafor", Time: "10:24", Transcript: "Beta is up"}}
	tt = ui.NewTester(padded(func(c *ui.Context) { Voicemail(c, messages, &playing, VoicemailOptions{}) }), 500, 400)
	tt.Advance(time.Second)
	click(t, tt, "Pause")
	tt.Advance(2 * time.Second)
	if tt.FramePending() || tt.HasText("Beta is up") {
		t.Errorf("paused voicemail still plays: texts %q", tt.Texts())
	}
}

func TestClockPastAnHour(t *testing.T) {
	if got := clock(time.Hour + 65*time.Second); got != "1:01:05" {
		t.Errorf("clock %q", got)
	}
	if got := clock(65 * time.Second); got != "01:05" {
		t.Errorf("clock %q", got)
	}
}

func TestUnsubscribeDoneNote(t *testing.T) {
	reason, confirmed := "", false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Unsubscribe(c, []string{"Too many", "Never signed up"}, &reason, &confirmed, UnsubscribeOptions{})
	}), 500, 500)
	settles(t, tt)
	click(t, tt, "Never signed up")
	click(t, tt, "Unsubscribe")
	if !confirmed || !exact(tt, "You are unsubscribed") {
		t.Errorf("confirmed %v texts %q", confirmed, tt.Texts())
	}
	panics(t, `repeated "Other"`, func(c *ui.Context) {
		Unsubscribe(c, []string{"Other", "Other"}, &reason, &confirmed, UnsubscribeOptions{})
	})
}

func TestPriorityInboxEmptyAndRepeats(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { PriorityInbox(c, twoMails(), nil, PriorityInboxOptions{}) }), 500, 500)
	tt.Frame()
	if !exact(tt, "Nothing here") {
		t.Errorf("texts %q", tt.Texts())
	}
	panics(t, `repeated "a"`, func(c *ui.Context) {
		PriorityInbox(c, twoMails()[:1], twoMails()[:1], PriorityInboxOptions{})
	})
}

func TestReadingPaneKeepsIdentity(t *testing.T) {
	mails := append(twoMails(), Mail{ID: "c", Sender: "Priya Nair", Subject: "Copy pass", Time: "Mon"})
	mails[0].Body = "Three items left"
	selected := "a"
	tt := ui.NewTester(padded(func(c *ui.Context) { ReadingPane(c, mails, &selected) }), 700, 400)
	settles(t, tt)
	swapUnderPress(t, tt, "Crash fix", func() { mails[1], mails[2] = mails[2], mails[1] })
	if selected == "c" {
		t.Error("release selected the mail that slid in")
	}
	panics(t, `selection "zz"`, func(c *ui.Context) {
		gone := "zz"
		ReadingPane(c, mails, &gone)
	})
}

func TestRecipientsLabelsAndRepeats(t *testing.T) {
	to, cc, bcc := []string{"maya@atlas.app"}, []string{}, []string{}
	tt := ui.NewTester(padded(func(c *ui.Context) { Recipients(c, &to, &cc, &bcc, RecipientsOptions{}) }), 600, 300)
	settles(t, tt)
	click(t, tt, "Remove maya@atlas.app")
	if len(to) != 0 {
		t.Errorf("to %v", to)
	}
	panics(t, `Cc has an invalid address "nope"`, func(c *ui.Context) {
		bad := []string{"nope"}
		Recipients(c, &to, &bad, &bcc, RecipientsOptions{})
	})
}
