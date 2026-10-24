package gallery

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestMailDemosLog(t *testing.T) {
	press := func(name string) func(tt *ui.Tester) error {
		return func(tt *ui.Tester) error { return tt.Click(name) }
	}
	acts := map[string]struct {
		do   func(tt *ui.Tester) error
		want string
	}{
		"Inbox": {press("Release"), "selected m1"},
		"MailRow": {func(tt *ui.Tester) error {
			r, ok := tt.Find("Release")
			if !ok {
				return fmt.Errorf("subject not found: %q", tt.Texts())
			}
			tt.Press(r.X+5, r.Y+5)
			tt.Move(r.X+600, r.Y+5)
			tt.Release(r.X+600, r.Y+5)
			tt.Advance(3 * time.Second)
			return nil
		}, "archived true"},
		"ComposeWindow":        {press("Send"), "sent"},
		"Recipients":           {press("ana@example.com"), "recipients 1"},
		"Folders":              {press("Sent"), "folder sent"},
		"MailLabels":           {press("Work"), "labels 1"},
		"MailThread":           {press("Ana"), "open 2"},
		"Snooze":               {press("Tomorrow"), "snoozed"},
		"SendLater":            {press("Tomorrow morning"), "scheduled"},
		"Notification":         {press("New reply"), "opened true"},
		"NotificationSettings": {press("Replies Push"), "push true"},
		"Signature": {func(tt *ui.Tester) error {
			if err := tt.Click("Signature"); err != nil {
				return err
			}
			tt.Type("!")
			return nil
		}, "runes 4"},
		"OutOfOffice":      {press("Out of office"), "away true"},
		"Contacts":         {press("Ana"), "contact ana"},
		"ContactCard":      {press("Call"), "action call"},
		"CallScreen":       {press("End"), "ended true"},
		"IncomingCall":     {press("Accept"), "accepted true"},
		"CallControls":     {press("Mute"), "muted true"},
		"ScreenShare":      {press("Stop"), "sharing false"},
		"Voicemail":        {press("Play"), "playing 0"},
		"MessageTemplates": {press("Thanks"), "reply 0"},
		"Unsubscribe": {func(tt *ui.Tester) error {
			if err := tt.Click("Too many"); err != nil {
				return err
			}
			tt.Frame()
			return tt.Click("Unsubscribe")
		}, "confirmed"},
		"Spam":        {press("Not spam"), "spam false"},
		"MailSearch":  {press("Has attachment"), "attachments true"},
		"UndoSend":    {press("Undo"), "visible false"},
		"ReadingPane": {press("Second"), "reading m2"},
	}
	quiet := map[string]bool{"MailAttachments": true, "VideoGrid": true, "Participants": true, "PriorityInbox": true}
	seen := 0
	for _, e := range Entries() {
		if e.Category != "mail" {
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
		tt.Advance(2 * time.Second)
		if err := act.do(tt); err != nil {
			t.Errorf("%s: %v", e.Name, err)
			continue
		}
		tt.Frame()
		if !strings.Contains(strings.Join(logs, "\n"), act.want) {
			t.Errorf("%s logs %q, want %q", e.Name, logs, act.want)
		}
	}
	if seen != 30 {
		t.Errorf("checked %d mail demos", seen)
	}
}
