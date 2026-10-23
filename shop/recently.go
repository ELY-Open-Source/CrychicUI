package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Recently shows a recently viewed product strip.
func Recently(c *ui.Context, label string, products []Product, currency string) ui.Element {
	if label == "" || len(products) == 0 || currency == "" {
		panic("shop: Recently needs products and currency")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, product := range products {
			validProduct("Recently", product)
			if product.Image == nil {
				panic("shop: Recently needs product images")
			}
			c.Key(i)
			card := ui.Column(c).Grow(1).Basis(0).Radius(t.Radius).Clip().Border(1, t.Border)
			kit.Arrive(card, ui.Pose{X: t.Space(3), Fade: 1}, 0)
			card.Children(func() {
				ui.Image(c, product.Image).FillWidth().AspectRatio(1.4).Fit(ui.Cover)
				ui.Text(c, product.Name).FontSize(t.Rem(0.8)).SingleLine().Padding(t.Space(1))
				money(c, product.Price, currency).Padding(0, t.Space(1), t.Space(1), t.Space(1))
			})
		}
	})
	return root
}
