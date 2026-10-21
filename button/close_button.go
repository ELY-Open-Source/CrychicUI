package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CloseButton turns its cross ninety degrees on hover.
func CloseButton(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("button: CloseButton needs a label")
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Center().
		Radius(t.Radius).Background(t.Surface).Border(1, t.Border))
	root.Children(func() {
		cross := ui.Text(c, "×")
		cross.Rotate(cross.AnimateTo("turn", 90*unit(root.Hovered()), ui.Snappy))
	})
	if root.Clicked() {
		sound.Play(sound.Close)
	}
	return root
}
