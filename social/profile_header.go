package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProfileHeaderOptions are the follow labels.
type ProfileHeaderOptions struct{ Follow, Following string }

// ProfileHeader shows a profile and follows; cue Toggle.
func ProfileHeader(c *ui.Context, user User, following *bool, o ProfileHeaderOptions) ui.Element {
	validUser(user)
	if following == nil {
		panic("social: ProfileHeader needs state")
	}
	def(&o.Follow, "Follow")
	def(&o.Following, "Following")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		if user.Cover != nil {
			ui.Image(c, user.Cover).FillWidth().AspectRatio(3).Fit(ui.Cover).Radius(t.Radius).Clip()
		}
		c.Key("profile")
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			avatar(c, user, t.Space(14))
			ui.Column(c).Grow(1).Children(func() {
				ui.Text(c, user.Name).Bold()
				kit.Muted(c, user.Handle)
			})
			label := o.Follow
			if *following {
				label = o.Following
			}
			b := button(c, label).Checked(*following)
			b.Children(func() { ui.Text(c, label) })
			if b.Clicked() {
				*following = !*following
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("ProfileHeader", "following", *following)
			}
		})
		ui.Text(c, user.Bio).Wrap()
	})
	return root
}
