package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SplitButton pairs a main action with a chevron menu.
func SplitButton(c *ui.Context, label string, targets []string, selected *int) (root, main ui.Element) {
	if len(targets) == 0 || selected == nil || *selected < -1 || *selected >= len(targets) {
		panic("button: SplitButton needs targets and valid selection")
	}
	t := c.Theme()
	root = ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Stretch)
	root.Children(func() {
		main = Solid(c, label).Radius(t.Radius, 0, 0, t.Radius)
		chevron := textBase(c, label+" menu").FillHeight().Radius(0, t.Radius, t.Radius, 0).
			Background(t.Accent).TextColor(t.AccentText)
		chevron.Children(func() { ui.Icon(c, icon.Must("chevron-down")).Size(16, 16) })
		chevron.Menu(func(m *ui.Menu) {
			chosen := false
			for i, target := range targets {
				if m.Item(target).Chosen() {
					*selected = i
					root.SetChanged()
					kit.Log("SplitButton", "selected", "index", i)
					chosen = true
				}
			}
			if !chosen {
				sound.Play(sound.Open)
			}
		})
	})
	return root, main
}
