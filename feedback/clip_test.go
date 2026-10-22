package feedback

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestUnlockToastShineStaysInside(t *testing.T) {
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(40).Children(func() { UnlockToast(c, "First streak", "Seven days.", &open) })
	}, 500, 220)
	var frames [][]uint8
	for range 12 {
		tt.Advance(50 * time.Millisecond)
		frames = append(frames, append([]uint8(nil), tt.Image().Pix...))
	}
	tt.Advance(2 * time.Second)
	img, w := tt.Image(), float32(40)+ui.LightTheme().Space(90)
	for i, f := range frames[5:] {
		for y := range 220 {
			for x := range 500 {
				if x > 36 && float32(x) < w+4 {
					continue
				}
				o := img.PixOffset(x, y)
				if f[o] != img.Pix[o] || f[o+1] != img.Pix[o+1] || f[o+2] != img.Pix[o+2] {
					t.Fatalf("frame %d: shine changes %d,%d beside the toast", i+5, x, y)
				}
			}
		}
	}
}
