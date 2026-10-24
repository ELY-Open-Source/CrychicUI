package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VariantStock is one variant, its stock and optional swatch.
type VariantStock struct {
	Name   string
	Stock  int
	Swatch ui.Color
}

// VariantsOptions word the groups; Stock takes the count left.
type VariantsOptions struct{ Color, Size, Stock string }

// Variants picks a swatch and size, sold-out struck; cue Select.
func Variants(c *ui.Context, colors, sizes []VariantStock, color, size *int, o VariantsOptions) ui.Element {
	if len(colors) == 0 || len(sizes) == 0 || color == nil || size == nil ||
		*color < -1 || *color >= len(colors) || *size < -1 || *size >= len(sizes) {
		panic("shop: Variants needs options and selections or -1")
	}
	for _, list := range [][]VariantStock{colors, sizes} {
		seen := map[string]bool{}
		for _, v := range list {
			if v.Name == "" || seen[v.Name] || v.Stock < 0 {
				panic(fmt.Sprintf("shop: Variants needs unique named variants with stock, got %q", v.Name))
			}
			seen[v.Name] = true
		}
	}
	def(&o.Color, "Color")
	def(&o.Size, "Size")
	def(&o.Stock, "%d left")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		for g, group := range []struct {
			name     string
			items    []VariantStock
			selected *int
		}{{o.Color, colors, color}, {o.Size, sizes, size}} {
			c.Key(g)
			ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, group.name).FontSize(t.Rem(0.8)).Bold()
					if s := *group.selected; s >= 0 {
						kit.Muted(c, group.items[s].Name).FontSize(t.Rem(0.8)).ContentTransition(ui.ContentTransition{})
						if n := group.items[s].Stock; n < 5 {
							kit.Chip(c, fmt.Sprintf(o.Stock, n), t.Warning)
						}
					}
				})
				ui.Row(c).Wrap().Gap(t.Space(2)).FocusGroup(ui.Horizontal).Children(func() {
					for i, v := range group.items {
						c.Key(v.Name)
						on := *group.selected == i
						out := v.Stock == 0
						name := fmt.Sprintf("%s (%s)", v.Name, fmt.Sprintf(o.Stock, v.Stock))
						swatch := v.Swatch.A > 0
						b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(v.Name).Radius(999).Checked(on)).Disabled(out)
						if swatch {
							b.Size(t.Space(9), t.Space(9)).Padding(t.Space(0.75))
						} else {
							b.MinWidth(t.Space(11)).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)
						}
						b.Border(t.Space(0.5), b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)).Mix(t.TextMuted, 0.4*unit(b.Hovered() && !on)), ui.Snappy))
						b.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
							if !out {
								return
							}
							var slash ui.Path
							slash.MoveTo(r.X+r.W*0.2, r.Y+r.H*0.8).LineTo(r.X+r.W*0.8, r.Y+r.H*0.2)
							p.StrokePath(&slash, 1, t.TextMuted)
						}))
						b.Children(func() {
							if swatch {
								ui.Box(c).Fill().Radius(999).Background(v.Swatch)
								return
							}
							ui.Text(c, v.Name).FontSize(t.Rem(0.85)).Bold().Shrink(0).TextAlign(ui.Center)
						})
						if b.Clicked() && !out && !on {
							*group.selected = i
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("Variants", "selected", "group", g, "index", i)
						}
					}
				})
			})
		}
	})
	return root
}
