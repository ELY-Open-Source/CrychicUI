package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SavedReply is one reusable message template.
type SavedReply struct {
	ID, Title, Body string
}

// MessageTemplates selects saved replies; cue Select.
func MessageTemplates(c *ui.Context, replies []SavedReply, selected *int) ui.Element {
	if len(replies) == 0 || selected == nil || *selected < -1 || *selected >= len(replies) {
		panic("mail: MessageTemplates needs replies and selection")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, reply := range replies {
			if reply.ID == "" || reply.Title == "" || reply.Body == "" {
				panic("mail: invalid saved reply")
			}
			c.Key(reply.ID)
			card := button(c, reply.Title).Column().FillWidth().Gap(t.Space(1)).AlignItems(ui.Start).
				Border(1, t.Border).Checked(*selected == i)
			card.Background(card.AnimateColor("selected", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*selected == i)), ui.Snappy))
			card.Children(func() {
				ui.Text(c, reply.Title).Bold()
				ui.Text(c, reply.Body).Wrap()
			})
			if card.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("MessageTemplates", "selected", "index", i)
			}
		}
	})
	return root
}
