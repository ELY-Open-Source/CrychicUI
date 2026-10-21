package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AnchorHeading adds a hover anchor that copies its fragment.
func AnchorHeading(c *ui.Context, text, anchor string) (root, copyAnchor ui.Element) {
	if anchor == "" {
		panic("typography: AnchorHeading needs an anchor")
	}
	t := c.Theme()
	root = ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2))
	root.Children(func() {
		Heading(c, text, 2)
		copyAnchor = ui.ButtonBase(c).Label("Copy anchor").TextColor(t.Accent)
		target := float32(0)
		if root.Hovered() || copyAnchor.FocusVisible() {
			target = 1
		}
		kit.Press(copyAnchor).Opacity(copyAnchor.AnimateTo("show", target, ui.Snappy)).Children(func() {
			ui.Text(c, "#")
		})
		if kit.Clicked(copyAnchor, sound.Pop) {
			c.WriteClipboard("#" + anchor)
		}
	})
	return root, copyAnchor
}
