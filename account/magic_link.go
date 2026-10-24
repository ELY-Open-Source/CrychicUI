package account

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MagicLinkOptions are the magic link labels.
type MagicLinkOptions struct{ Title, Message, Resend string }

// MagicLink shows check-email state with resend; cue Send.
func MagicLink(c *ui.Context, email string, resent *int, o MagicLinkOptions) ui.Element {
	if !validEmail(email) || resent == nil || *resent < 0 {
		panic("account: MagicLink needs email and resend state")
	}
	def(&o.Title, "Check your email")
	def(&o.Message, "We sent a magic link")
	def(&o.Resend, "Resend")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, o.Title).Bold()
		ui.Text(c, o.Message)
		kit.Muted(c, email)
		b := button(c, o.Resend).Disabled(*resent == math.MaxInt)
		b.Children(func() { ui.Text(c, o.Resend) })
		if b.Clicked() && *resent < math.MaxInt {
			*resent++
			root.SetChanged()
			sound.Play(sound.Send)
			kit.Log("MagicLink", "resent", "count", *resent)
		}
	})
	return root
}
