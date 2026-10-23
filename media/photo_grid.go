package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PhotoGridOptions tune a PhotoGrid.
type PhotoGridOptions struct {
	// Columns per row; zero gives 3.
	Columns int
}

// PhotoGrid lays out mixed sizes; tiles lift.
func PhotoGrid(c *ui.Context, label string, images []*ui.Bitmap, o PhotoGridOptions) ui.Element {
	pictures("PhotoGrid", label, images)
	if o.Columns < 0 {
		panic(fmt.Sprintf("media: PhotoGrid %q columns must not be negative", label))
	}
	if o.Columns == 0 {
		o.Columns = 3
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(o.Columns).FillWidth().Gap(t.Space(2)).Label(label)
	return root.Children(func() {
		for i, img := range images {
			c.Key(i)
			tile := ui.Box(c).AspectRatio(1).Radius(t.Radius).Clip().Background(t.Surface)
			if o.Columns > 1 && i%6 == 0 {
				tile.ColumnSpan(2).RowSpan(2)
			}
			pose := ui.Pose{}
			if tile.Hovered() {
				pose = ui.Pose{Y: -t.Space(1), Scale: 1.02}
			}
			kit.Arrive(tile.AnimatePose(pose, ui.Snappy), ui.Pose{Fade: 1, Scale: 0.96}, time.Duration(i)*40*time.Millisecond).
				Children(func() { ui.Image(c, img).Fill().Fit(ui.Cover) })
		}
	})
}
