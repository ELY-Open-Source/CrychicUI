package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FollowersOptions are the follower labels.
type FollowersOptions struct{ FollowBack, Following string }

// Followers follows back people in a list; cue Toggle.
func Followers(c *ui.Context, users []User, following *map[string]bool, o FollowersOptions) ui.Element {
	if len(users) == 0 || following == nil || *following == nil {
		panic("social: Followers needs users and state")
	}
	def(&o.FollowBack, "Follow back")
	def(&o.Following, "Following")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList)
	root.Children(func() {
		for _, user := range users {
			validUser(user)
			c.Key(user.ID)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				avatar(c, user, t.Space(8))
				ui.Text(c, user.Name).Grow(1)
				on := (*following)[user.ID]
				label := o.FollowBack
				if on {
					label = o.Following
				}
				b := button(c, user.ID+" "+label).Checked(on)
				b.Children(func() { ui.Text(c, label) })
				if b.Clicked() {
					(*following)[user.ID] = !on
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Followers", "following", "on", !on)
				}
			})
		}
	})
	return root
}
