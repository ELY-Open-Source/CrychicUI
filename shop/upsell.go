package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UpsellOptions word the add buttons, "%s" taking the product name.
type UpsellOptions struct{ Add, Added string }

// Upsell is a row of related products; cue Pop.
func Upsell(c *ui.Context, label string, products []Product, currency string, added *[]bool, o UpsellOptions) ui.Element {
	if label == "" || len(products) == 0 || currency == "" || added == nil || len(*added) != len(products) {
		panic(fmt.Sprintf("shop: Upsell %q needs products, a currency and one added flag each", label))
	}
	for i, p := range products {
		validProduct("Upsell", p)
		if p.Image == nil || slices.IndexFunc(products, func(q Product) bool { return q.Name == p.Name }) != i {
			panic(fmt.Sprintf("shop: Upsell %q product %q needs an image and a unique name", label, p.Name))
		}
	}
	def(&o.Add, "Add %s")
	def(&o.Added, "Added %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must("sparkles")).Size(16, 16).TextColor(t.Accent)
			ui.Text(c, label).Bold().Shrink(0)
		})
		row := ui.Row(c).FillWidth().Gap(t.Space(2.5))
		row.Children(func() {
			for i, p := range products {
				c.Key(p.Name)
				on := (*added)[i]
				card := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1)).Padding(0, 0, t.Space(2), 0).
					Radius(t.Radius).Clip().Border(1, t.Border)
				var b ui.Element
				kit.Stagger(c, row, card, i).Children(func() {
					ui.Grid(c).Columns(1).FillWidth().Children(func() {
						ui.Image(c, p.Image).ColumnStart(1).RowStart(1).FillWidth().AspectRatio(1).Fit(ui.Cover)
						name := fmt.Sprintf(o.Add, p.Name)
						if on {
							name = fmt.Sprintf(o.Added, p.Name)
						}
						b = kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).ColumnStart(1).RowStart(1).AlignSelf(ui.End).JustifySelf(ui.End).
							Margin(t.Space(1.5)).Size(t.Space(8), t.Space(8)).Center().Radius(999).Checked(on).Disabled(p.Stock == 0).
							Shadow(0, 1, 3, 0, kit.Shade(0.12)).TextColor(t.AccentText))
						b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())).Mix(t.Success, unit(on)), ui.Snappy))
						pop := max(0, b.AnimateTo("pop", unit(on), ui.Bouncy))
						s := 1 + 0.8*pop*(1-pop)
						b.Children(func() {
							glyph := "plus"
							if on {
								glyph = "check"
							}
							ui.Icon(c, icon.Must(glyph)).Size(16, 16).Scale(s, s).ContentTransition(ui.ContentTransition{})
						})
					})
					ui.Text(c, p.Name).FontSize(t.Rem(0.8)).Padding(0, t.Space(2)).FillWidth().SingleLine()
					money(c, p.Price, currency).Bold().FontSize(t.Rem(0.85)).Padding(0, t.Space(2))
				})
				if b.Clicked() && !on && p.Stock > 0 {
					(*added)[i] = true
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("Upsell", "added", "index", i)
				}
			}
		})
	})
	return root
}
