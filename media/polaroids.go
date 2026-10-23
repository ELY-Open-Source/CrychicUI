package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PolaroidsOptions tune a Polaroids.
type PolaroidsOptions struct {
	// Size of a photo; zero gives the default.
	Size float32
}

// Polaroids stacks tilted prints; the top one lifts on hover.
func Polaroids(c *ui.Context, label string, images []*ui.Bitmap, captions []string, o PolaroidsOptions) ui.Element {
	pictures("Polaroids", label, images)
	if len(captions) != len(images) {
		panic(fmt.Sprintf("media: Polaroids %q needs one caption per image", label))
	}
	t := c.Theme()
	size("Polaroids", label, &o.Size, t.Space(40))
	tilts := [...]float32{-7, 5, -3, 8, -5, 3}
	root := ui.Grid(c).Columns(1).Padding(t.Space(6)).Label(label)
	return root.Children(func() {
		top := len(images) - 1
		for i, img := range images {
			c.Key(i)
			card := ui.Column(c).ColumnStart(1).RowStart(1).JustifySelf(ui.Center).AlignSelf(ui.Center).Gap(t.Space(2)).
				Padding(t.Space(2), t.Space(2), t.Space(5), t.Space(2)).Background(t.Background).Border(1, t.Border).
				Shadow(0, t.Space(1), t.Space(3), 0, kit.Shade(0.12))
			pose := ui.Pose{Rotate: tilts[i%len(tilts)]}
			if i == top && card.Hovered() {
				pose = ui.Pose{Y: -t.Space(3), Scale: 1.04}
			}
			kit.Arrive(card.AnimatePose(pose, ui.Bouncy), ui.Pose{Fade: 1, Y: t.Space(4)}, time.Duration(i)*60*time.Millisecond).
				Children(func() {
					ui.Image(c, img).Size(o.Size, o.Size).Fit(ui.Cover)
					kit.Muted(c, captions[i]).SingleLine().FontSize(t.Rem(0.8))
				})
		}
	})
}
