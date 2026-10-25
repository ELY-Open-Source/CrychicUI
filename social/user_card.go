package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UserCard reveals a user hover card; cue Toggle.
func UserCard(c *ui.Context, user User, open *bool) ui.Element {
	validUser("UserCard", user)
	if open == nil {
		panic("social: UserCard " + user.Name + " needs state")
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2)).AlignItems(ui.Start).Role(ui.RoleGroup)
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(user.Name).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(0.5), t.Space(2), t.Space(0.5), t.Space(0.5)).Radius(999).Checked(*open))
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.1), max(unit(b.Hovered()), unit(*open))), ui.Snappy))
		b.Children(func() {
			avatar(c, user, t.Space(6))
			ui.Text(c, user.Name).FontSize(t.Rem(0.85)).Bold().TextColor(t.Accent).Shrink(0)
		})
		if b.Clicked() {
			*open = !*open
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("UserCard", "open", *open)
		}
		c.Key("card")
		ui.Column(c).Children(func() {
			if !*open && !root.Hovered() {
				return
			}
			pop := ui.Column(c).Width(t.Space(72)).Radius(t.Radius).Border(1, t.Border).Background(t.Background).Clip().
				Shadow(0, 8, 24, 0, t.Text.Alpha(0.14)).Transition(kit.Enter)
			pop.Children(func() {
				cover := ui.Box(c).FillWidth().Height(t.Space(14)).Clip()
				cover.Children(func() {
					if user.Cover != nil {
						ui.Image(c, user.Cover).Fill().Fit(ui.Cover)
					}
				})
				if user.Cover == nil {
					cover.Background(t.Accent.Alpha(0.2))
				}
				ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(0, t.Space(3), t.Space(3), t.Space(3)).Children(func() {
					ui.Box(c).Margin(-t.Space(6), 0, 0, 0).Padding(t.Space(0.75)).Radius(999).Background(t.Background).AlignSelf(ui.Start).
						Children(func() { avatar(c, user, t.Space(12)) })
					ui.Text(c, user.Name).FontSize(t.Rem(0.95)).Bold().Wrap().FillWidth()
					kit.Muted(c, user.Handle).FontSize(t.Rem(0.78))
					if user.Bio != "" {
						ui.Text(c, user.Bio).FontSize(t.Rem(0.82)).Wrap().FillWidth()
					}
				})
			})
		})
	})
	return root
}
