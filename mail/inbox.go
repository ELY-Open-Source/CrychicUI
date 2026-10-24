package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Inbox lists unread mail and selects rows; cue Select.
func Inbox(c *ui.Context, mails []Mail, selected *string) ui.Element {
	if len(mails) == 0 || selected == nil {
		panic("mail: Inbox needs messages and selection")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList)
	root.Children(func() {
		for _, m := range mails {
			validMail(m)
			c.Key(m.ID)
			row := button(c, m.Subject).Column().FillWidth().Gap(t.Space(1)).AlignItems(ui.Start).
				Border(1, t.Border).Checked(*selected == m.ID).Role(ui.RoleRow)
			row.Background(row.AnimateColor("selected", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*selected == m.ID)), ui.Snappy))
			row.Children(func() { mailText(c, m) })
			if row.Clicked() && *selected != m.ID {
				*selected = m.ID
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Inbox", "selected")
			}
		}
	})
	return root
}
