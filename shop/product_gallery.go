package shop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProductGalleryOptions word the gallery; Item takes the number.
type ProductGalleryOptions struct {
	Item, Previous, Next string
}

// ProductGallery shows a main image over thumbnails; cue Select.
func ProductGallery(c *ui.Context, images []*ui.Bitmap, selected *int, o ProductGalleryOptions) ui.Element {
	if len(images) == 0 || selected == nil || *selected < 0 || *selected >= len(images) {
		panic(fmt.Sprintf("shop: ProductGallery needs images and a selection, got %d of %d", deref(selected), len(images)))
	}
	for _, image := range images {
		if image == nil {
			panic("shop: ProductGallery images must not be nil")
		}
	}
	def(&o.Item, "Image %d")
	def(&o.Previous, "Previous image")
	def(&o.Next, "Next image")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	pick := func(i int, event string) {
		*selected = i
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("ProductGallery", event, "index", i)
	}
	var prev, next ui.Element
	root.Children(func() {
		ui.Grid(c).Columns(1).FillWidth().Radius(t.Radius).Clip().Background(t.Surface).Children(func() {
			// Keyed by image: the new one fades over the old.
			c.Key(*selected)
			ui.Image(c, images[*selected]).ColumnStart(1).RowStart(1).FillWidth().AspectRatio(4.0 / 3).Fit(ui.Cover).
				Transition(ui.ElementTransition{Duration: 220 * time.Millisecond, Enter: &ui.Motion{}, Exit: &ui.Motion{}})
			ui.Row(c).ColumnStart(1).RowStart(1).FillWidth().AlignSelf(ui.Center).Justify(ui.SpaceBetween).Padding(0, t.Space(2)).Children(func() {
				arrow := func(name, glyph string, off bool) ui.Element {
					b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(9), t.Space(9)).Center().Radius(999)).Disabled(off)
					b.Background(b.AnimateColor("bg", t.Background.Alpha(0.75+0.25*unit(b.Hovered())), ui.Snappy)).
						Shadow(0, 1, 3, 0, kit.Shade(0.08))
					return b.Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16) })
				}
				prev = arrow(o.Previous, "chevron-left", *selected == 0)
				next = arrow(o.Next, "chevron-right", *selected == len(images)-1)
			})
			ui.Box(c).ColumnStart(1).RowStart(1).AlignSelf(ui.End).JustifySelf(ui.End).Margin(t.Space(2)).Children(func() {
				kit.Numeric(c, strconv.Itoa(*selected+1)+" / "+strconv.Itoa(len(images))).FontSize(t.Rem(0.7)).Bold().
					Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.Background.Alpha(0.85)).ContentTransition(ui.ContentTransition{Numeric: true})
			})
		})
		strip := ui.ScrollHorizontal(c).FillWidth().Gap(t.Space(2)).Padding(0, 0, t.Space(1), 0).FocusGroup(ui.Horizontal)
		strip.Children(func() {
			seen := map[*ui.Bitmap]int{}
			for i, image := range images {
				// Keyed by image and repeat, so reorders keep identity.
				c.Key(thumbKey{image, seen[image]})
				seen[image]++
				on := i == *selected
				name := fmt.Sprintf(o.Item, i+1)
				thumb := kit.Press(ui.ButtonBase(c).Label(name).Size(t.Space(14), t.Space(14)).Shrink(0).Radius(t.Radius).Clip().Checked(on))
				kit.Stagger(c, strip, thumb, i)
				thumb.Border(t.Space(0.5), thumb.AnimateColor("border", t.Border.Alpha(0).Mix(t.Accent, unit(on)), ui.Snappy))
				// Dims the image, so the thumb keeps its arrival fade.
				dim := thumb.AnimateTo("dim", 0.6+0.4*unit(on || thumb.Hovered()), ui.Snappy)
				thumb.Children(func() { ui.Image(c, image).Fill().Fit(ui.Cover).Opacity(dim) })
				if thumb.Clicked() && !on {
					pick(i, "selected")
				}
			}
		})
	})
	switch {
	case *selected > 0 && prev.Clicked():
		pick(*selected-1, "previous")
	case *selected < len(images)-1 && next.Clicked():
		pick(*selected+1, "next")
	}
	return root
}

type thumbKey struct {
	image *ui.Bitmap
	n     int
}

func deref(p *int) int {
	if p == nil {
		return -1
	}
	return *p
}
