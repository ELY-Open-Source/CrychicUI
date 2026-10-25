package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PriorityInboxOptions word the sections.
type PriorityInboxOptions struct{ Important, Other, Empty string }

// PriorityInbox groups important and other messages.
func PriorityInbox(c *ui.Context, important, other []Mail, o PriorityInboxOptions) ui.Element {
	def(&o.Important, "Important")
	def(&o.Other, "Other")
	def(&o.Empty, "Nothing here")
	var ids []string
	for _, m := range append(append([]Mail{}, important...), other...) {
		validMail("PriorityInbox", m)
		ids = append(ids, m.ID)
	}
	unique("PriorityInbox", ids)
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		groups := []struct {
			title, mark string
			ink         ui.Color
			mails       []Mail
		}{{o.Important, "star", t.Warning, important}, {o.Other, "inbox", t.TextMuted, other}}
		for g, group := range groups {
			c.Key(g)
			section := ui.Column(c).FillWidth().Padding(t.Space(1.5))
			if g == 0 {
				section.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
			}
			section.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(1.5)).Children(func() {
					glyph(c, group.mark, 14).TextColor(group.ink)
					ui.Text(c, group.title).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
					kit.Chip(c, fmt.Sprint(len(group.mails)), group.ink)
				})
				if len(group.mails) == 0 {
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.8)).Padding(t.Space(1.5), t.Space(2.5))
				}
				list := ui.Column(c).FillWidth().Gap(t.Space(0.5))
				list.Children(func() {
					for i, m := range group.mails {
						c.Key(m.ID)
						line := ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius)
						kit.Stagger(c, list, line, i).Children(func() { mailFace(c, m) })
					}
				})
			})
		}
	})
	return root
}
