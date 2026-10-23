package files

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ImageFile pairs a file with its thumbnail.
type ImageFile struct {
	Item  FileItem
	Image *ui.Bitmap
}

// ImageFiles shows selectable thumbnails; cue Select.
func ImageFiles(c *ui.Context, label string, images []ImageFile, selected *int) ui.Element {
	if label == "" || len(images) == 0 || selected == nil || *selected < -1 || *selected >= len(images) {
		panic("files: ImageFiles needs images and a valid selection")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, image := range images {
			validFile("ImageFiles", image.Item)
			if image.Image == nil {
				panic("files: ImageFiles needs bitmap data")
			}
			c.Key(i)
			tile := kit.Press(ui.ButtonBase(c).Label(image.Item.Name).Column().Radius(t.Radius).Clip().Border(1, t.Border).
				Checked(i == *selected))
			tile.Background(tile.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			tile.Children(func() {
				ui.Image(c, image.Image).FillWidth().AspectRatio(1).Fit(ui.Cover)
				ui.Text(c, image.Item.Name).FontSize(t.Rem(0.8)).SingleLine().Padding(t.Space(1))
			})
			if tile.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ImageFiles", "selected", "index", i)
			}
		}
	})
	return root
}
