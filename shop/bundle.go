package shop

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BundleOptions word the total; Off is saved from two items.
type BundleOptions struct {
	Title, Total, Save string
	Off                float64
}

// Bundle picks items for a bundle price; cue Toggle.
func Bundle(c *ui.Context, items []Product, selected *[]bool, currency string, o BundleOptions) ui.Element {
	if len(items) == 0 || selected == nil || len(*selected) != len(items) || currency == "" || !(o.Off >= 0 && o.Off < 1) {
		panic(fmt.Sprintf("shop: Bundle needs items, one flag each, a currency and Off in [0,1), got %v", o.Off))
	}
	for i, p := range items {
		validProduct("Bundle", p)
		if slices.IndexFunc(items, func(q Product) bool { return q.Name == p.Name }) != i {
			panic(fmt.Sprintf("shop: Bundle product %q must be unique", p.Name))
		}
	}
	def(&o.Title, "Frequently bought together")
	def(&o.Total, "Bundle price")
	def(&o.Save, "Save %s")
	sum, n := 0.0, 0
	for i, p := range items {
		if (*selected)[i] {
			sum += p.Price
			n++
		}
	}
	off := 0.0
	if n >= 2 {
		off = o.Off
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, o.Title).Bold()
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			for i, p := range items {
				c.Key(p.Name)
				if i > 0 {
					ui.Icon(c, icon.Must("plus")).Size(14, 14).TextColor(t.TextMuted)
				}
				img := thumb(c, p.Image, t.Space(14))
				img.Opacity(img.AnimateTo("dim", 0.35+0.65*unit((*selected)[i]), ui.Snappy))
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5))
		list.Children(func() {
			for i, p := range items {
				c.Key(p.Name)
				on := (*selected)[i]
				b := kit.Press(ui.ButtonBase(c).Label(p.Name).Row().FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).
					Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					checkbox(c, on)
					ui.Text(c, p.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					money(c, p.Price, currency).FontSize(t.Rem(0.85)).Shrink(0)
				})
				if b.Clicked() {
					(*selected)[i] = !on
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Bundle", "toggled", "index", i, "selected", !on)
				}
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), 0, 0, 0).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Text(c, o.Total).Bold().Grow(1).MinWidth(0)
			if off > 0 {
				kit.Chip(c, fmt.Sprintf(o.Save, moneyText(c, sum*off, currency)), t.Success)
				money(c, sum, currency).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).Strikethrough().Shrink(0)
			}
			money(c, sum*(1-off), currency).Bold().FontSize(t.Rem(1.05)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root
}
