package data

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RowActions is a "…" button opening actions.
func RowActions(c *ui.Context, label string, actions []Action, chosen *int, open *bool) ui.Element {
	actionsValid("RowActions", actions)
	if label == "" || chosen == nil || open == nil || *chosen < -1 || *chosen >= len(actions) {
		panic("data: RowActions needs a label, a choice and open")
	}
	t := c.Theme()
	root := ui.Box(c).Shrink(0)
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius)).Expanded(*open)
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || b.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("ellipsis")).Size(16, 16)
		})
		menu(c, root, b, "RowActions", label, actions, open, func(i int) {
			*chosen = i
			sound.Play(sound.Select)
			kit.Log("RowActions", "chose", "index", i)
		})
	})
	return root
}
