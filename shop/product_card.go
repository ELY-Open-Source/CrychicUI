package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProductCardOptions are the card labels.
type ProductCardOptions struct{ Add string }

// ProductCard shows a product with add; cue Pop.
func ProductCard(c *ui.Context, p Product, currency string, added *bool, o ProductCardOptions) (root, add ui.Element) {
	validProduct("ProductCard", p)
	if p.Image == nil || added == nil {
		panic("shop: ProductCard needs an image and bound state")
	}
	def(&o.Add, "Add to cart")
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(52)).Radius(t.Radius).Clip().Border(1, t.Border).Background(t.Surface)
	root.Children(func() {
		image := ui.Image(c, p.Image).FillWidth().AspectRatio(1).Fit(ui.Cover)
		image.AnimatePose(ui.Pose{Scale: 1 + 0.05*unit(root.Hovered())}, ui.Snappy)
		ui.Column(c).Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			ui.Text(c, p.Name).Bold().SingleLine()
			money(c, p.Price, currency)
			add = addButton(c, o.Add).Disabled(p.Stock == 0)
			add.Children(func() { ui.Text(c, o.Add).SingleLine() })
		})
	})
	if add.Clicked() && p.Stock > 0 {
		*added = true
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("ProductCard", "added")
	}
	return root, add
}
