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
	root := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(3)).
		Radius(t.Radius).Background(t.Text).TextColor(t.Background).Shadow(0, 4, 12, 0, kit.Shade(0.18)).Role(ui.RoleStatus)
	shown := root.AnimateTo("visible", unit(*visible), ui.Snappy)
	root.Opacity(shown).Translate(0, t.Space(3)*(1-shown))
	root.Children(func() {
		glyph(c, "circle-check", 16).TextColor(t.Success.Mix(t.Background, 0.25))
		ui.Text(c, o.Sent).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
		undo := kit.Press(ui.ButtonBase(c).Label(o.Undo).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius).Shrink(0).Disabled(!*visible))
		undo.Background(undo.AnimateColor("bg", t.Background.Alpha(0.12*unit(undo.Hovered())), ui.Snappy))
		undo.TextColor(t.Accent.Mix(t.Background, 0.35)).Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "undo-2", 14)
				ui.Text(c, o.Undo).FontSize(t.Rem(0.82)).Bold().Shrink(0)
			})
		})
		if undo.Clicked() && *visible {
			*visible = false
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("UndoSend", "undone")
		}
	})
	return root
}
