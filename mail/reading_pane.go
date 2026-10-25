package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReadingPane shows selected mail beside the list.
func ReadingPane(c *ui.Context, mails []Mail, selected *string) ui.Element {
	if len(mails) == 0 || selected == nil {
		panic("mail: ReadingPane needs messages and selection")
	}
	ids := make([]string, len(mails))
	var current Mail
	found := false
	for i, m := range mails {
		validMail("ReadingPane", m)
		ids[i] = m.ID
		if m.ID == *selected {
			current, found = m, true
		}
	}
	unique("ReadingPane", ids)
	if !found || current.Body == "" {
		panic(fmt.Sprintf("mail: ReadingPane selection %q needs a message with a body", *selected))
	}
	t := c.Theme()
	root := card(c).Row().AlignItems(ui.Stretch)
	root.Children(func() {
		list := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(1.5)).
			BorderWidth(0, 1, 0, 0).BorderColor(t.Border).Role(ui.RoleList)
		list.Children(func() {
			for i, m := range mails {
				c.Key(m.ID)
				r := row(c, m.Subject, m.ID == *selected).Role(ui.RoleRow)
				kit.Stagger(c, list, r, i).Children(func() { mailFace(c, m) })
				if r.Clicked() && *selected != m.ID {
					*selected = m.ID
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("ReadingPane", "selected", "index", i)
				}
			}
		})
		ui.Column(c).Grow(1.4).Basis(0).MinWidth(0).Padding(t.Space(3)).Children(func() {
			pane := ui.Column(c).Key(current.ID).FillWidth().Gap(t.Space(2.5)).Transition(kit.Enter)
			pane.Children(func() {
				ui.Text(c, current.Subject).FontSize(t.Rem(1)).Bold().Wrap().FillWidth()
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					avatar(c, current.Sender, 8)
					ui.Text(c, current.Sender).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
					kit.Muted(c, current.Time).FontSize(t.Rem(0.72)).Shrink(0)
				})
				ui.Text(c, current.Body).FontSize(t.Rem(0.82)).Wrap().FillWidth().TextColor(t.Text.Mix(t.TextMuted, 0.25))
			})
		})
	})
	return root
}
