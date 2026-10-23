package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StickerOptions tune a Sticker.
type StickerOptions struct {
	// Size of the square side; zero gives the default.
	Size float32
}

// Sticker is a large image without a bubble; pops in.
func Sticker(c *ui.Context, label string, image *ui.Bitmap, o StickerOptions) ui.Element {
	if label == "" || image == nil {
		panic(fmt.Sprintf("chat: Sticker %q needs a label and an image", label))
	}
	t := c.Theme()
	size("Sticker", &o.Size, t.Space(32))
	root := ui.Box(c).AlignSelf(ui.Start).Label(label).Role(ui.RoleImage)
	kit.Arrive(root, ui.Pose{Scale: 0.5, Fade: 1}, 0)
	p := ui.Pose{}
	if root.Hovered() {
		p.Rotate = 4
	}
	return root.AnimatePose(p, ui.Bouncy).Children(func() {
		ui.Image(c, image).Size(o.Size, o.Size).Fit(ui.Contain)
	})
}
