package content

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CaptionedImageOptions credit the image; zero Ratio is 3:2.
type CaptionedImageOptions struct {
	Credit string
	Ratio  float32
}

// CaptionedImage is an image with a caption below.
func CaptionedImage(c *ui.Context, img *ui.Bitmap, caption string, o CaptionedImageOptions) ui.Element {
	if o.Ratio == 0 {
		o.Ratio = 1.5
	}
	if img == nil || strings.TrimSpace(caption) == "" || math.IsNaN(float64(o.Ratio)) || o.Ratio < 0.25 || o.Ratio > 4 {
		panic(fmt.Sprintf("content: CaptionedImage %q needs an image, a caption and a ratio of 0.25 to 4", caption))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleImage).Label(caption)
	return root.Children(func() {
		frame := ui.Box(c).FillWidth().AspectRatio(o.Ratio).Clip().Radius(t.Radius*1.5).Border(1, t.Border)
		frame.Children(func() {
			kit.Arrive(ui.Image(c, img).Fill().Fit(ui.Cover), ui.Pose{Scale: 1.06, Fade: 1}, 0)
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
			ui.Box(c).Width(t.Space(0.5)).Height(t.Space(4)).Shrink(0).Radius(999).Background(t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Arrive(kit.Muted(c, caption).FontSize(t.Rem(0.8)).LineHeight(1.5).FillWidth().Wrap(), ui.Pose{Y: t.Space(1), Fade: 1}, 150)
				if o.Credit != "" {
					ui.Text(c, o.Credit).FontSize(t.Rem(0.68)).TextColor(t.TextMuted.Alpha(0.8)).FillWidth().MinWidth(0).SingleLine()
				}
			})
		})
	})
}
