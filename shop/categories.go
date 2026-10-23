package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Category is one storefront category tile.
type Category struct {
	Name  string
	Image *ui.Bitmap
}

// Categories shows image category tiles; cue Tap.
func Categories(c *ui.Context, categories []Category, selected *int) ui.Element {
	if len(categories) == 0 || selected == nil || *selected < -1 || *selected >= len(categories) {
		panic("shop: Categories needs categories and selection")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, category := range categories {
			if category.Name == "" || category.Image == nil {
				panic("shop: Categories needs names and images")
			}
			c.Key(i)
			tile := kit.Press(ui.ButtonBase(c).Label(category.Name).Column().Radius(t.Radius).Clip().Border(1, t.Border).
				Checked(i == *selected))
			tile.Background(tile.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			tile.Children(func() {
				ui.Image(c, category.Image).FillWidth().AspectRatio(1.4).Fit(ui.Cover)
				ui.Text(c, category.Name).Bold().Padding(t.Space(2))
			})
			if tile.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Categories", "selected", "index", i)
			}
		}
	})
	return root
}
