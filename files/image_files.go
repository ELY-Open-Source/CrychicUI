package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ImageFile pairs a file with its thumbnail.
type ImageFile struct {
	Item  FileItem
	Image *ui.Bitmap
}

// ImageFiles shows selectable thumbnails with names; cue Select.
func ImageFiles(c *ui.Context, label string, images []ImageFile, selected *int) ui.Element {
	if label == "" || len(images) == 0 || selected == nil || *selected < -1 || *selected >= len(images) {
		panic(fmt.Sprintf("files: ImageFiles %q needs images and a selection or -1", label))
	}
	list := make([]FileItem, len(images))
	for i, img := range images {
		validFile("ImageFiles", img.Item)
		if img.Image == nil {
			panic(fmt.Sprintf("files: ImageFiles %q needs bitmap data for %q", label, img.Item.Name))
		}
		list[i] = img.Item
	}
	unique("ImageFiles", label, list)
	t := c.Theme()
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).FocusGroup(ui.Horizontal).Role(ui.RoleList).Label(label)
	root.Children(func() {
		for i, img := range images {
			c.Key(fileKey(img.Item))
			on := i == *selected
			tile := kit.Press(ui.ButtonBase(c).Label(img.Item.Name).Column().Gap(t.Space(1.5)).Padding(t.Space(1)).Radius(t.Radius).Checked(on))
			kit.Lift(tile)
			kit.Stagger(c, root, tile, i)
			tile.Background(tile.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(tile.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
			tile.Border(1, tile.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			tile.Children(func() {
				ui.Grid(c).Columns(1).FillWidth().Children(func() {
					ui.Image(c, img.Image).ColumnStart(1).RowStart(1).FillWidth().AspectRatio(4.0 / 3).Fit(ui.Cover).Radius(t.Radius)
					badge := ui.Box(c).ColumnStart(1).RowStart(1).AlignSelf(ui.Start).JustifySelf(ui.End).Margin(t.Space(1.5)).
						Size(t.Space(5), t.Space(5)).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText)
					k := badge.AnimateTo("on", unit(on), ui.Bouncy)
					badge.Scale(k, k).Opacity(min(1, k)).Children(func() { ui.Icon(c, icon.Must("check")).Size(12, 12) })
				})
				ui.Column(c).FillWidth().Padding(0, t.Space(1), t.Space(1), t.Space(1)).Children(func() {
					ui.Text(c, img.Item.Name).FontSize(t.Rem(0.8)).Bold().FillWidth().SingleLine()
					sizeText(c, img.Item.Size).FontSize(t.Rem(0.7))
				})
			})
			if tile.Clicked() && !on {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ImageFiles", "selected", "index", i)
			}
		}
	})
	return root
}
