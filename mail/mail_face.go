package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// mailFace is a message's avatar, sender, subject, preview and time.
func mailFace(c *ui.Context, m Mail) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		dot := ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999)
		dot.Background(dot.AnimateColor("unread", t.Accent.Alpha(unit(m.Unread)), ui.Snappy))
		avatar(c, m.Sender, 9)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				sender := ui.Text(c, m.Sender).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
				if m.Attached {
					glyph(c, "paperclip", 12).TextColor(t.TextMuted)
				}
				when := kit.Muted(c, m.Time).FontSize(t.Rem(0.7)).Shrink(0)
				if m.Unread {
					sender.Bold()
					when.Bold().TextColor(t.Accent)
				}
			})
			subject := ui.Text(c, m.Subject).FontSize(t.Rem(0.8)).FillWidth().MinWidth(0).SingleLine()
			if m.Unread {
				subject.Bold()
			}
			if m.Preview != "" {
				kit.Muted(c, m.Preview).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
			}
		})
	})
}
