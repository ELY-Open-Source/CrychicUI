package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReadingPane shows selected mail beside the list.
func ReadingPane(c *ui.Context, mails []Mail, selected *string) ui.Element {
	if len(mails) == 0 || selected == nil {
		panic("mail: ReadingPane needs messages and selection")
	}
	t := c.Theme()
	var current Mail
	found := false
	for _, m := range mails {
		validMail(m)
		if m.ID == *selected {
			current, found = m, true
		}
	}
	if !found || current.Body == "" {
		panic("mail: ReadingPane needs a selected body")
	}
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Children(func() {
			for _, m := range mails {
				c.Key(m.ID)
				row := button(c, m.Subject).Column().FillWidth().AlignItems(ui.Start).Checked(m.ID == *selected)
				row.Background(row.AnimateColor("selected", t.Surface.Mix(t.Accent.Alpha(0.12), unit(m.ID == *selected)), ui.Snappy))
				row.Children(func() {
					ui.Text(c, m.Sender).SingleLine()
					ui.Text(c, m.Subject).SingleLine()
				})
				if row.Clicked() && *selected != m.ID {
					*selected = m.ID
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("ReadingPane", "selected")
				}
			}
		})
		ui.Column(c).Grow(2).Basis(0).Gap(t.Space(2)).Key(current.ID).Transition(kit.Enter).Children(func() {
			ui.Text(c, current.Subject).Bold()
			kit.Muted(c, current.Sender+" · "+current.Time)
			ui.Text(c, current.Body).Wrap()
		})
	})
	return root
}
