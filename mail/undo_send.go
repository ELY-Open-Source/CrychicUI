package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UndoSendOptions are the sent toast labels.
type UndoSendOptions struct{ Sent, Undo string }

// UndoSend undoes a sent message; cue Close.
func UndoSend(c *ui.Context, visible *bool, o UndoSendOptions) ui.Element {
	if visible == nil {
		panic("mail: UndoSend needs state")
	}
	def(&o.Sent, "Message sent")
	def(&o.Undo, "Undo")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
		Background(t.Text).TextColor(t.Background)
	root.Opacity(root.AnimateTo("visible", float32(unit(*visible)), ui.Snappy))
	root.Children(func() {
		ui.Text(c, o.Sent).Grow(1)
		undo := button(c, o.Undo).Disabled(!*visible)
		undo.Children(func() { ui.Text(c, o.Undo) })
		if undo.Clicked() && *visible {
			*visible = false
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("UndoSend", "undone")
		}
	})
	return root
}
