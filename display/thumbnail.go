package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Thumbnail is a small image with an icon fallback.
func Thumbnail(c *ui.Context, image *ui.Bitmap, iconName string) ui.Element {
	if image == nil && iconName == "" {
		panic("display: Thumbnail needs an image or icon")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Radius(t.Radius).Clip().
		Background(t.Surface).Border(1, t.Border).Center()
	root.Children(func() {
		if image != nil {
			c.Key("image")
			img := ui.Image(c, image).Fill().Fit(ui.Cover).ContentTransition(ui.ContentTransition{})
			kit.Arrive(img, ui.Pose{Fade: 1}, 0)
			return
		}
		c.Key("icon")
		kit.Arrive(ui.Icon(c, icon.Must(iconName)).Size(16, 16).TextColor(t.TextMuted), ui.Pose{Fade: 1}, 0)
	})
	return kit.Arrive(root, ui.Pose{Fade: 1}, 0)
}
