package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MediaPost shows an image grid of up to four.
func MediaPost(c *ui.Context, images []ui.ImageSource) ui.Element {
	if len(images) == 0 || len(images) > 4 {
		panic("social: MediaPost needs one to four images")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, image := range images {
			if image == nil {
				panic("social: MediaPost image must not be nil")
			}
			c.Key(i)
			cell := ui.Box(c).AspectRatio(1).Radius(t.Radius).Clip()
			kit.Arrive(cell, ui.Pose{Scale: 0.96, Fade: 1}, 0)
			cell.Children(func() { ui.Image(c, image).Fill().Fit(ui.Cover) })
		}
	})
	return root
}
