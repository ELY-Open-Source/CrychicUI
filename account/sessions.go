package account

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Session is one signed-in session; Location and Icon are optional.
type Session struct {
	ID, Device, Location, Last, Icon string
	Current                          bool
}

// SessionsOptions word the card; Label takes action and device.
type SessionsOptions struct{ Title, SignOut, Current, Empty, Label string }

// Sessions signs out other sessions; cue Close.
func Sessions(c *ui.Context, sessions *[]Session, o SessionsOptions) ui.Element {
	if sessions == nil {
		panic("account: Sessions needs sessions")
	}
	ids := make([]string, len(*sessions))
	for i, s := range *sessions {
		if s.Device == "" || s.Last == "" {
			panic(fmt.Sprintf("account: Sessions %q needs a device and last seen", s.ID))
		}
		ids[i] = s.ID
	}
	unique("Sessions", ids)
	def(&o.Title, "Active sessions")
	def(&o.SignOut, "Sign out")
	def(&o.Current, "This device")
	def(&o.Empty, "No other sessions")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	remove := -1
	root.Children(func() {
		kit.CardHead(c, "monitor-smartphone", o.Title, func() { kit.Chip(c, fmt.Sprint(len(*sessions)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(*sessions) == 0 {
				empty(c, "monitor-smartphone", o.Empty)
			}
			for i, s := range *sessions {
				c.Key(s.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Transition(rowExit)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					mark := s.Icon
					def(&mark, "monitor")
					ink := t.TextMuted
					if s.Current {
						ink = t.Success
					}
					tile(c, mark, ink, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, s.Device).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							if s.Location != "" {
								kit.Muted(c, s.Location).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
								dot(c)
							}
							kit.Muted(c, s.Last).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
						})
					})
					if s.Current {
						kit.Chip(c, o.Current, t.Success)
						return
					}
					b := outline(c, fmt.Sprintf(o.Label, o.SignOut, s.Device), o.SignOut, "log-out", t.Danger)
					if b.Clicked() {
						remove = i
					}
				})
			}
		})
	})
	if remove >= 0 {
		*sessions = slices.Delete(slices.Clone(*sessions), remove, remove+1)
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("Sessions", "signedOut", "count", len(*sessions))
	}
	return root
}
