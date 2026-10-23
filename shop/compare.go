package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CompareSpec is one specification row across products.
type CompareSpec struct {
	Name   string
	Values []string
}

// Compare shows products side by side with specifications.
func Compare(c *ui.Context, products []Product, specs []CompareSpec, currency string) ui.Element {
	if len(products) != 2 || currency == "" {
		panic("shop: Compare needs two products and currency")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for i, product := range products {
				validProduct("Compare", product)
				c.Key(i)
				ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Children(func() {
					ui.Text(c, product.Name).Bold()
					money(c, product.Price, currency)
				})
			}
		})
		for i, spec := range specs {
			if spec.Name == "" || len(spec.Values) != len(products) {
				panic("shop: Compare has an invalid specification")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				kit.Muted(c, spec.Name).FontSize(t.Rem(0.75))
				ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
					for j, value := range spec.Values {
						c.Key(j)
						ui.Text(c, value).Grow(1).Basis(0).Wrap()
					}
				})
			})
		}
	})
	return root
}
