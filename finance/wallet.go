package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Coin is one wallet asset.
type Coin struct {
	Symbol string
	Amount float64
	Value  float64
}

// WalletOptions are the wallet labels.
type WalletOptions struct{ Send string }

// Wallet shows coin balances and send; cue Send.
func Wallet(c *ui.Context, coins []Coin, currency string, sent *int, o WalletOptions) ui.Element {
	if len(coins) == 0 || currency == "" || sent == nil || *sent < -1 || *sent >= len(coins) {
		panic("finance: Wallet needs coins and send state")
	}
	def(&o.Send, "Send")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, coin := range coins {
			if coin.Symbol == "" || coin.Amount < 0 || coin.Value < 0 || !finite(coin.Amount) || !finite(coin.Value) {
				panic("finance: Wallet has an invalid coin")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, coin.Symbol).Bold().Grow(1)
				kit.Numeric(c, fmt.Sprintf("%g %s", coin.Amount, coin.Symbol))
				money(c, coin.Value, currency)
				b := kit.Press(ui.ButtonBase(c).Label(o.Send).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
				b.Children(func() { ui.Text(c, o.Send).SingleLine() })
				if b.Clicked() {
					*sent = i
					root.SetChanged()
					sound.Play(sound.Send)
					kit.Log("Wallet", "sent", "index", i)
				}
			})
		}
	})
	return root
}
