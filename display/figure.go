package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Figure is an image with a caption fading in.
func Figure(c *ui.Context, image *ui.Bitmap, caption string) ui.Element {
	if image == nil {
		panic("display: Figure needs an image")
	}
	t := c.Theme()
	return ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(2)).Children(func() {
		Picture(c, image, PictureOptions{})
		kit.Arrive(kit.Muted(c, caption).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 60*time.Millisecond)
	})
}
