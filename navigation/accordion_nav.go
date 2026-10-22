package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AccordionNav opens one group of links at once.
func AccordionNav(c *ui.Context, label string, gs []LinkGroup, picked *[2]int, opened *int) ui.Element {
	groups("AccordionNav", label, gs, picked)
	if opened == nil || *opened < -1 || *opened >= len(gs) {
		panic(fmt.Sprintf("navigation: AccordionNav %q needs an open group or -1", label))
	}
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(64)).Label(label)
	return root.Children(func() {
		ui.Accordion(c, func() {
			for g, grp := range gs {
				c.Key(g)
				open := *opened == g
				item := ui.AccordionItem(c, grp.Title, &open, func() {
					ui.Column(c).FillWidth().Gap(t.Space(0.5)).Transition(kit.Enter).Children(func() {
						for i, l := range grp.Links {
							c.Key(i)
							if kit.Clicked(row(c, l), sound.Tap) {
								*picked = [2]int{g, i}
								root.SetChanged()
								kit.Log("AccordionNav", "picked", "group", g, "index", i)
							}
						}
					})
				})
				if item.Changed() {
					*opened = -1
					if open {
						*opened = g
					}
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("AccordionNav", "toggled", "group", g, "open", open)
				}
			}
		})
	})
}
