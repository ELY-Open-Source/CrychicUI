package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CompareSpec is one specification row across products.
type CompareSpec struct {
	Name   string
	Values []string
}

// Compare shows products side by side; differing specs are bold.
func Compare(c *ui.Context, products []Product, specs []CompareSpec, currency string) ui.Element {
	if len(products) < 2 || len(products) > 3 || currency == "" {
		panic(fmt.Sprintf("shop: Compare needs two or three products and a currency, got %d", len(products)))
	}
	cheapest, ties := 0, 0
	for i, p := range products {
		validProduct("Compare", p)
		switch {
		case p.Price < products[cheapest].Price:
			cheapest, ties = i, 0
		case i != cheapest && p.Price == products[cheapest].Price:
			ties++
		}
	}
	for _, s := range specs {
		if s.Name == "" || len(s.Values) != len(products) {
			panic(fmt.Sprintf("shop: Compare spec %q needs one value per product", s.Name))
		}
	}
	t := c.Theme()
	lead := t.Space(22)
	root := kit.Card(c).FillWidth().Padding(0, 0, t.Space(2), 0).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Box(c).Width(lead).Shrink(0)
			for i, p := range products {
				c.Key(i)
				ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1)).Children(func() {
					if p.Image != nil {
						ui.Image(c, p.Image).FillWidth().AspectRatio(1.3).Fit(ui.Cover).Radius(t.Radius)
					}
					if p.Brand != "" {
						kit.Muted(c, p.Brand).FontSize(t.Rem(0.7)).FillWidth().SingleLine()
					}
					ui.Text(c, p.Name).Bold().FontSize(t.Rem(0.9)).FillWidth().Wrap()
					if p.Reviews > 0 {
						ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
							stars(c, p.Rating, t.Space(2.5))
							kit.Numeric(c, fmt.Sprintf("(%d)", p.Reviews)).FontSize(t.Rem(0.7)).TextColor(t.TextMuted).Shrink(0)
						})
					}
					price := money(c, p.Price, currency).Bold()
					if i == cheapest && ties == 0 {
						price.TextColor(t.Success)
					}
				})
			}
		})
		for i, s := range specs {
			c.Key(i)
			same := !slices.ContainsFunc(s.Values, func(v string) bool { return v != s.Values[0] })
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(3)).Background(t.Surface.Alpha(0.6 * unit(i%2 == 0)))
			kit.Stagger(c, root, row, i).Children(func() {
				kit.Muted(c, s.Name).FontSize(t.Rem(0.8)).Width(lead).Shrink(0).Wrap()
				for j, v := range s.Values {
					c.Key(j)
					cell := ui.Text(c, v).FontSize(t.Rem(0.85)).Grow(1).Basis(0).MinWidth(0).Wrap()
					if same {
						cell.TextColor(t.TextMuted)
					} else {
						cell.Bold()
					}
				}
			})
		}
	})
	return root
}
