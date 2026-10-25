package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FollowersOptions word the card; Label takes action, name.
type FollowersOptions struct{ Title, FollowBack, Following, FollowsYou, Label string }

// Followers follows back people in a list; cue Toggle.
func Followers(c *ui.Context, users []User, following *map[string]bool, o FollowersOptions) ui.Element {
	if len(users) == 0 || following == nil || *following == nil {
		panic("social: Followers needs users and state")
	}
	ids := make([]string, len(users))
	for i, user := range users {
		validUser("Followers", user)
		ids[i] = user.ID
	}
	unique("Followers", ids)
	def(&o.Title, "Followers")
	def(&o.FollowBack, "Follow back")
	def(&o.Following, "Following")
	def(&o.FollowsYou, "Follows you")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "users", o.Title, func() { kit.Chip(c, number(c, len(users)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, user := range users {
				c.Key(user.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					avatar(c, user, t.Space(10))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, user.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							kit.Muted(c, user.Handle).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
							ui.Text(c, o.FollowsYou).FontSize(t.Rem(0.65)).Bold().Shrink(0).Padding(0, t.Space(1.5)).Radius(t.Radius).
								Background(t.SurfaceHover).TextColor(t.TextMuted)
						})
					})
					on := (*following)[user.ID]
					text := o.FollowBack
					if on {
						text = o.Following
					}
					b := follow(c, fmt.Sprintf(o.Label, text, user.Name), text, "user-plus", on)
					if b.Clicked() {
						(*following)[user.ID] = !on
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("Followers", "following", "index", i, "on", !on)
					}
				})
			}
		})
	})
	return root
}
