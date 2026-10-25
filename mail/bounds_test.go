package mail

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/egoist/mygo/ui"
)

// narrow puts a view in a 320-DIP column.
func narrow(view func(c *ui.Context)) func(c *ui.Context) {
	return padded(func(c *ui.Context) { ui.Column(c).Width(320).Children(func() { view(c) }) })
}

func inside(t *testing.T, tt *ui.Tester, name string) {
	t.Helper()
	r, ok := tt.Find(name)
	if !ok || r.W == 0 {
		t.Fatalf("%q missing: %q", name, tt.Texts())
	}
	if r.X+r.W > 330.5 {
		t.Errorf("%q ends at %.1f, past 330", name, r.X+r.W)
	}
}

func TestContactsDottedI(t *testing.T) {
	selected := ""
	people := []Contact{
		{ID: "irene", Name: "Irene Adler", Email: "irene@atlas.app"},
		{ID: "jules", Name: "Jules Martin", Email: "jules@atlas.app"},
		{ID: "irmak", Name: "ırmak Çelik", Email: "irmak@atlas.app"},
	}
	tt := ui.NewTester(padded(func(c *ui.Context) { Contacts(c, people, &selected, ContactsOptions{}) }), 500, 400)
	settles(t, tt)
	n := 0
	for _, s := range tt.Texts() {
		if s == "I" {
			n++
		}
	}
	if n != 1 {
		t.Errorf("%d I groups: %q", n, tt.Texts())
	}
}

func TestLaterTodayLocalHour(t *testing.T) {
	kathmandu := time.FixedZone("NPT", 5*3600+45*60)
	early, _ := laterToday(time.Date(2026, 10, 24, 10, 5, 0, 0, kathmandu))
	late, _ := laterToday(time.Date(2026, 10, 24, 10, 50, 0, 0, kathmandu))
	if !early.Equal(late) || early.Hour() != 13 || early.Minute() != 0 {
		t.Errorf("10:05 gives %v, 10:50 gives %v", early, late)
	}
	helsinki, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2027, 10, 31, 0, 30, 0, 0, helsinki)
	if later, _ := laterToday(now); later.Sub(now) != 150*time.Minute {
		t.Errorf("repeated hour: %v is %v after %v", later, later.Sub(now), now)
	}
}

func TestRepeatChoiceReportsOnce(t *testing.T) {
	for _, tc := range []struct {
		name, choice string
		view         func(c *ui.Context, at *time.Time) ui.Element
	}{
		{"Snooze", "Tomorrow", func(c *ui.Context, at *time.Time) ui.Element { return Snooze(c, at, SnoozeOptions{}) }},
		{"SendLater", "Tomorrow morning", func(c *ui.Context, at *time.Time) ui.Element { return SendLater(c, at, SendLaterOptions{}) }},
	} {
		var at time.Time
		changes := 0
		tt := ui.NewTester(padded(func(c *ui.Context) {
			if at.IsZero() {
				at = c.Now()
			}
			if tc.view(c, &at).Changed() {
				changes++
			}
		}), 500, 300)
		settles(t, tt)
		click(t, tt, tc.choice)
		tt.Advance(time.Second)
		click(t, tt, tc.choice)
		if changes != 1 {
			t.Errorf("%s: %d changes for one new value", tc.name, changes)
		}
	}
}

func TestChipsStayInside(t *testing.T) {
	long := "Customer support and escalation for Europe and Asia Pacific"
	selected := map[string]bool{}
	labels := []MailLabel{{ID: "work", Name: "Work", Color: ui.Color{R: 1, A: 1}}, {ID: "cs", Name: long, Color: ui.Color{B: 1, A: 1}}}
	tt := ui.NewTester(narrow(func(c *ui.Context) { MailLabels(c, labels, &selected) }), 600, 300)
	settles(t, tt)
	inside(t, tt, long)

	to := []string{"customer.support.escalations.europe.asiapacific@atlas.app"}
	cc, bcc := []string{}, []string{}
	tt = ui.NewTester(narrow(func(c *ui.Context) { Recipients(c, &to, &cc, &bcc, RecipientsOptions{}) }), 600, 300)
	settles(t, tt)
	inside(t, tt, "Remove "+to[0])
}

func TestVideoCaptionInsideTile(t *testing.T) {
	long := "Bartholomew Montgomery-Fitzwilliam"
	people := []Participant{
		{ID: "a", Name: "Maya Chen"}, {ID: "b", Name: "Daniel Okafor"}, {ID: "c", Name: "Priya Nair"},
		{ID: "d", Name: "Lucas Martin"}, {ID: "e", Name: "Sofia Rossi"}, {ID: "f", Name: long, Muted: true},
	}
	tt := ui.NewTester(narrow(func(c *ui.Context) { VideoGrid(c, people) }), 600, 400)
	settles(t, tt)
	inside(t, tt, long)
}

func TestAttachmentNamesUnique(t *testing.T) {
	panics(t, "MailAttachments", func(c *ui.Context) {
		MailAttachments(c, []Attachment{{Name: "launch.pdf", Size: 10}, {Name: "launch.pdf", Size: 20}}, MailAttachmentsOptions{})
	})
}
