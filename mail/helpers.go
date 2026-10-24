// Package mail holds mail and communication components.
package mail

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Mail is one message summary and body.
type Mail struct {
	ID, Sender, Subject, Preview, Body, Time string
	Unread                                   bool
}

// Contact is one person's contact details.
type Contact struct {
	ID, Name, Email, Phone string
}

// Participant is one call participant.
type Participant struct {
	ID, Name        string
	Muted, Speaking bool
}

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func validMail(m Mail) {
	if m.ID == "" || m.Sender == "" || m.Subject == "" || m.Time == "" {
		panic("mail: message needs id, sender, subject, and time")
	}
}

func validContact(p Contact) {
	if p.ID == "" || p.Name == "" || p.Email == "" {
		panic("mail: contact needs id, name, and email")
	}
}

func initial(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && n == 0 {
		panic("mail: name must not be empty")
	}
	return string(r)
}

func button(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	return kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
}

func mailText(c *ui.Context, m Mail) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
		sender := ui.Text(c, m.Sender).Grow(1).SingleLine()
		if m.Unread {
			sender.Bold()
		}
		kit.Muted(c, m.Time).SingleLine()
	})
	subject := ui.Text(c, m.Subject).SingleLine()
	if m.Unread {
		subject.Bold()
	}
	kit.Muted(c, m.Preview).SingleLine().FontSize(t.Rem(0.8))
}

func bytes(n int64) string {
	if n < 0 {
		panic("mail: attachment size must not be negative")
	}
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	value := float64(n)
	for _, suffix := range []string{"KB", "MB", "GB", "TB"} {
		value /= unit
		if value < unit {
			return fmt.Sprintf("%.1f %s", value, suffix)
		}
	}
	return fmt.Sprintf("%.1f PB", value/unit)
}
