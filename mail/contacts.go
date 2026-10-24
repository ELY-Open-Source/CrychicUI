package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Contacts selects people from a contact list; cue Select.
func Contacts(c *ui.Context, contacts []Contact, selected *string) ui.Element {
	if len(contacts) == 0 || selected == nil {
		panic("mail: Contacts needs contacts and selection")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList)
	root.Children(func() {
		for _, person := range contacts {
			validContact(person)
			c.Key(person.ID)
			row := button(c, person.Name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Checked(*selected == person.ID)
			row.Background(row.AnimateColor("selected", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*selected == person.ID)), ui.Snappy))
			row.Children(func() {
				ui.Box(c).Size(t.Space(7), t.Space(7)).Radius(999).Background(t.Accent.Alpha(0.15)).Center().Children(func() {
					ui.Text(c, initial(person.Name)).TextColor(t.Accent)
				})
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, person.Name)
					kit.Muted(c, person.Email).SingleLine().FontSize(t.Rem(0.75))
				})
			})
			if row.Clicked() && *selected != person.ID {
				*selected = person.ID
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Contacts", "selected")
			}
		}
	})
	return root
}
