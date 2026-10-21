package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PictureOptions tune a Picture.
type PictureOptions struct {
	// Fit is Contain, Cover or FillBox.
	Fit ui.Fit
}

// Picture is a rounded fading image with Fit.
func Picture(c *ui.Context, image *ui.Bitmap, o PictureOptions) ui.Element {
	if image == nil || o.Fit < ui.Contain || o.Fit > ui.FillBox {
		panic("display: Picture needs an image and valid fit")
	}
	img := ui.Image(c, image).Radius(c.Theme().Radius).Clip().Fit(o.Fit).ContentTransition(ui.ContentTransition{})
	return kit.Arrive(img, ui.Pose{Fade: 1}, 0)
}
