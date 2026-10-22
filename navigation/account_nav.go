package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AccountNav is an avatar opening an account menu.
func AccountNav(c *ui.Context, name string, image *ui.Bitmap, items []string, picked *int, open *bool) ui.Element {
	if name == "" || len(items) == 0 || picked == nil || open == nil || *picked < -1 || *picked >= len(items) {
		panic(fmt.Sprintf("navigation: AccountNav %q needs items, a pick or -1, and open", name))
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	was := *open
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Radius(999).Padding(t.Space(0.5))).Expanded(*open)
		trigger.Background(trigger.AnimateColor("ring", t.Accent.Alpha(0.3*unit(*open || trigger.Hovered())), ui.Snappy)).Children(func() {
			display.Avatar(c, name, image)
		})
		popup(c, trigger, open, func() {
			ui.Column(c).Width(t.Space(52)).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).Bold().SingleLine().Padding(t.Space(1.5), t.Space(3))
				ui.Box(c).FillWidth().Height(1).Background(t.Border)
				for i, item := range items {
					c.Key(i)
					if kit.Clicked(row(c, item), sound.Tap) {
						*picked, *open = i, false
						root.SetChanged()
						kit.Log("AccountNav", "picked", "index", i)
					}
				}
			})
		})
	})
	reopened(root, was, *open, "AccountNav")
	return root
}
