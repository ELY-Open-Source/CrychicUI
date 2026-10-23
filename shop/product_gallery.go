package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProductGalleryOptions are the thumbnail labels.
type ProductGalleryOptions struct{ Item string }

// ProductGallery shows a main image and thumbnails; cue Select.
func ProductGallery(c *ui.Context, images []*ui.Bitmap, selected *int, o ProductGalleryOptions) ui.Element {
	if len(images) == 0 || selected == nil || *selected < 0 || *selected >= len(images) {
		panic("shop: ProductGallery needs images and selection")
	}
	def(&o.Item, "Image %d")
	for _, image := range images {
		if image == nil {
			panic("shop: ProductGallery image must not be nil")
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		c.Key(*selected)
		ui.Image(c, images[*selected]).FillWidth().Height(t.Space(80)).Fit(ui.Cover).Radius(t.Radius).Transition(kit.Enter)
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			for i, image := range images {
				c.Key(i)
				name := fmt.Sprintf(o.Item, i+1)
				thumb := kit.Press(ui.ButtonBase(c).Label(name).Size(t.Space(14), t.Space(14)).Radius(t.Radius).Clip().
					Checked(i == *selected))
				thumb.Border(1, thumb.AnimateColor("border", t.Border.Mix(t.Accent, unit(i == *selected)), ui.Snappy))
				thumb.Children(func() { ui.Image(c, image).Fill().Fit(ui.Cover) })
				if thumb.Clicked() && *selected != i {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("ProductGallery", "selected", "index", i)
				}
			}
		})
	})
	return root
}
