package chat

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Mention is a highlighted "@name" in text.
func Mention(c *ui.Context, name string) ui.Element {
	if name == "" {
		panic("chat: Mention needs a name")
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label("@"+name).Padding(0, t.Space(1)).Radius(t.Radius).TextColor(t.Accent))
	root.Background(root.AnimateColor("bg", t.Accent.Alpha(0.1+0.1*unit(root.Hovered())), ui.Snappy))
	root.Children(func() { ui.Text(c, "@"+name).Bold().SingleLine() })
	if kit.Clicked(root, sound.Tap) {
		kit.Log("Mention", "opened")
	}
	return root
}
