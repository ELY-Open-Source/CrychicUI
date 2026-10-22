package overlay

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SessionTimeoutOptions tune a SessionTimeout.
type SessionTimeoutOptions struct {
	Title, Stay, Leave string
	// Text takes the seconds left.
	Text string
}

// SessionTimeout counts down to deadline; it closes itself at zero.
func SessionTimeout(c *ui.Context, deadline time.Time, open *bool, o SessionTimeoutOptions) (root, stay, leave ui.Element) {
	if deadline.IsZero() || open == nil {
		panic("overlay: SessionTimeout needs a deadline and open")
	}
	def(&o.Title, "Still there?")
	def(&o.Text, "You will be signed out in %d s.")
	def(&o.Stay, "Stay signed in")
	def(&o.Leave, "Sign out")
	t := c.Theme()
	root = ui.Box(c)
	now := c.Now()
	left := deadline.Sub(now)
	if *open && left <= 0 {
		*open = false
		root.SetChanged()
		kit.Log("SessionTimeout", "expired")
	}
	secs := int64((left + time.Second - 1) / time.Second)
	if *open {
		c.After(left - time.Duration(secs-1)*time.Second)
	}
	shut := func(why string) {
		*open = false
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("SessionTimeout", why)
	}
	if dialog(c, open, t.Space(90), func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			clock := ui.Icon(c, icon.Must("timer")).Size(16, 16).TextColor(t.Warning)
			clock.Rotate(12 * (clock.Loop("tick", time.Second, ui.Bounce(ui.EaseInOut)) - 0.5))
			ui.Text(c, o.Title).Bold().FontSize(t.Rem(1.1))
		})
		kit.Numeric(c, fmt.Sprintf(o.Text, max(secs, 0))).TextColor(t.TextMuted)
		ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			leave = quiet(c, o.Leave)
			if leave.Clicked() {
				shut("left")
			}
			stay = primary(c, o.Stay, t.Accent, t.AccentText)
			if stay.Clicked() {
				shut("stayed")
			}
		})
	}).Dismissed() {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "SessionTimeout"); opened {
		sound.Play(sound.Warning)
	}
	return root, stay, leave
}
