package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Category is one storefront category tile.
type Category struct {
	Name  string
	Image *ui.Bitmap
	Count int
}

// CategoriesOptions word the count, "%d" taking it.
type CategoriesOptions struct{ Item, Items string }

// Categories shows image tiles that zoom on hover; cue Tap.
func Categories(c *ui.Context, categories []Category, selected *int, o CategoriesOptions) ui.Element {
	if len(categories) == 0 || selected == nil || *selected < -1 || *selected >= len(categories) {
		panic("shop: Categories needs categories and a selection or -1")
	}
	def(&o.Item, "%d item")
	def(&o.Items, "%d items")
	t := c.Theme()
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, cat := range categories {
			if cat.Name == "" || cat.Image == nil || cat.Count < 0 || slices.IndexFunc(categories, func(o Category) bool { return o.Name == cat.Name }) != i {
				panic(fmt.Sprintf("shop: Categories tile %q needs a unique name, an image and a count", cat.Name))
			}
			c.Key(cat.Name)
			on := i == *selected
			tile := kit.Press(ui.ButtonBase(c).Label(cat.Name).Column().AlignItems(ui.Stretch).Radius(t.Radius).Clip().Checked(on))
			tile.Border(1, tile.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			tile.Background(tile.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(tile.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
			kit.Stagger(c, root, tile, i)
			tile.Children(func() {
				ui.Box(c).FillWidth().AspectRatio(1.4).Clip().Children(func() {
					img := ui.Image(c, cat.Image).FillWidth().FillHeight().Fit(ui.Cover)
					z := 1 + 0.06*img.AnimateTo("zoom", unit(tile.Hovered()), ui.Smooth)
					img.Scale(z, z)
					badge := ui.Box(c).Absolute().Top(t.Space(1.5)).Right(t.Space(1.5)).Size(t.Space(6), t.Space(6)).Center().
						Radius(999).Background(t.Accent).TextColor(t.AccentText)
					s := max(0, badge.AnimateTo("on", unit(on), ui.Bouncy))
					badge.Scale(s, s).Children(func() { ui.Icon(c, icon.Must("check")).Size(14, 14) })
				})
				ui.Column(c).FillWidth().Gap(t.Space(0.25)).Padding(t.Space(2), t.Space(2.5)).Children(func() {
					ui.Text(c, cat.Name).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
					items := o.Items
					if cat.Count == 1 {
						items = o.Item
					}
					kit.Muted(c, fmt.Sprintf(items, cat.Count)).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
				})
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
