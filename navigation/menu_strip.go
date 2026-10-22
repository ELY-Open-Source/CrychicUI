package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MenuStrip is a bar of titles opening links.
func MenuStrip(c *ui.Context, label string, menus []LinkGroup, picked *[2]int, opened *int) ui.Element {
	groups("MenuStrip", label, menus, picked)
	if opened == nil || *opened < -1 || *opened >= len(menus) {
		panic(fmt.Sprintf("navigation: MenuStrip %q needs an open menu or -1", label))
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(0.5)).Padding(t.Space(0.5)).Radius(t.Radius).Background(t.Surface).
		FocusGroup(ui.Horizontal).Label(label)
	was := *opened
	root.Children(func() {
		for g, m := range menus {
			c.Key(g)
			open := *opened == g
			trigger := opener(c, m.Title, open, false)
			toggled(trigger, &open)
			ui.Popover(c, trigger, &open, func() {
				ui.Column(c).Width(t.Space(44)).Gap(t.Space(0.5)).Children(func() {
					for i, l := range m.Links {
						c.Key(i)
						if kit.Clicked(row(c, l), sound.Tap) {
							*picked, open = [2]int{g, i}, false
							root.SetChanged()
							kit.Log("MenuStrip", "picked", "menu", g, "index", i)
						}
					}
				})
			})
			switch {
			case open:
				*opened = g
			case *opened == g:
				*opened = -1
			}
		}
	})
	reopened(root, was, *opened, "MenuStrip")
	return root
}
