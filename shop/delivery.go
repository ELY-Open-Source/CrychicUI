package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DeliveryOption is one delivery speed; a zero Price reads Free.
type DeliveryOption struct {
	Name, Date string
	Price      float64
}

// DeliveryOptions word free delivery.
type DeliveryOptions struct{ Free string }

// Delivery picks a speed from priced radio cards; cue Select.
func Delivery(c *ui.Context, label string, choices []DeliveryOption, currency string, selected *int, o DeliveryOptions) ui.Element {
	if label == "" || len(choices) == 0 || currency == "" || selected == nil || *selected < 0 || *selected >= len(choices) {
		panic(fmt.Sprintf("shop: Delivery %q needs choices, a currency and a selection", label))
	}
	for i, ch := range choices {
		if ch.Name == "" || ch.Date == "" || ch.Price < 0 || !finite(ch.Price) ||
			slices.IndexFunc(choices, func(d DeliveryOption) bool { return d.Name == ch.Name }) != i {
			panic(fmt.Sprintf("shop: Delivery %q choice %q needs a unique name, a date and a price", label, ch.Name))
		}
	}
	def(&o.Free, "Free")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).FocusGroup(ui.Vertical).Role(ui.RoleRadioGroup).Label(label)
	root.Children(func() {
		for i, ch := range choices {
			c.Key(ch.Name)
			on := i == *selected
			b := kit.Press(ui.ButtonBase(c).Label(ch.Name).Row().FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).
				Padding(t.Space(3)).Radius(t.Radius).Checked(on))
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			kit.Stagger(c, root, b, i).Children(func() {
				radio(c, on)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, ch.Name).Bold().FillWidth().SingleLine()
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
						ui.Icon(c, icon.Must("truck")).Size(12, 12)
						ui.Text(c, ch.Date).FontSize(t.Rem(0.75)).Shrink(0)
					})
				})
				if ch.Price == 0 {
					kit.Chip(c, o.Free, t.Success)
				} else {
					money(c, ch.Price, currency).Bold().Shrink(0)
				}
			})
			if b.Clicked() && !on {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Delivery", "selected", "index", i)
			}
		}
	})
	return root
}
