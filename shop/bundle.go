package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BundleOptions are the bundle labels.
type BundleOptions struct{ Total string }

// Bundle picks items for a bundle price; cue Toggle.
func Bundle(c *ui.Context, items []Product, selected *[]bool, currency string, o BundleOptions) ui.Element {
	if len(items) == 0 || selected == nil || len(*selected) != len(items) || currency == "" {
		panic("shop: Bundle needs items and selection")
	}
	def(&o.Total, "Bundle price")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	total := 0.0
	root.Children(func() {
		for i, product := range items {
			validProduct("Bundle", product)
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(product.Name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).Checked((*selected)[i]))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit((*selected)[i])), ui.Snappy))
			b.Children(func() {
				ui.Text(c, product.Name).Grow(1)
				money(c, product.Price, currency)
			})
			if (*selected)[i] {
				total += product.Price
			}
			if b.Clicked() {
				(*selected)[i] = !(*selected)[i]
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Bundle", "toggled", "index", i, "selected", (*selected)[i])
			}
		}
		ui.Row(c).FillWidth().Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1)
			money(c, total, currency).Bold()
		})
	})
	return root
}
