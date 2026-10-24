package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WhoToFollowOptions are the suggestion labels.
type WhoToFollowOptions struct{ Follow, Following string }

// WhoToFollow follows suggested people; cue Toggle.
func WhoToFollow(c *ui.Context, users []User, following *map[string]bool, o WhoToFollowOptions) ui.Element {
	if len(users) == 0 || following == nil || *following == nil {
		panic("social: WhoToFollow needs users and state")
	}
	def(&o.Follow, "Follow")
	def(&o.Following, "Following")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, user := range users {
			validUser(user)
			c.Key(user.ID)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				avatar(c, user, t.Space(8))
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, user.Name)
					kit.Muted(c, user.Handle).SingleLine()
				})
				on := (*following)[user.ID]
				label := o.Follow
				if on {
					label = o.Following
				}
				b := button(c, user.ID+" "+label).Checked(on)
				b.Children(func() { ui.Text(c, label) })
				if b.Clicked() {
					(*following)[user.ID] = !on
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("WhoToFollow", "following", "on", !on)
				}
			})
		}
	})
	return root
}
