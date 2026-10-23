package chat

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MessageMenuOptions name the four actions.
type MessageMenuOptions struct {
	Reply, Copy, Edit, Delete string
}

// MessageMenu shows reply, copy, edit and delete over hovered content.
func MessageMenu(c *ui.Context, label string, chosen *int, o MessageMenuOptions, content func()) ui.Element {
	if label == "" || chosen == nil || content == nil {
		panic(fmt.Sprintf("chat: MessageMenu %q needs a label, a choice and content", label))
	}
	def(&o.Reply, "Reply")
	def(&o.Copy, "Copy")
	def(&o.Edit, "Edit")
	def(&o.Delete, "Delete")
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Padding(t.Space(5), 0, 0, 0).MinWidth(t.Space(31)).Role(ui.RoleGroup).Label(label)
	shown := root.Hovered()
	names := [4]string{o.Reply, o.Copy, o.Edit, o.Delete}
	glyphs := [4]string{"reply", "copy", "pencil", "trash-2"}
	root.Children(func() {
		content()
		if !shown {
			return
		}
		menu := ui.Row(c).Key("menu").Absolute().Top(0).End(0).Gap(t.Space(0.5)).Padding(t.Space(0.5)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border).Shadow(0, 4, 12, 0, kit.Shade(0.12)).FocusGroup(ui.Horizontal).
			Role(ui.RoleToolbar).Transition(ui.ElementTransition{Duration: 140 * time.Millisecond, Enter: &ui.Motion{Y: 4}, Exit: &ui.Motion{Y: 4}})
		menu.Children(func() {
			for i, name := range names {
				c.Key(i)
				b := round(c, name, glyphs[i]).Size(t.Space(7), t.Space(7)).Radius(t.Radius)
				if i == 3 {
					b.TextColor(t.Danger)
				}
				if b.Clicked() {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("MessageMenu", "chosen", "index", i)
				}
			}
		})
	})
	return root
}
