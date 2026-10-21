package button

import (
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TextButton is an inline text action with animated color.
func TextButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	root := textBase(c, label).Padding(t.Space(1), t.Space(1)).Background(t.Background.Alpha(0))
	color := t.Accent
	if root.Pressed() {
		color = t.AccentPressed
	} else if root.Hovered() {
		color = t.AccentHover
	}
	root.TextColor(root.AnimateColor("color", color, ui.Snappy)).Children(func() { ui.Text(c, label).SingleLine() })
	if root.Clicked() {
		sound.Play(sound.Tap)
	}
	return root
}
