package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShopFiltersOptions are the filter copy and price range.
type ShopFiltersOptions struct {
	Title, Price, Brands, Rating string
	// Stars names a rating row; Up follows its stars.
	Stars, Up      string
	Currency       string
	Min, Max, Step float64
}

// ShopFilters edits price, brands and a minimum rating; cue Toggle.
func ShopFilters(c *ui.Context, low, high *float64, brands []string, chosen *[]bool, rating *int, o ShopFiltersOptions) ui.Element {
	if low == nil || high == nil || len(brands) == 0 || chosen == nil || len(*chosen) != len(brands) || rating == nil || *rating < 0 || *rating > 4 {
		panic(fmt.Sprintf("shop: ShopFilters needs bound filters, one flag per brand and a rating in 0..4, got %d brands", len(brands)))
	}
	def(&o.Title, "Filters")
	def(&o.Price, "Price")
	def(&o.Brands, "Brands")
	def(&o.Rating, "Rating")
	def(&o.Stars, "%d stars")
	def(&o.Up, "& up")
	def(&o.Currency, "USD")
	if o.Min == 0 && o.Max == 0 {
		o.Min, o.Max, o.Step = 0, 1000, 10
	}
	active := 0
	for _, on := range *chosen {
		active += int(unit(on))
	}
	active += int(unit(*rating > 0)) + int(unit(*low > o.Min || *high < o.Max))
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	section := func(title string, last bool, body func()) {
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(!last))).Children(func() {
			ui.Text(c, title).Bold().FontSize(t.Rem(0.85)).Shrink(0)
			body()
		})
	}
	root.Children(func() {
		kit.CardHead(c, "sliders-horizontal", o.Title, func() {
			c.Key("active")
			kit.Chip(c, fmt.Sprint(active), t.Accent).Opacity(unit(active > 0))
		})
		section(o.Price, false, func() {
			slider := choice.RangeSlider(c, o.Price, low, high, o.Min, o.Max, o.Step, choice.RangeSliderOptions{})
			if slider.Changed() {
				root.SetChanged()
			}
			kit.Muted(c, moneyText(c, *low, o.Currency)+" – "+moneyText(c, *high, o.Currency)).FontSize(t.Rem(0.8)).
				ContentTransition(ui.ContentTransition{Numeric: true})
		})
		section(o.Brands, false, func() {
			ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(1)).Children(func() {
				for i, brand := range brands {
					if brand == "" || slices.Index(brands, brand) != i {
						panic(fmt.Sprintf("shop: ShopFilters brand %q must be unique text", brand))
					}
					c.Key(brand)
					on := (*chosen)[i]
					b := kit.Press(ui.ButtonBase(c).Label(brand).Row().Gap(t.Space(2)).AlignItems(ui.Center).
						Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius).Checked(on))
					b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
					b.Children(func() {
						box := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Center().Radius(t.Space(1)).TextColor(t.AccentText)
						box.Background(box.AnimateColor("fill", t.Accent.Alpha(unit(on)), ui.Snappy))
						box.Border(1, box.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
						s := max(0, box.AnimateTo("check", unit(on), ui.Bouncy))
						box.Children(func() { ui.Icon(c, icon.Must("check")).Size(12, 12).Scale(s, s) })
						ui.Text(c, brand).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					})
					if b.Clicked() {
						(*chosen)[i] = !on
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("ShopFilters", "brand", "index", i, "chosen", !on)
					}
				}
			})
		})
		section(o.Rating, true, func() {
			for n := 4; n >= 1; n-- {
				c.Key(n)
				on := *rating == n
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Stars, n)).FillWidth().Row().Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.Start).
					Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
				b.Children(func() {
					ring := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Center().Radius(999)
					ring.Border(1, ring.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
					ring.Children(func() {
						dot := ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Background(t.Accent)
						s := max(0, dot.AnimateTo("dot", unit(on), ui.Bouncy))
						dot.Scale(s, s)
					})
					stars(c, float64(n), t.Space(3.5))
					kit.Muted(c, o.Up).FontSize(t.Rem(0.8)).Shrink(0)
				})
				if b.Clicked() {
					*rating = n
					if on {
						*rating = 0
					}
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("ShopFilters", "rating", "value", *rating)
				}
			}
		})
	})
	return root
}
