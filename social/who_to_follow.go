package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WhoToFollowOptions word the card; Label takes action, name.
type WhoToFollowOptions struct{ Title, Follow, Following, Label string }

// WhoToFollow follows suggested people; cue Toggle.
func WhoToFollow(c *ui.Context, users []User, following *map[string]bool, o WhoToFollowOptions) ui.Element {
	if len(users) == 0 || following == nil || *following == nil {
		panic("social: WhoToFollow needs users and state")
	}
	ids := make([]string, len(users))
	for i, user := range users {
		validUser("WhoToFollow", user)
		ids[i] = user.ID
	}
	unique("WhoToFollow", ids)
	def(&o.Title, "Who to follow")
	def(&o.Follow, "Follow")
	def(&o.Following, "Following")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "user-plus", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, user := range users {
				c.Key(user.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					avatar(c, user, t.Space(10))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						ui.Text(c, user.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, user.Handle).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
						if user.Bio != "" {
							ui.Text(c, user.Bio).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
						}
					})
					on := (*following)[user.ID]
					text := o.Follow
					if on {
						text = o.Following
					}
					b := follow(c, fmt.Sprintf(o.Label, text, user.Name), text, on)
					if b.Clicked() {
						(*following)[user.ID] = !on
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("WhoToFollow", "following", "index", i, "on", !on)
					}
				})
			}
		})
	})
	return root
}
