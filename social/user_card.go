package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UserCard reveals a user hover card; cue Toggle.
func UserCard(c *ui.Context, user User, open *bool) ui.Element {
	validUser(user)
	if open == nil {
		panic("social: UserCard needs state")
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		b := button(c, user.Name)
		b.Children(func() { ui.Text(c, user.Name).TextColor(t.Accent) })
		if b.Clicked() {
			*open = !*open
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("UserCard", "open", *open)
		}
		if *open || root.Hovered() {
			card := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Background(t.Background).Transition(kit.Enter)
			card.Children(func() {
				avatar(c, user, t.Space(10))
				ui.Column(c).Children(func() {
					ui.Text(c, user.Name).Bold()
					kit.Muted(c, user.Handle)
					ui.Text(c, user.Bio).Wrap()
				})
			})
		}
	})
	return root
}
