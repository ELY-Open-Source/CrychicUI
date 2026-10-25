package account

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PasskeyOptions word the prompt.
type PasskeyOptions struct{ Title, Description, Add, Added string }

// Passkey adds a passkey; cue Success.
func Passkey(c *ui.Context, added *bool, o PasskeyOptions) ui.Element {
	if added == nil {
		panic("account: Passkey needs state")
	}
	def(&o.Title, "Add a passkey")
	def(&o.Description, "Sign in with your fingerprint, face or screen lock. No password to remember.")
	def(&o.Add, "Add passkey")
	def(&o.Added, "Added")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(6), t.Space(5)).AlignItems(ui.Center)
	root.Children(func() {
		ink := t.Accent
		if *added {
			ink = t.Success
		}
		c.Key(*added)
		halo := ui.Box(c).Size(t.Space(18), t.Space(18)).Center().Radius(999).Background(ink.Alpha(0.08))
		kit.Arrive(halo, ui.Pose{Scale: 0.7, Fade: 1}, 0).Children(func() {
			disc := ui.Box(c).Size(t.Space(12), t.Space(12)).Center().Radius(999).Background(ink.Alpha(0.16)).TextColor(ink)
			kit.Arrive(disc, ui.Pose{Scale: 0.5}, 120*time.Millisecond).Children(func() {
				mark := "fingerprint"
				if *added {
					mark = "circle-check"
				}
				glyph(c, mark, 28)
			})
		})
		ui.Text(c, o.Title).FontSize(t.Rem(1.1)).Bold().Shrink(0)
		kit.Muted(c, o.Description).FontSize(t.Rem(0.85)).Wrap().TextAlign(ui.Center).MaxWidth(t.Space(80))
		label, mark := o.Add, "key-round"
		if *added {
			label, mark = o.Added, "check"
		}
		b := primary(c, label, mark).Disabled(*added)
		if b.Clicked() && !*added {
			*added = true
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("Passkey", "added")
		}
	})
	return root
}
