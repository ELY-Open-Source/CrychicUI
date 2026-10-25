package account

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LockoutOptions word the card.
type LockoutOptions struct{ Title, Note, Wait, Ready string }

// Lockout counts down a sign-in lock from from to until.
func Lockout(c *ui.Context, from, until time.Time, o LockoutOptions) ui.Element {
	if !until.After(from) || until.Unix() <= from.Unix() {
		panic(fmt.Sprintf("account: Lockout %q needs until at least a second after from", o.Title))
	}
	def(&o.Title, "Too many attempts")
	def(&o.Note, "For your security, sign-in is paused.")
	def(&o.Wait, "Try again in")
	def(&o.Ready, "You can try again now")
	t := c.Theme()
	now := c.Now()
	left := max(0, until.Sub(now))
	// Whole seconds from Unix times: Sub saturates past 292 years.
	secs := until.Unix() - now.Unix()
	frac := time.Duration((until.Nanosecond() - now.Nanosecond() + 1e9) % 1e9)
	if until.Nanosecond() > now.Nanosecond() {
		secs++
	}
	secs = max(0, secs)
	if left > 0 {
		if frac == 0 {
			frac = time.Second
		}
		c.After(min(left, frac+time.Millisecond))
	}
	total := until.Unix() - from.Unix()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(6), t.Space(5)).AlignItems(ui.Center).Role(ui.RoleStatus)
	root.Children(func() {
		locked := secs > 0
		ink, mark := t.Danger, "lock"
		if !locked {
			ink, mark = t.Success, "lock-open"
		}
		c.Key(locked)
		halo := ui.Box(c).Size(t.Space(16), t.Space(16)).Center().Radius(999).Background(ink.Alpha(0.1)).TextColor(ink)
		kit.Arrive(halo, ui.Pose{Scale: 0.6, Rotate: -12, Fade: 1}, 0).Children(func() { glyph(c, mark, 28) })
		ui.Text(c, o.Title).FontSize(t.Rem(1.1)).Bold().Wrap().TextAlign(ui.Center)
		kit.Muted(c, o.Note).FontSize(t.Rem(0.85)).Wrap().TextAlign(ui.Center)
		if !locked {
			c.Key("ready")
			note(c, "circle-check", o.Ready, t.Success)
			return
		}
		c.Key("wait")
		ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			kit.Muted(c, o.Wait).FontSize(t.Rem(0.78))
			parts := fmt.Sprintf("%02d:%02d", secs/60%60, secs%60)
			if secs >= 3600 {
				parts = fmt.Sprintf("%d:%02d:%02d", secs/3600, secs/60%60, secs%60)
			}
			kit.Numeric(c, parts).FontSize(t.Rem(1.8)).Bold().Font("monospace").Shrink(0)
			meter(c, float32(float64(secs)/float64(total)), t.Danger)
		})
	})
	return root
}
