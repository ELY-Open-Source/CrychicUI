package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThreadMessage is one message in a mail thread.
type ThreadMessage struct {
	ID, Sender, Body, Time string
}

// MailThread toggles collapsed messages; cue Toggle.
func MailThread(c *ui.Context, messages []ThreadMessage, open *map[string]bool) ui.Element {
	if len(messages) == 0 || open == nil || *open == nil {
		panic("mail: MailThread needs messages and state")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, m := range messages {
			if m.ID == "" || m.Sender == "" || m.Body == "" || m.Time == "" {
				panic("mail: invalid thread message")
			}
			c.Key(m.ID)
			expanded, set := (*open)[m.ID]
			if !set && i == len(messages)-1 {
				expanded = true
			}
			card := ui.Column(c).FillWidth().Radius(t.Radius).Border(1, t.Border)
			card.Children(func() {
				header := button(c, m.Sender).Row().FillWidth().AlignItems(ui.Center).Checked(expanded)
				header.Children(func() {
					ui.Text(c, m.Sender).Bold().Grow(1)
					kit.Muted(c, m.Time)
				})
				if header.Clicked() {
					(*open)[m.ID] = !expanded
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("MailThread", "toggled", "open", (*open)[m.ID])
				}
				if expanded {
					ui.Text(c, m.Body).Wrap().Padding(t.Space(2)).Transition(kit.Enter)
				}
			})
		}
	})
	return root
}
