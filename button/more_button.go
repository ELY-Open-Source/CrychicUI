package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MoreButton shows three dots and opens a menu.
func MoreButton(c *ui.Context, label string, targets []string, selected *int) ui.Element {
	if label == "" || len(targets) == 0 || selected == nil || *selected < -1 || *selected >= len(targets) {
		panic("button: MoreButton needs a label, targets and selection")
	}
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Center().
		Radius(t.Radius).Background(t.Surface).Border(1, t.Border))
	root.Children(func() { ui.Icon(c, icon.Must("ellipsis")).Size(16, 16) })
	root.Menu(func(m *ui.Menu) {
		chosen := false
		for i, target := range targets {
			if m.Item(target).Chosen() {
				*selected = i
				root.SetChanged()
				kit.Log("MoreButton", "selected", "index", i)
				chosen = true
			}
		}
		if !chosen {
			sound.Play(sound.Open)
		}
	})
	return root
}
