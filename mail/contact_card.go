package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ContactCardOptions are the contact actions.
type ContactCardOptions struct{ Call, Message string }

// ContactCard shows contact details and actions.
func ContactCard(c *ui.Context, person Contact, action *string, o ContactCardOptions) ui.Element {
	validContact(person)
	if action == nil {
		panic("mail: ContactCard needs action state")
	}
	def(&o.Call, "Call")
	def(&o.Message, "Message")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Radius(999).Background(t.Accent.Alpha(0.15)).Center().Children(func() {
				ui.Text(c, initial(person.Name)).TextColor(t.Accent)
			})
			ui.Column(c).Children(func() {
				ui.Text(c, person.Name).Bold()
				kit.Muted(c, person.Email)
				kit.Muted(c, person.Phone)
			})
		})
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			actions := []struct {
				label, id string
			}{{o.Call, "call"}, {o.Message, "message"}}
			for _, item := range actions {
				b := button(c, item.label)
				b.Children(func() { ui.Text(c, item.label) })
				if b.Clicked() {
					*action = item.id
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("ContactCard", "action", "id", item.id)
				}
			}
		})
	})
	return root
}
