package mail

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Attachment is one attached file.
type Attachment struct {
	Name string
	Size int64
}

// MailAttachmentsOptions word the summary line.
type MailAttachmentsOptions struct{ One, Many string }

// MailAttachments shows attachment cards with sizes.
func MailAttachments(c *ui.Context, attachments []Attachment, o MailAttachmentsOptions) ui.Element {
	if len(attachments) == 0 {
		panic("mail: MailAttachments needs attachments")
	}
	def(&o.One, "1 attachment · %s")
	def(&o.Many, "%d attachments · %s")
	total := int64(0)
	names := make([]string, len(attachments))
	for i, a := range attachments {
		names[i] = a.Name
		if a.Name == "" || a.Size < 0 || total > math.MaxInt64-a.Size {
			panic(fmt.Sprintf("mail: MailAttachments has an invalid file %q", a.Name))
		}
		total += a.Size
	}
	unique("MailAttachments", names)
	summary := fmt.Sprintf(o.One, bytes(total))
	if len(attachments) > 1 {
		summary = fmt.Sprintf(o.Many, len(attachments), bytes(total))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
			glyph(c, "paperclip", 14)
			ui.Text(c, summary).FontSize(t.Rem(0.78)).Shrink(0)
		})
		cards := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2))
		cards.Children(func() {
			for i, a := range attachments {
				c.Key(a.Name)
				mark, ink := fileLook(t, a.Name)
				box := ui.Row(c).Width(t.Space(54)).Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Border(1, t.Border).Background(t.Background)
				kit.Stagger(c, cards, box, i).Children(func() {
					tile(c, mark, ink)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, a.Name).FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, bytes(a.Size)).FontSize(t.Rem(0.7))
					})
				})
			}
		})
	})
	return root
}

// fileLook picks an icon and ink from a file's extension.
func fileLook(t *ui.Theme, name string) (string, ui.Color) {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".pdf", ".doc", ".docx", ".txt", ".md":
		return "file-text", t.Danger
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return "file-image", t.Accent
	case ".zip", ".tar", ".gz", ".7z":
		return "file-archive", t.Warning
	case ".csv", ".xlsx", ".xls", ".numbers":
		return "file-spreadsheet", t.Success
	}
	return "file", t.TextMuted
}
