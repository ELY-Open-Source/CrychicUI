package display

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestPicture(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Picture(c, testBitmap(), PictureOptions{}) }, 300, 200)
	settles(t, tt, "Picture")
	panics(t, "Picture image", func() { ui.Render(func(c *ui.Context) { Picture(c, nil, PictureOptions{}) }, 10, 10, 1) })
}

func TestThumbnail(t *testing.T) {
	var image *ui.Bitmap
	tt := ui.NewTester(func(c *ui.Context) { Thumbnail(c, image, "image") }, 300, 200)
	settles(t, tt, "Thumbnail fallback")
	image = testBitmap()
	tt.Frame()
	settles(t, tt, "Thumbnail image")
}

func TestFigure(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Figure(c, testBitmap(), "A caption") }, 300, 200)
	settles(t, tt, "Figure")
	find(t, tt, "A caption")
	panics(t, "Figure image", func() { ui.Render(func(c *ui.Context) { Figure(c, nil, "A caption") }, 10, 10, 1) })
}
