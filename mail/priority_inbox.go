package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PriorityInboxOptions are the priority section labels.
type PriorityInboxOptions struct{ Important, Other string }

// PriorityInbox groups important and other messages.
func PriorityInbox(c *ui.Context, important, other []Mail, o PriorityInboxOptions) ui.Element {
	if len(important) == 0 && len(other) == 0 {
		panic("mail: PriorityInbox needs messages")
	}
	def(&o.Important, "Important")
	def(&o.Other, "Other")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		groups := []struct {
			title string
			mails []Mail
		}{{o.Important, important}, {o.Other, other}}
		for i, group := range groups {
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Text(c, group.title).Bold()
				for _, m := range group.mails {
					validMail(m)
					c.Key(m.ID)
					ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).Children(func() {
						mailText(c, m)
					})
				}
			})
		}
	})
	return root
}
