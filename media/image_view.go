package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ImageViewOptions tune an ImageView.
type ImageViewOptions struct {
	// Ratio is width over height; zero gives 4:3.
	Ratio float32
	// Fit is Contain, Cover or FillBox.
	Fit ui.Fit
}

// ImageView fades an image in; nil shows a loading placeholder.
func ImageView(c *ui.Context, label string, image *ui.Bitmap, o ImageViewOptions) ui.Element {
	if label == "" || o.Fit < ui.Contain || o.Fit > ui.FillBox {
		panic(fmt.Sprintf("media: ImageView %q needs a label and a valid fit", label))
	}
	size("ImageView", label, &o.Ratio, 4.0/3)
	t := c.Theme()
	root := ui.Grid(c).Columns(1).FillWidth().AspectRatio(o.Ratio).Radius(t.Radius).Clip().Background(t.Surface).
		Role(ui.RoleImage).Label(label)
	return root.Children(func() {
		if image == nil {
			c.Key("wait")
			feedback.Skeleton(c).ColumnStart(1).RowStart(1).Fill().Radius(0).
				Transition(ui.ElementTransition{Duration: 250 * time.Millisecond, Exit: &ui.Motion{}})
			return
		}
		c.Key("image")
		img := ui.Image(c, image).ColumnStart(1).RowStart(1).Fill().Fit(o.Fit)
		kit.Arrive(img, ui.Pose{Fade: 1, Scale: 1.02}, 0)
	})
}
