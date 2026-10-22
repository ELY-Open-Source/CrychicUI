package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WorkspaceSwitcher picks a workspace from a list with avatars.
func WorkspaceSwitcher(c *ui.Context, label string, spaces []string, picked *int, open *bool) ui.Element {
	if label == "" || len(spaces) == 0 || picked == nil || open == nil || !valid(*picked, len(spaces)) {
		panic(fmt.Sprintf("navigation: WorkspaceSwitcher %q needs workspaces, a pick and open", label))
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	was, shown := *open, *picked
	avatar := func(name string) { display.Avatar(c, name, nil).Size(t.Space(6), t.Space(6)) }
	root.Children(func() {
		trigger := kit.Press(ui.ButtonBase(c).Label(label).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)).Expanded(*open)
		trigger.Background(trigger.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || trigger.Hovered())), ui.Snappy)).Children(func() {
			avatar(spaces[*picked])
			ui.Text(c, spaces[*picked]).Bold().SingleLine().ContentTransition(ui.ContentTransition{})
			ui.Icon(c, icon.Must("chevrons-up-down")).Size(16, 16).TextColor(t.TextMuted)
		})
		popup(c, trigger, open, func() {
			ui.Column(c).Width(t.Space(56)).Gap(t.Space(0.5)).Children(func() {
				for i, name := range spaces {
					c.Key(i)
					if kit.Clicked(checked(c, name, i == shown, func() { avatar(name) }), sound.Select) {
						*picked, *open = i, false
						root.SetChanged()
						kit.Log("WorkspaceSwitcher", "picked", "index", i)
					}
				}
			})
		})
	})
	reopened(root, was, *open, "WorkspaceSwitcher")
	return root
}
