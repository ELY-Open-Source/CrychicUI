package finance

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// OrderTicketOptions word the ticket.
type OrderTicketOptions struct{ Title, Quantity, Decrease, Increase, Buy, Sell, Price, Total string }

// OrderTicket sets a quantity, then buys or sells; cue Success.
func OrderTicket(c *ui.Context, quantity *int, price float64, currency string, o OrderTicketOptions) (root, buy, sell ui.Element) {
	if quantity == nil || *quantity < 1 || price <= 0 || !finite(price) || currency == "" || !finite(price*float64(*quantity)) {
		panic(fmt.Sprintf("finance: OrderTicket needs a quantity of one or more and a positive price, got %v", price))
	}
	def(&o.Title, "Place order")
	def(&o.Quantity, "Quantity")
	def(&o.Decrease, "Decrease")
	def(&o.Increase, "Increase")
	def(&o.Buy, "Buy")
	def(&o.Sell, "Sell")
	def(&o.Price, "Market price")
	def(&o.Total, "Estimated total")
	more := *quantity < math.MaxInt && finite(price*float64(*quantity+1))
	t := c.Theme()
	root = kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	line := func(label string, body func()) {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			kit.Muted(c, label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0)
			body()
		})
	}
	act := func(name, glyph string, ink ui.Color) ui.Element {
		b := kit.Press(ui.ButtonBase(c).Label(name).Grow(1).Basis(0).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).TextColor(t.AccentText))
		b.Background(b.AnimateColor("bg", ink.Mix(t.Text, 0.15*unit(b.Hovered())), ui.Snappy))
		return b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must(glyph)).Size(16, 16)
				ui.Text(c, name).Bold().Shrink(0)
			})
		})
	}
	root.Children(func() {
		kit.CardHead(c, "arrow-left-right", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			line(o.Quantity, func() {
				ui.Row(c).Shrink(0).AlignItems(ui.Center).Radius(999).Border(1, t.Border).Children(func() {
					for _, step := range []struct {
						name, glyph string
						delta       int
					}{{o.Decrease, "minus", -1}, {o.Increase, "plus", 1}} {
						can := step.delta < 0 && *quantity > 1 || step.delta > 0 && more
						b := kit.Press(ui.ButtonBase(c).Label(step.name).Tooltip(step.name).Size(t.Space(8), t.Space(8)).Center().Radius(999)).Disabled(!can)
						b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
						b.Children(func() { ui.Icon(c, icon.Must(step.glyph)).Size(14, 14) })
						if step.delta < 0 {
							kit.Numeric(c, strconv.Itoa(*quantity)).Bold().MinWidth(t.Space(8)).TextAlign(ui.Center).
								ContentTransition(ui.ContentTransition{Numeric: true})
						}
						if b.Clicked() && can {
							*quantity += step.delta
							root.SetChanged()
							sound.Play(sound.Tick)
							kit.Log("OrderTicket", "quantity", "value", *quantity)
						}
					}
				})
			})
			line(o.Price, func() { money(c, price, currency).FontSize(t.Rem(0.9)).Shrink(0) })
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), 0, 0, 0).
				BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
				ui.Text(c, o.Total).Bold().Grow(1).MinWidth(0)
				money(c, price*float64(*quantity), currency).Bold().FontSize(t.Rem(1.15)).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				buy = act(o.Buy, "arrow-up-right", t.Success)
				sell = act(o.Sell, "arrow-down-right", t.Danger)
			})
		})
	})
	if buy.Clicked() {
		sound.Play(sound.Success)
		kit.Log("OrderTicket", "buy", "quantity", *quantity)
	}
	if sell.Clicked() {
		sound.Play(sound.Success)
		kit.Log("OrderTicket", "sell", "quantity", *quantity)
	}
	return root, buy, sell
}
