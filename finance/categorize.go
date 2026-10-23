package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Categorize selects a transaction category; cue Select.
func Categorize(c *ui.Context, transaction Transaction, categories []string, selected *int) ui.Element {
	if transaction.Merchant == "" || len(categories) == 0 || selected == nil || *selected < -1 || *selected >= len(categories) {
		panic("finance: Categorize has invalid data")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, transaction.Merchant).Bold()
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			for i, category := range categories {
				if category == "" {
					panic("finance: Categorize category must not be empty")
				}
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(category).Padding(t.Space(1), t.Space(2)).Radius(999).
					Border(1, t.Border).Checked(i == *selected))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
				b.Children(func() { ui.Text(c, category).FontSize(t.Rem(0.8)) })
				if b.Clicked() && *selected != i {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Categorize", "selected", "index", i)
				}
			}
		})
	})
	return root
}
