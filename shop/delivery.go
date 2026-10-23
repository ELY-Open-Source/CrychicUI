package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DeliveryOption is one delivery speed and price.
type DeliveryOption struct {
	Name, Date string
	Price      float64
}

// Delivery shows speed options with dates and prices; cue Select.
func Delivery(c *ui.Context, label string, choices []DeliveryOption, currency string, selected *int) ui.Element {
	if label == "" || len(choices) == 0 || currency == "" || selected == nil || *selected < 0 || *selected >= len(choices) {
		panic("shop: Delivery needs choices and selection")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, choice := range choices {
			if choice.Name == "" || choice.Date == "" || choice.Price < 0 || !finite(choice.Price) {
				panic("shop: Delivery has an invalid choice")
			}
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(choice.Name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).Checked(i == *selected))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			b.Children(func() {
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, choice.Name).Bold()
					kit.Muted(c, choice.Date).FontSize(t.Rem(0.75))
				})
				money(c, choice.Price, currency)
			})
			if b.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Delivery", "selected", "index", i)
			}
		}
	})
	return root
}
