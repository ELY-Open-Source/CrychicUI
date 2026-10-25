package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// InboxOptions word the card; Unread formats the count.
type InboxOptions struct{ Title, Unread, Empty string }

// Inbox lists mail, unread bold; cue Select.
func Inbox(c *ui.Context, mails []Mail, selected *string, o InboxOptions) ui.Element {
	if selected == nil {
		panic("mail: Inbox needs selection")
	}
	def(&o.Title, "Inbox")
	def(&o.Unread, "%d unread")
	def(&o.Empty, "No mail here")
	ids := make([]string, len(mails))
	unread := 0
	for i, m := range mails {
		validMail("Inbox", m)
		ids[i] = m.ID
		if m.Unread {
			unread++
		}
	}
	unique("Inbox", ids)
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "inbox", o.Title, func() {
			if unread > 0 {
				kit.Chip(c, fmt.Sprintf(o.Unread, unread), t.Accent)
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5)).Role(ui.RoleList)
		list.Children(func() {
			if len(mails) == 0 {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).TextColor(t.TextMuted).Children(func() {
					glyph(c, "mail-open", 16)
					ui.Text(c, o.Empty).FontSize(t.Rem(0.8)).Shrink(0)
				})
			}
			for i, m := range mails {
				c.Key(m.ID)
				r := row(c, m.Subject, *selected == m.ID).Role(ui.RoleRow)
				kit.Stagger(c, list, r, i).Children(func() { mailFace(c, m) })
				if r.Clicked() && *selected != m.ID {
					*selected = m.ID
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Inbox", "selected", "index", i)
				}
			}
		})
	})
	return root
}
