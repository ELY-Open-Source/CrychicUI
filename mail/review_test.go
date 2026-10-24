package mail

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestContactCardLocalizedAction(t *testing.T) {
	action := ""
	person := Contact{ID: "ana", Name: "Ana", Email: "ana@example.com"}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ContactCard(c, person, &action, ContactCardOptions{Call: "呼叫"})
	}), 500, 250)
	settles(t, tt)
	click(t, tt, "呼叫")
	if action != "call" {
		t.Errorf("action %q", action)
	}
}

func TestUnsubscribeRejectsUnknownReason(t *testing.T) {
	reason := "stale"
	confirmed, reported := false, false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		_, confirm := Unsubscribe(c, []string{"Too many"}, &reason, &confirmed, UnsubscribeOptions{})
		reported = reported || confirm.Changed()
	}), 500, 300)
	settles(t, tt)
	if err := tt.Click("Unsubscribe"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if confirmed || reported {
		t.Error("unknown reason confirmed")
	}
}

func TestEndedCallsStopUpdating(t *testing.T) {
	ended, accepted := true, true
	call := ui.NewTester(padded(func(c *ui.Context) {
		CallScreen(c, "Ana", c.Now().Add(-time.Minute), &ended, CallScreenOptions{})
	}), 500, 250)
	call.Advance(2 * time.Second)
	if call.FramePending() {
		t.Error("ended call timer still running")
	}
	incoming := ui.NewTester(padded(func(c *ui.Context) {
		declined := false
		IncomingCall(c, "Ana", &accepted, &declined, IncomingCallOptions{})
	}), 500, 300)
	incoming.Advance(2 * time.Second)
	if incoming.FramePending() {
		t.Error("answered call ring still running")
	}
}

func TestSendLaterNextMonday(t *testing.T) {
	var at time.Time
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if at.IsZero() {
			at = c.Now()
		}
		SendLater(c, &at, SendLaterOptions{})
	}), 500, 300)
	settles(t, tt)
	click(t, tt, "Monday morning")
	if at.Weekday() != time.Monday || !at.After(tt.Now()) || at.Hour() != 9 {
		t.Errorf("at %v", at)
	}
}

func TestAttachmentSizeBoundaries(t *testing.T) {
	cases := map[int64]string{0: "0 B", 1023: "1023 B", 1024: "1.0 KB", 1048576: "1.0 MB"}
	for size, want := range cases {
		if got := bytes(size); got != want {
			t.Errorf("bytes(%d) = %q, want %q", size, got, want)
		}
	}
}
