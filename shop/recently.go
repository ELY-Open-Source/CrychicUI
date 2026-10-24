package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Recently is a scrolling strip of recently viewed products.
func Recently(c *ui.Context, label string, products []Product, currency string) ui.Element {
	if label == "" || len(products) == 0 || currency == "" {
		panic(fmt.Sprintf("shop: Recently %q needs products and a currency", label))
	}
	for i, p := range products {
		validProduct("Recently", p)
		if p.Image == nil || slices.IndexFunc(products, func(q Product) bool { return q.Name == p.Name }) != i {
			panic(fmt.Sprintf("shop: Recently %q product %q needs an image and a unique name", label, p.Name))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must("history")).Size(16, 16).TextColor(t.TextMuted)
			ui.Text(c, label).Bold().Shrink(0)
			kit.Chip(c, fmt.Sprint(len(products)), t.TextMuted)
		})
		strip := ui.ScrollHorizontal(c).FillWidth().Gap(t.Space(2.5)).Padding(0, 0, t.Space(1), 0)
		strip.Children(func() {
			for i, p := range products {
				c.Key(p.Name)
				card := ui.Column(c).Width(t.Space(28)).Shrink(0).Gap(t.Space(1))
				kit.Stagger(c, strip, card, i).Children(func() {
					ui.Image(c, p.Image).FillWidth().AspectRatio(1).Fit(ui.Cover).Radius(t.Radius)
					ui.Text(c, p.Name).FontSize(t.Rem(0.8)).FillWidth().SingleLine()
					if p.Reviews > 0 {
						stars(c, p.Rating, t.Space(2.5))
					}
					money(c, p.Price, currency).Bold().FontSize(t.Rem(0.8))
				})
			}
		})
	})
	return root
}
