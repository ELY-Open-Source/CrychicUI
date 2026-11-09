package media

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// stacked is red on top, blue below.
func stacked(h int) *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, 64, h))
	for y := range h {
		for x := range 64 {
			c := redInk
			if y >= h/2 {
				c = blueInk
			}
			img.SetRGBA(x, y, c)
		}
	}
	return ui.NewBitmap(img)
}

// lower checks the image fills its frame.
func lower(t *testing.T, view func(c *ui.Context)) {
	t.Helper()
	tt := ui.NewTester(func(c *ui.Context) { ui.Column(c).Width(480).Children(func() { view(c) }) }, 520, 720)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tt.Advance(2 * time.Second)
	r := find(t, tt, "View")
	if !at(tt, r.X+r.W/2, r.Y+r.H*0.6, blueInk) {
		t.Errorf("image does not fill its %v frame", r)
	}
}

func TestImageViewFillsRatio(t *testing.T) {
	img := stacked(48)
	lower(t, func(c *ui.Context) { ImageView(c, "View", img, ImageViewOptions{Ratio: 16.0 / 9, Fit: ui.FillBox}) })
}

func TestMediaGridFillsTile(t *testing.T) {
	items := []MediaItem{{Name: "Portrait", Image: stacked(96)}}
	picked := 0
	lower(t, func(c *ui.Context) { MediaGrid(c, "View", items, &picked, MediaGridOptions{Columns: 1}) })
}
