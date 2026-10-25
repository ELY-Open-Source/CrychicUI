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
	validUser("ProfileHeader", user)
	if following == nil {
		panic("social: ProfileHeader " + user.Name + " needs state")
	}
	def(&o.Follow, "Follow")
	def(&o.Following, "Following")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		cover := ui.Box(c).FillWidth().Height(t.Space(24)).Clip()
		cover.Children(func() {
			if user.Cover != nil {
				ui.Image(c, user.Cover).Fill().Fit(ui.Cover)
			}
		})
		if user.Cover == nil {
			cover.Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillGradient(r, ui.LinearGradient{Angle: 135, Stops: []ui.GradientStop{
					{Color: t.Accent.Alpha(0.35), At: 0}, {Color: t.Success.Alpha(0.25), At: 1}}}, 0)
			})
		}
		c.Key("profile")
		ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(0, t.Space(4), t.Space(4), t.Space(4)).Children(func() {
			ui.Row(c).FillWidth().AlignItems(ui.End).Children(func() {
				face := ui.Box(c).Shrink(0).Margin(-t.Space(9), 0, 0, 0).Radius(999).Background(t.Background).Padding(t.Space(1))
				face.Children(func() { avatar(c, user, t.Space(18)) })
				ui.Box(c).Grow(1)
				text := o.Follow
				if *following {
					text = o.Following
				}
				b := follow(c, text, text, "user-plus", *following)
				if b.Clicked() {
					*following = !*following
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("ProfileHeader", "following", *following)
				}
			})
			ui.Column(c).FillWidth().Gap(t.Space(0.25)).Children(func() {
				ui.Text(c, user.Name).FontSize(t.Rem(1.15)).Bold().Wrap().FillWidth()
				kit.Muted(c, user.Handle).FontSize(t.Rem(0.8))
			})
			if user.Bio != "" {
				ui.Text(c, user.Bio).FontSize(t.Rem(0.88)).Wrap().FillWidth()
			}
		})
	})
	return root
}
