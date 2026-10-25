package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThreadMessage is one message in a mail thread.
type ThreadMessage struct {
	ID, Sender, Body, Time string
}

// MailThread toggles collapsed messages; the last starts open; cue Toggle.
func MailThread(c *ui.Context, messages []ThreadMessage, open *map[string]bool) ui.Element {
	if len(messages) == 0 || open == nil || *open == nil {
		panic("mail: MailThread needs messages and state")
	}
	ids := make([]string, len(messages))
	for i, m := range messages {
		if m.Sender == "" || m.Body == "" || m.Time == "" {
			panic(fmt.Sprintf("mail: MailThread message %q needs sender, body, and time", m.ID))
		}
		ids[i] = m.ID
	}
	unique("MailThread", ids)
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		for i, m := range messages {
			c.Key(m.ID)
			expanded, set := (*open)[m.ID]
			if !set && i == len(messages)-1 {
				expanded = true
			}
			item := ui.Column(c).FillWidth()
			if i < len(messages)-1 {
				item.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
			}
			kit.Stagger(c, root, item, i).Children(func() {
				header := kit.Press(ui.ButtonBase(c).Label(m.Sender).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(2.5), t.Space(3)).Checked(expanded))
				header.Background(header.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(header.Hovered())), ui.Snappy))
				header.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						avatar(c, m.Sender, 8)
						ui.Text(c, m.Sender).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
						kit.Muted(c, m.Time).FontSize(t.Rem(0.72)).Shrink(0)
					})
				})
				if header.Clicked() {
					(*open)[m.ID] = !expanded
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("MailThread", "toggled", "index", i, "open", (*open)[m.ID])
				}
				layout.Collapse(c, expanded, func() {
					ui.Text(c, m.Body).FontSize(t.Rem(0.85)).Wrap().FillWidth().Padding(0, t.Space(3), t.Space(3), t.Space(13.5))
				})
			})
		}
	})
	return root
}
