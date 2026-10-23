package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShopFiltersOptions are the filter labels and range.
type ShopFiltersOptions struct {
	Price, Brands, Rating string
	Stars                 string
	Min, Max, Step        float64
}

// ShopFilters edits price, brands, and rating; cue Toggle.
func ShopFilters(c *ui.Context, low, high *float64, brands []string, chosen *[]bool, rating *int, o ShopFiltersOptions) ui.Element {
	if low == nil || high == nil || len(brands) == 0 || chosen == nil || len(*chosen) != len(brands) || rating == nil || *rating < 0 || *rating > 5 {
		panic("shop: ShopFilters needs bound filters")
	}
	def(&o.Price, "Price")
	def(&o.Brands, "Brands")
	def(&o.Rating, "Minimum rating")
	def(&o.Stars, "%d stars")
	if o.Min == 0 && o.Max == 0 {
		o.Min, o.Max, o.Step = 0, 1000, 10
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Price).Bold()
		slider := choice.RangeSlider(c, o.Price, low, high, o.Min, o.Max, o.Step, choice.RangeSliderOptions{})
		if slider.Changed() {
			root.SetChanged()
		}
		ui.Text(c, o.Brands).Bold()
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			for i, brand := range brands {
				if brand == "" {
					panic("shop: ShopFilters brand must not be empty")
				}
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(brand).Padding(t.Space(1), t.Space(2)).Radius(999).
					Border(1, t.Border).Checked((*chosen)[i]))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit((*chosen)[i])), ui.Snappy))
				b.Children(func() { ui.Text(c, brand).FontSize(t.Rem(0.8)) })
				if b.Clicked() {
					(*chosen)[i] = !(*chosen)[i]
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("ShopFilters", "brand", "index", i, "chosen", (*chosen)[i])
				}
			}
		})
		ui.Text(c, o.Rating).Bold()
		ui.Row(c).Gap(t.Space(1)).Children(func() {
			for i := 1; i <= 5; i++ {
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Stars, i)).Padding(t.Space(1), t.Space(2)).Radius(999).
					Border(1, t.Border).Checked(*rating == i))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*rating == i)), ui.Snappy))
				b.Children(func() { ui.Text(c, fmt.Sprint(i)).SingleLine() })
				if b.Clicked() && *rating != i {
					*rating = i
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("ShopFilters", "rating", "value", i)
				}
			}
		})
	})
	return root
}
