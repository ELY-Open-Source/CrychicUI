package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Card is one saved payment card.
type Card struct {
	ID, Brand, Last4, Expiry string
}

// PaymentMethodsOptions word the card; Ending, Expires and Label are formats.
type PaymentMethodsOptions struct{ Title, Ending, Expires, Default, Label string }

// PaymentMethods picks the default card; cue Select.
func PaymentMethods(c *ui.Context, cards []Card, chosen *string, o PaymentMethodsOptions) ui.Element {
	if len(cards) == 0 || chosen == nil {
		panic("account: PaymentMethods needs cards and a choice")
	}
	ids := make([]string, len(cards))
	for i, card := range cards {
		if card.Brand == "" || len(card.Last4) != 4 || card.Expiry == "" {
			panic(fmt.Sprintf("account: PaymentMethods %q needs a brand, four last digits and expiry", card.ID))
		}
		ids[i] = card.ID
	}
	unique("PaymentMethods", ids)
	def(&o.Title, "Payment methods")
	def(&o.Ending, "•••• %s")
	def(&o.Expires, "Expires %s")
	def(&o.Default, "Default")
	def(&o.Label, "%s ending %s")
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Warning, t.Success, t.Danger}
	root := card(c).Role(ui.RoleRadioGroup)
	root.Children(func() {
		kit.CardHead(c, "credit-card", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		list.Children(func() {
			for i, card := range cards {
				on := *chosen == card.ID
				c.Key(card.ID)
				b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, card.Brand, card.Last4)).FillWidth().
					Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.06), unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						ink := inks[i%len(inks)]
						face := ui.Box(c).Size(t.Space(16), t.Space(10)).Shrink(0).Radius(t.Radius * 0.75).Padding(t.Space(1)).
							Background(ink.Alpha(0.16)).TextColor(ink)
						face.Children(func() {
							ui.Text(c, card.Brand).FontSize(t.Rem(0.6)).Bold().FillWidth().MinWidth(0).SingleLine()
						})
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							ui.Text(c, fmt.Sprintf(o.Ending, card.Last4)).Font("monospace").FontSize(t.Rem(0.85)).Bold().
								FillWidth().MinWidth(0).SingleLine()
							kit.Muted(c, fmt.Sprintf(o.Expires, card.Expiry)).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
						})
						if on {
							c.Key("default")
							kit.Arrive(kit.Chip(c, o.Default, t.Accent), ui.Pose{Scale: 0.6, Fade: 1}, 0)
						}
					})
				})
				if b.Clicked() && !on {
					*chosen = card.ID
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("PaymentMethods", "default", "index", i)
				}
			}
		})
	})
	return root
}
