package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReturnsOptions are the return labels.
type ReturnsOptions struct{ Reason string }

// Returns selects items and return reasons; cue Select.
func Returns(c *ui.Context, items []Product, reasons []string, selected *int, reason *string, o ReturnsOptions) ui.Element {
	if len(items) == 0 || len(reasons) == 0 || selected == nil || reason == nil || *selected < -1 || *selected >= len(items) {
		panic("shop: Returns needs items, reasons, and selection")
	}
	def(&o.Reason, "Reason")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, product := range items {
			validProduct("Returns", product)
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(product.Name).Padding(t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Checked(i == *selected))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			b.Children(func() { ui.Text(c, product.Name).SingleLine() })
			if b.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Returns", "selected", "index", i)
			}
		}
		ui.Text(c, o.Reason).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			for _, item := range reasons {
				if item == "" {
					panic("shop: Returns reason must not be empty")
				}
				b := kit.Press(ui.ButtonBase(c).Label(item).Padding(t.Space(1), t.Space(2)).Radius(999).
					Border(1, t.Border).Checked(*reason == item))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*reason == item)), ui.Snappy))
				b.Children(func() { ui.Text(c, item).FontSize(t.Rem(0.8)) })
				if b.Clicked() && *reason != item {
					*reason = item
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Returns", "reason", "reason", item)
				}
			}
		})
	})
	return root
}
