package finance

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Coin is one wallet asset; Name is optional.
type Coin struct {
	Symbol, Name string
	Amount       float64
	Value        float64
}

// WalletOptions word the card; SendFor names each send button.
type WalletOptions struct{ Title, Send, SendFor, Sent string }

// Wallet shows coin balances under their total; cue Send.
func Wallet(c *ui.Context, coins []Coin, currency string, sent *int, o WalletOptions) ui.Element {
	if len(coins) == 0 || currency == "" || sent == nil || *sent < -1 || *sent >= len(coins) {
		panic("finance: Wallet needs coins, a currency and send state")
	}
	total := 0.0
	for i, coin := range coins {
		if coin.Symbol == "" || coin.Amount < 0 || coin.Value < 0 || !finite(coin.Amount) || !finite(coin.Value) ||
			slices.IndexFunc(coins, func(x Coin) bool { return x.Symbol == coin.Symbol }) != i {
			panic(fmt.Sprintf("finance: Wallet coin %q needs a unique symbol and amounts of zero or more", coin.Symbol))
		}
		total += coin.Value
	}
	if !finite(total) {
		panic("finance: Wallet total overflows")
	}
	def(&o.Title, "Wallet")
	def(&o.Send, "Send")
	def(&o.SendFor, "Send %s")
	def(&o.Sent, "Sent")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "wallet", o.Title, func() {
			money(c, total, currency).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3))
		list.Children(func() {
			for i, coin := range coins {
				c.Key(coin.Symbol)
				var send ui.Element
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(coins)-1)))
				kit.Stagger(c, list, row, i).Children(func() {
					symbolTile(c, coin.Symbol)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, coin.Symbol).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
						if coin.Name != "" {
							kit.Muted(c, coin.Name).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
						}
					})
					ui.Column(c).Shrink(0).AlignItems(ui.End).Gap(t.Space(0.5)).Children(func() {
						money(c, coin.Value, currency).Bold().FontSize(t.Rem(0.85)).Shrink(0)
						kit.Numeric(c, strconv.FormatFloat(coin.Amount, 'f', -1, 64)+" "+coin.Symbol).
							FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Shrink(0)
					})
					if *sent == i {
						c.Key("sent")
						kit.Arrive(kit.Chip(c, o.Sent, t.Success), ui.Pose{Scale: 0.8, Fade: 1}, 0)
					}
					c.Key("send")
					send = kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.SendFor, coin.Symbol)).Tooltip(o.Send).
						Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Border(1, t.Border))
					send.Background(send.AnimateColor("bg", t.Accent.Alpha(0.12*unit(send.Hovered())), ui.Snappy)).
						TextColor(send.AnimateColor("ink", t.Text.Mix(t.Accent, unit(send.Hovered())), ui.Snappy))
					send.Children(func() { ui.Icon(c, icon.Must("send")).Size(14, 14) })
				})
				if send.Clicked() {
					*sent = i
					root.SetChanged()
					sound.Play(sound.Send)
					kit.Log("Wallet", "sent", "index", i)
				}
			}
		})
	})
	return root
}
