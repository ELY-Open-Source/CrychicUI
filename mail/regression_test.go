package mail

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestRecipientIdentityDuringReorder(t *testing.T) {
	to := []string{"ana@example.com", "ben@example.com"}
	cc, bcc := []string{}, []string{}
	tt := ui.NewTester(padded(func(c *ui.Context) { Recipients(c, &to, &cc, &bcc, RecipientsOptions{}) }), 600, 250)
	settles(t, tt)
	r, _ := tt.Find("ana@example.com")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	to[0], to[1] = to[1], to[0]
	tt.Frame()
	r, _ = tt.Find("ana@example.com")
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if len(to) != 1 || to[0] != "ben@example.com" {
		t.Errorf("wrong recipient removed: %v", to)
	}
}

func TestNotificationSettingIdentityDuringReorder(t *testing.T) {
	settings := []NotificationSetting{{Event: "Replies"}, {Event: "Mentions"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { NotificationSettings(c, &settings, NotificationSettingsOptions{}) }), 600, 250)
	settles(t, tt)
	r, _ := tt.Find("Replies Push")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	settings[0], settings[1] = settings[1], settings[0]
	tt.Frame()
	r, _ = tt.Find("Replies Push")
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if !settings[1].Push || settings[0].Push {
		t.Errorf("wrong setting toggled: %+v", settings)
	}
}

func TestUnsubscribeReasonIdentityDuringReorder(t *testing.T) {
	reason := ""
	confirmed := false
	reasons := []string{"Too many", "Other"}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Unsubscribe(c, reasons, &reason, &confirmed, UnsubscribeOptions{})
	}), 500, 300)
	settles(t, tt)
	r, _ := tt.Find("Too many")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Frame()
	reasons[0], reasons[1] = reasons[1], reasons[0]
	tt.Frame()
	r, _ = tt.Find("Too many")
	tt.Release(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	if reason != "Too many" {
		t.Errorf("reason %q", reason)
	}
}

func TestMailRowArchiveConsumed(t *testing.T) {
	archived := false
	tt := ui.NewTester(padded(func(c *ui.Context) { MailRow(c, testMail(), &archived) }), 500, 220)
	settles(t, tt)
	r, _ := tt.Find("Release")
	tt.Press(r.X+5, r.Y+5)
	tt.Move(r.X+500, r.Y+5)
	tt.Release(r.X+500, r.Y+5)
	tt.Advance(3 * time.Second)
	if !archived {
		t.Fatal("message not archived")
	}
	before := sound.Recent()
	archived = false
	tt.Frame()
	tt.Advance(3 * time.Second)
	if archived || len(sound.Recent()) != len(before) {
		t.Errorf("archive replayed: archived %v", archived)
	}
}

func TestEndedCallFreezesDuration(t *testing.T) {
	ended := true
	started := time.Now().Add(-65 * time.Second)
	tt := ui.NewTester(padded(func(c *ui.Context) { CallScreen(c, "Ana", started, &ended, CallScreenOptions{}) }), 500, 250)
	tt.Advance(2 * time.Second)
	if !tt.HasText("01:05") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Advance(2 * time.Second)
	if !tt.HasText("01:05") {
		t.Errorf("ended duration advanced: %q", tt.Texts())
	}
}

func TestCallerEndedCallFreezesCurrentDuration(t *testing.T) {
	ended := false
	started := time.Now().Add(-65 * time.Second)
	tt := ui.NewTester(padded(func(c *ui.Context) { CallScreen(c, "Ana", started, &ended, CallScreenOptions{}) }), 500, 250)
	tt.Advance(2 * time.Second)
	ended = true
	tt.Frame()
	if !tt.HasText("01:07") {
		t.Fatalf("wrong frozen duration: %q", tt.Texts())
	}
	tt.Advance(2 * time.Second)
	if !tt.HasText("01:07") {
		t.Errorf("external end duration advanced: %q", tt.Texts())
	}
}

func TestLaterTodayBoundary(t *testing.T) {
	zone := time.FixedZone("test", 0)
	evening := time.Date(2026, 10, 23, 22, 0, 0, 0, zone)
	if _, ok := laterToday(evening); ok {
		t.Error("late evening offers tomorrow as later today")
	}
	morning := time.Date(2026, 10, 23, 10, 0, 0, 0, zone)
	later, ok := laterToday(morning)
	if !ok || later.Day() != morning.Day() || later.Hour() != 13 {
		t.Errorf("later %v available %v", later, ok)
	}
}

func TestMailThreadLastDefault(t *testing.T) {
	open := map[string]bool{}
	messages := []ThreadMessage{{ID: "m1", Sender: "Ana", Body: "First", Time: "09:00"}, {ID: "m2", Sender: "Ben", Body: "Last", Time: "10:00"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { MailThread(c, messages, &open) }), 500, 300)
	settles(t, tt)
	if tt.HasText("First") || !tt.HasText("Last") {
		t.Errorf("initial texts %q", tt.Texts())
	}
	click(t, tt, "Ben")
	if open["m2"] || tt.HasText("Last") {
		t.Errorf("last did not collapse: %v %q", open, tt.Texts())
	}
}
