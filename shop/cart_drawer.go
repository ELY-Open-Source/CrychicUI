package shop

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CartDrawerOptions word the drawer.
type CartDrawerOptions struct {
	Title, Empty, Subtotal, Checkout, Note string
}

// CartDrawer shows cart items, the subtotal and checkout; cue Open.
func CartDrawer(c *ui.Context, items []CartItem, currency string, checkout *bool, o CartDrawerOptions) (root, send ui.Element) {
	if checkout == nil || currency == "" {
		panic("shop: CartDrawer needs a currency and checkout state")
	}
	// Quantities may sum past int, so count wide.
	total, count := 0.0, new(big.Int)
	for _, item := range items {
		validProduct("CartDrawer", item.Product)
		if item.Quantity <= 0 || item.Quantity > item.Product.Stock {
			panic(fmt.Sprintf("shop: CartDrawer %q needs a quantity within stock", item.Product.Name))
		}
		total += item.Product.Price * float64(item.Quantity)
		count.Add(count, big.NewInt(int64(item.Quantity)))
	}
	def(&o.Title, "Your cart")
	def(&o.Empty, "Your cart is empty")
	def(&o.Subtotal, "Subtotal")
	def(&o.Checkout, "Checkout")
	def(&o.Note, "Shipping and taxes are added at checkout")
	t := c.Theme()
	root = kit.Card(c).FillWidth().Role(ui.RoleDialog).Label(o.Title)
	kit.Arrive(root, ui.Pose{X: t.Space(4), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "shopping-bag", o.Title, func() {
			kit.Chip(c, count.String(), t.Accent)
		})
		ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).Children(func() {
			if len(items) == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(5)).Transition(kit.Enter).Children(func() {
					ui.Icon(c, icon.Must("shopping-bag")).Size(22, 22).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
				})
			}
			for i, item := range items {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1))
				kit.Stagger(c, root, row, i).Children(func() {
					if item.Product.Image != nil {
						ui.Image(c, item.Product.Image).Size(t.Space(12), t.Space(12)).Shrink(0).Fit(ui.Cover).Radius(t.Radius)
					}
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, item.Product.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().SingleLine()
						ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
							kit.Numeric(c, strconv.Itoa(item.Quantity)+" ×").FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
							money(c, item.Product.Price, currency).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
						})
					})
					money(c, item.Product.Price*float64(item.Quantity), currency).FontSize(t.Rem(0.85)).Bold()
				})
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Background(t.Surface).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Subtotal).Bold().Grow(1)
				money(c, total, currency).Bold().FontSize(t.Rem(1.05)).ContentTransition(ui.ContentTransition{Numeric: true})
			})
			kit.Muted(c, o.Note).FontSize(t.Rem(0.75))
			c.Key("checkout")
			send = kit.Press(ui.ButtonBase(c).Label(o.Checkout).FillWidth().Padding(t.Space(2), t.Space(3)).Radius(t.Radius)).Disabled(len(items) == 0)
			send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy)).TextColor(t.AccentText).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, o.Checkout).Bold().Shrink(0)
					arrow := ui.Icon(c, icon.Must("arrow-right")).Size(14, 14)
					arrow.Translate(arrow.AnimateTo("nudge", t.Space(1)*unit(send.Hovered()), ui.Snappy), 0)
				})
			})
		})
	})
	if send.Clicked() && len(items) > 0 {
		*checkout = true
		root.SetChanged()
		sound.Play(sound.Open)
		kit.Log("CartDrawer", "checkout", "items", len(items))
	}
	return root, send
}
