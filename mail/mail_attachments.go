package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Attachment is one attached file.
type Attachment struct {
	Name string
	Size int64
}

// MailAttachments shows attachment cards with sizes.
func MailAttachments(c *ui.Context, attachments []Attachment) ui.Element {
	if len(attachments) == 0 {
		panic("mail: MailAttachments needs attachments")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, attachment := range attachments {
			if attachment.Name == "" || attachment.Size < 0 {
				panic("mail: invalid attachment")
			}
			c.Key(i)
			card := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
			kit.Arrive(card, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			card.Children(func() {
				ui.Icon(c, icon.Must("paperclip")).Size(16, 16).TextColor(t.Accent)
				ui.Column(c).Children(func() {
					ui.Text(c, attachment.Name)
					kit.Muted(c, bytes(attachment.Size)).FontSize(t.Rem(0.75))
				})
			})
		}
	})
	return root
}
