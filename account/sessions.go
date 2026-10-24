package account

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Session is one signed-in session.
type Session struct {
	ID, Device, Location, Last string
	Current                    bool
}

// SessionsOptions are the session labels.
type SessionsOptions struct{ SignOut, Current, Empty string }

// Sessions signs out sessions; cue Close.
func Sessions(c *ui.Context, sessions *[]Session, o SessionsOptions) ui.Element {
	if sessions == nil {
		panic("account: Sessions needs sessions")
	}
	def(&o.SignOut, "Sign out")
	def(&o.Current, "Current")
	def(&o.Empty, "No sessions")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	remove := -1
	root.Children(func() {
		if len(*sessions) == 0 {
			kit.Muted(c, o.Empty)
		}
		for i, session := range *sessions {
			if session.ID == "" || session.Device == "" || session.Last == "" {
				panic("account: invalid session")
			}
			c.Key(session.ID)
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, session.Device)
					kit.Muted(c, session.Location+" · "+session.Last)
				})
				ui.Box(c).Children(func() {
					if session.Current {
						kit.Muted(c, o.Current)
					}
				})
				b := button(c, session.ID+" "+o.SignOut)
				b.Children(func() { ui.Text(c, o.SignOut) })
				if b.Clicked() {
					remove = i
				}
			})
		}
	})
	if remove >= 0 {
		*sessions = slices.Delete(*sessions, remove, remove+1)
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("Sessions", "signedOut", "count", len(*sessions))
	}
	return root
}
