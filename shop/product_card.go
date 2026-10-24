package shop

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProductCardOptions word the card.
type ProductCardOptions struct {
	Add, Added, SoldOut string
	// Sale words the badge, "%d" taking the percent off.
	Sale string
}

// ProductCard shows image, price, rating and add; cue Pop.
func ProductCard(c *ui.Context, p Product, currency string, added *bool, o ProductCardOptions) (root, add ui.Element) {
	validProduct("ProductCard", p)
	if p.Image == nil || added == nil {
		panic(fmt.Sprintf("shop: ProductCard %q needs an image and bound state", p.Name))
	}
	def(&o.Add, "Add to cart")
	def(&o.Added, "Added")
	def(&o.SoldOut, "Sold out")
	def(&o.Sale, "-%d%%")
	t := c.Theme()
	root = kit.Lift(kit.Card(c).FillWidth().Clip().Role(ui.RoleGroup).Label(p.Name))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	zoom := root.AnimateTo("zoom", unit(root.Hovered()), ui.Smooth)
	root.Children(func() {
		ui.Grid(c).Columns(1).FillWidth().Clip().Children(func() {
			ui.Image(c, p.Image).ColumnStart(1).RowStart(1).FillWidth().AspectRatio(4.0/3).Fit(ui.Cover).Scale(1+0.06*zoom, 1+0.06*zoom)
			if off := percentOff(p.Price, p.OldPrice); p.OldPrice > p.Price && off > 0 {
				ui.Box(c).ColumnStart(1).RowStart(1).AlignSelf(ui.Start).JustifySelf(ui.Start).Margin(t.Space(2)).Children(func() {
					SaleBadge(c, off, SaleBadgeOptions{Format: o.Sale})
				})
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(3)).Children(func() {
			if p.Brand != "" {
				kit.Muted(c, p.Brand).FontSize(t.Rem(0.7)).Bold()
			}
			ui.Text(c, p.Name).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
			if p.Reviews > 0 {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					stars(c, p.Rating, t.Space(3))
					kit.Muted(c, "("+strconv.Itoa(p.Reviews)+")").FontSize(t.Rem(0.7))
				})
			}
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.End).Children(func() {
				money(c, p.Price, currency).Bold()
				if p.OldPrice > p.Price {
					money(c, p.OldPrice, currency).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Strikethrough()
				}
			})
			name, glyph := o.Add, "shopping-bag"
			switch {
			case p.Stock == 0:
				name, glyph = o.SoldOut, "ban"
			case *added:
				name, glyph = o.Added, "check"
			}
			// Keyed: optional rows above must not move its focus.
			c.Key("add")
			add = kit.Press(ui.ButtonBase(c).Label(name).FillWidth().Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)).Disabled(p.Stock == 0)
			ink := t.Accent.Mix(t.AccentHover, unit(add.Hovered())).Mix(t.Success, unit(*added))
			add.Background(add.AnimateColor("bg", ink, ui.Snappy)).TextColor(t.AccentText).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(14, 14).ContentTransition(ui.ContentTransition{})
					ui.Text(c, name).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
				})
			})
		})
	})
	if add.Clicked() && p.Stock > 0 && !*added {
		*added = true
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("ProductCard", "added")
	}
	return root, add
}

// percentOff is the whole percent saved, rounded down.
func percentOff(price, old float64) int {
	if old <= price {
		return 0
	}
	return min(100, int(math.Floor((old-price)/old*100+1e-9)))
}
