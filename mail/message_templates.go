package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SavedReply is one reusable message template.
type SavedReply struct {
	ID, Title, Body string
}

// MessageTemplatesOptions word the card.
type MessageTemplatesOptions struct{ Title string }

// MessageTemplates selects saved replies; cue Select.
func MessageTemplates(c *ui.Context, replies []SavedReply, selected *int, o MessageTemplatesOptions) ui.Element {
	if len(replies) == 0 || selected == nil || *selected < -1 || *selected >= len(replies) {
		panic("mail: MessageTemplates needs replies and a selection in range")
	}
	ids := make([]string, len(replies))
	for i, r := range replies {
		if r.Title == "" || r.Body == "" {
			panic(fmt.Sprintf("mail: MessageTemplates reply %q needs a title and body", r.ID))
		}
		ids[i] = r.ID
	}
	unique("MessageTemplates", ids)
	def(&o.Title, "Saved replies")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "reply", o.Title, func() { kit.Chip(c, fmt.Sprint(len(replies)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2))
		list.Children(func() {
			for i, reply := range replies {
				c.Key(reply.ID)
				on := *selected == i
				b := row(c, reply.Title, on)
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start).Children(func() {
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.75)).Children(func() {
							ui.Text(c, reply.Title).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
							kit.Muted(c, reply.Body).FontSize(t.Rem(0.75)).Wrap().FillWidth()
						})
						c.Key(on)
						if on {
							badge := ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText)
							kit.Arrive(badge, ui.Pose{Scale: 0.4, Fade: 1}, 0).Children(func() { glyph(c, "check", 12) })
						}
					})
				})
				if b.Clicked() && !on {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("MessageTemplates", "selected", "index", i)
				}
			}
		})
	})
	return root
}
