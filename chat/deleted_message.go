package chat

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DeletedMessageOptions tune a DeletedMessage.
type DeletedMessageOptions struct {
	Text string
}

// DeletedMessage is a "Message deleted" placeholder; fades in.
func DeletedMessage(c *ui.Context, o DeletedMessageOptions) ui.Element {
	def(&o.Text, "Message deleted")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
		Radius(t.Radius * 2).TextColor(t.TextMuted).Role(ui.RoleStatus).Label(o.Text)
	root.Draw(kit.Inside(t.Radius*2, func(p *ui.Painter, r ui.Rect) { p.StrokeDashed(r, t.Border, t.Radius*2, 1) }))
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		ui.Icon(c, icon.Must("ban")).Size(16, 16)
		ui.Text(c, o.Text).Italic().FontSize(t.Rem(0.9))
	})
}
