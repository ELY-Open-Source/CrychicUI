package shop

import (
	"fmt"
	"math/big"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MiniCartOptions word the popover; Each takes quantity and price.
type MiniCartOptions struct{ Cart, Empty, Subtotal, Each, More string }

// MiniCart is a counted bag opening the cart; cue Open.
func MiniCart(c *ui.Context, items []CartItem, currency string, open *bool, o MiniCartOptions) ui.Element {
	if open == nil || currency == "" {
		panic("shop: MiniCart needs a currency and open state")
	}
	total, count := 0.0, new(big.Int)
	for i, item := range items {
		validProduct("MiniCart", item.Product)
		if item.Quantity <= 0 || item.Quantity > item.Product.Stock ||
			slices.IndexFunc(items, func(x CartItem) bool { return x.Product.Name == item.Product.Name }) != i {
			panic(fmt.Sprintf("shop: MiniCart %q needs a unique line with a quantity within stock", item.Product.Name))
		}
		total += item.Product.Price * float64(item.Quantity)
		count.Add(count, big.NewInt(int64(item.Quantity)))
	}
	def(&o.Cart, "Cart")
	def(&o.Empty, "Your cart is empty")
	def(&o.Subtotal, "Subtotal")
	def(&o.Each, "%d × %s")
	def(&o.More, "99+")
	badge := count.String()
	if count.Cmp(big.NewInt(99)) > 0 {
		badge = o.More
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(o.Cart).Tooltip(o.Cart).Size(t.Space(11), t.Space(11)).Radius(999).Expanded(*open))
		b.Background(b.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(b.Hovered() || *open)), ui.Snappy))
		b.Children(func() {
			ui.Grid(c).Columns(1).Fill().Children(func() {
				ui.Icon(c, icon.Must("shopping-bag")).Size(18, 18).ColumnStart(1).RowStart(1).AlignSelf(ui.Center).JustifySelf(ui.Center)
				if count.Sign() > 0 {
					pill := kit.Numeric(c, badge).FontSize(t.Rem(0.65)).Bold().MinWidth(t.Space(4.5)).TextAlign(ui.Center).
						Padding(t.Space(0.25), t.Space(1)).Radius(999).Background(t.Accent).TextColor(t.AccentText).
						ColumnStart(1).RowStart(1).AlignSelf(ui.Start).JustifySelf(ui.End).ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Arrive(pill, ui.Pose{Scale: 0.5, Fade: 1}, 0)
				}
			})
		})
		if b.Clicked() {
			*open = !*open
			root.SetChanged()
			if *open {
				sound.Play(sound.Open)
			} else {
				sound.Play(sound.Close)
			}
			kit.Log("MiniCart", "toggled", "open", *open)
		}
		if ui.PopoverBase(c, b, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).Width(t.Space(76)).Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Y: -t.Space(1), Fade: 1}, 0)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Cart).Bold().Grow(1).MinWidth(0)
				kit.Chip(c, count.String(), t.Accent)
			})
			if len(items) == 0 {
				ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(4), 0).TextColor(t.TextMuted).Children(func() {
					ui.Icon(c, icon.Must("shopping-bag")).Size(24, 24)
					ui.Text(c, o.Empty).FontSize(t.Rem(0.85)).Shrink(0)
				})
				return
			}
			list := ui.Scroll(c).FillWidth().MaxHeight(t.Space(72)).Gap(t.Space(2))
			list.Children(func() {
				for i, item := range items {
					c.Key(item.Product.Name)
					row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center)
					kit.Stagger(c, list, row, i).Children(func() {
						thumb(c, item.Product.Image, t.Space(10))
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							ui.Text(c, item.Product.Name).FontSize(t.Rem(0.85)).FillWidth().SingleLine()
							kit.Muted(c, fmt.Sprintf(o.Each, item.Quantity, moneyText(c, item.Product.Price, currency))).FontSize(t.Rem(0.75))
						})
						money(c, item.Product.Price*float64(item.Quantity), currency).FontSize(t.Rem(0.85)).Bold().Shrink(0)
					})
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2), 0, 0, 0).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
				ui.Text(c, o.Subtotal).Bold().Grow(1).MinWidth(0)
				money(c, total, currency).Bold().Shrink(0)
			})
		}).Dismissed() {
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("MiniCart", "dismissed")
		}
	})
	return root
}
