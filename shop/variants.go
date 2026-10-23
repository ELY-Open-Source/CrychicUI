package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VariantStock is one variant and its stock.
type VariantStock struct {
	Name  string
	Stock int
}

// VariantsOptions are the variant labels.
type VariantsOptions struct{ Color, Size, Stock string }

// Variants selects color and size with stock; cue Select.
func Variants(c *ui.Context, colors, sizes []VariantStock, color, size *int, o VariantsOptions) ui.Element {
	if len(colors) == 0 || len(sizes) == 0 || color == nil || size == nil ||
		*color < -1 || *color >= len(colors) || *size < -1 || *size >= len(sizes) {
		panic("shop: Variants needs options and selections")
	}
	def(&o.Color, "Color")
	def(&o.Size, "Size")
	def(&o.Stock, "%d left")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, group := range []struct {
			name     string
			items    []VariantStock
			selected *int
		}{{o.Color, colors, color}, {o.Size, sizes, size}} {
			ui.Text(c, group.name).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
			ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
				for i, item := range group.items {
					if item.Name == "" || item.Stock < 0 {
						panic("shop: Variants has an invalid variant")
					}
					c.Key(i)
					name := fmt.Sprintf("%s (%s)", item.Name, fmt.Sprintf(o.Stock, item.Stock))
					b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1), t.Space(2)).Radius(999).
						Border(1, t.Border).Checked(*group.selected == i).Disabled(item.Stock == 0))
					b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*group.selected == i)), ui.Snappy))
					b.Children(func() { ui.Text(c, name).FontSize(t.Rem(0.8)).SingleLine() })
					if b.Clicked() && item.Stock > 0 && *group.selected != i {
						*group.selected = i
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("Variants", "selected", "group", group.name, "index", i)
					}
				}
			})
		}
	})
	return root
}
