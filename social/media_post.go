package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MediaPost shows an image grid of up to four.
func MediaPost(c *ui.Context, images []ui.ImageSource) ui.Element {
	if len(images) == 0 || len(images) > 4 {
		panic(fmt.Sprintf("social: MediaPost needs one to four images, got %d", len(images)))
	}
	for i, image := range images {
		if image == nil {
			panic(fmt.Sprintf("social: MediaPost image %d is nil", i))
		}
	}
	t := c.Theme()
	gap := t.Space(1)
	root := ui.Row(c).FillWidth().AspectRatio(16.0 / 10).Gap(gap).AlignItems(ui.Stretch).Radius(t.Radius).Clip().Role(ui.RoleGroup)
	cell := func(i int) {
		c.Key(i)
		box := ui.Box(c).Grow(1).Basis(0).MinWidth(0).MinHeight(0).Background(t.SurfaceHover).Clip()
		kit.Stagger(c, root, box, i).Children(func() { ui.Image(c, images[i]).Fill().Fit(ui.Cover) })
	}
	// group lays cells from..to in a row or a column.
	group := func(e ui.Element, from, to int) {
		e.Grow(1).Basis(0).MinWidth(0).MinHeight(0).Gap(gap).AlignItems(ui.Stretch).Children(func() {
			for i := from; i < to; i++ {
				cell(i)
			}
		})
	}
	root.Children(func() {
		switch len(images) {
		case 1:
			cell(0)
		case 2:
			cell(0)
			cell(1)
		case 3:
			cell(0)
			c.Key("side")
			group(ui.Column(c), 1, 3)
		case 4:
			c.Key("grid")
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(gap).AlignItems(ui.Stretch).Children(func() {
				c.Key("top")
				group(ui.Row(c), 0, 2)
				c.Key("bottom")
				group(ui.Row(c), 2, 4)
			})
		}
	})
	return root
}
