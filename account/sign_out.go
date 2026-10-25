package account

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SignOutOptions word the card; Farewell takes the name.
type SignOutOptions struct{ Title, Farewell, SignIn string }

// SignOut shows the signed-out state with sign in; cue Tap.
func SignOut(c *ui.Context, name string, o SignOutOptions) (root, signIn ui.Element) {
	if name == "" {
		panic("account: SignOut needs a name")
	}
	def(&o.Title, "You're signed out")
	def(&o.Farewell, "See you soon, %s.")
	def(&o.SignIn, "Sign in again")
	t := c.Theme()
	root = card(c).Gap(t.Space(3)).Padding(t.Space(6), t.Space(5)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		halo := ui.Box(c).Size(t.Space(16), t.Space(16)).Center().Radius(999).Background(t.Accent.Alpha(0.1)).TextColor(t.Accent)
		kit.Arrive(halo, ui.Pose{X: -t.Space(4), Fade: 1}, 0).Children(func() {
			kit.Arrive(glyph(c, "log-out", 28), ui.Pose{X: -t.Space(3)}, 120*time.Millisecond)
		})
		ui.Text(c, o.Title).FontSize(t.Rem(1.1)).Bold().Wrap().TextAlign(ui.Center)
		kit.Muted(c, fmt.Sprintf(o.Farewell, name)).FontSize(t.Rem(0.85)).Wrap().TextAlign(ui.Center)
		signIn = primary(c, o.SignIn, "log-in")
		if signIn.Clicked() {
			signIn.SetChanged()
			root.SetChanged()
			sound.Play(sound.Tap)
			kit.Log("SignOut", "signIn")
		}
	})
	return root, signIn
}
