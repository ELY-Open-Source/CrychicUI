package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShareButton opens share targets with an open cue.
func ShareButton(c *ui.Context, label string, targets []string, selected *int) ui.Element {
	if len(targets) == 0 || selected == nil || *selected < -1 || *selected >= len(targets) {
		panic("button: ShareButton needs targets and valid selection")
	}
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	root.Children(func() {
		ui.Icon(c, icon.Must("share-2")).Size(16, 16)
		ui.Text(c, label).SingleLine()
	})
	root.Menu(func(m *ui.Menu) {
		chosen := false
		for i, target := range targets {
			if m.Item(target).Chosen() {
				*selected = i
				root.SetChanged()
				kit.Log("ShareButton", "selected", "index", i)
				chosen = true
			}
		}
		if !chosen {
			sound.Play(sound.Open)
		}
	})
	return root
}
