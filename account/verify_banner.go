package account

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VerifyBannerOptions word the banner; Message takes the email.
type VerifyBannerOptions struct{ Message, Resend, Sent string }

// VerifyBanner asks to verify an email with resend; cue Send.
func VerifyBanner(c *ui.Context, email string, resent *int, o VerifyBannerOptions) ui.Element {
	if !validEmail(email) || resent == nil || *resent < 0 {
		panic(fmt.Sprintf("account: VerifyBanner %q needs an email and resend count", email))
	}
	def(&o.Message, "Verify %s to keep your account secure.")
	def(&o.Resend, "Resend email")
	def(&o.Sent, "Email sent")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(4)).
		Radius(t.Radius*1.5).Background(t.Warning.Alpha(0.1)).Border(1, t.Warning.Alpha(0.35)).Role(ui.RoleStatus)
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		icon := glyph(c, "mail-warning", 20).TextColor(t.Warning)
		kit.Arrive(icon, ui.Pose{Rotate: -15, Scale: 1.2}, 200*time.Millisecond)
		ui.Text(c, fmt.Sprintf(o.Message, email)).FontSize(t.Rem(0.85)).Grow(1).MinWidth(t.Space(40)).Wrap()
		sent := *resent > 0
		text, mark := o.Resend, "send"
		if sent {
			text, mark = o.Sent, "check"
		}
		b := kit.Press(ui.ButtonBase(c).Label(o.Resend).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).Shrink(0).
			Disabled(*resent == math.MaxInt))
		b.Background(b.AnimateColor("bg", t.Warning.Alpha(0.12).Mix(t.Warning.Alpha(0.22), unit(b.Hovered())), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.Success, unit(sent)), ui.Snappy))
		b.Children(func() {
			c.Key(sent)
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.5, Fade: 1}, 0)
				ui.Text(c, text).FontSize(t.Rem(0.8)).Bold().Shrink(0)
			})
		})
		if b.Clicked() && *resent < math.MaxInt {
			*resent++
			root.SetChanged()
			sound.Play(sound.Send)
			kit.Log("VerifyBanner", "resent", "count", *resent)
		}
	})
	return root
}
