package account

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MagicLinkOptions word the card; Again takes the resend count.
type MagicLinkOptions struct{ Title, Message, Resend, Again string }

// MagicLink shows check-email state with resend; cue Send.
func MagicLink(c *ui.Context, email string, resent *int, o MagicLinkOptions) ui.Element {
	if !validEmail(email) || resent == nil || *resent < 0 {
		panic(fmt.Sprintf("account: MagicLink %q needs an email and resend state", email))
	}
	def(&o.Title, "Check your email")
	def(&o.Message, "We sent a sign-in link to")
	def(&o.Resend, "Resend")
	def(&o.Again, "Sent again ×%d")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(6), t.Space(5)).AlignItems(ui.Center)
	root.Children(func() {
		c.Key(*resent)
		env := ui.Box(c).Size(t.Space(16), t.Space(16)).Center().Radius(999).Background(t.Accent.Alpha(0.1)).TextColor(t.Accent)
		kit.Arrive(env, ui.Pose{Y: -t.Space(3), Scale: 0.7, Rotate: -8, Fade: 1}, 0).Children(func() {
			kit.Arrive(glyph(c, "mail-check", 32), ui.Pose{Scale: 0.5}, 150*time.Millisecond)
		})
		ui.Text(c, o.Title).FontSize(t.Rem(1.1)).Bold().Shrink(0)
		ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			kit.Muted(c, o.Message).FontSize(t.Rem(0.85)).Wrap().TextAlign(ui.Center)
			ui.Text(c, email).FontSize(t.Rem(0.88)).Bold().TextColor(t.Accent).Wrap().TextAlign(ui.Center)
		})
		b := outline(c, o.Resend, o.Resend, "rotate-cw", t.Text).Disabled(*resent == math.MaxInt)
		if b.Clicked() && *resent < math.MaxInt {
			*resent++
			root.SetChanged()
			sound.Play(sound.Send)
			kit.Log("MagicLink", "resent", "count", *resent)
		}
		if *resent > 0 {
			kit.Muted(c, fmt.Sprintf(o.Again, *resent)).FontSize(t.Rem(0.75))
		}
	})
	return root
}
