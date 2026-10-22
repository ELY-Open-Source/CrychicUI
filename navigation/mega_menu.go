package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LinkGroup is a titled group of links.
type LinkGroup struct {
	Title string
	Links []string
}

// groups checks link groups and a picked group and link.
func groups(component, label string, gs []LinkGroup, picked *[2]int) {
	if label == "" || len(gs) == 0 || picked == nil {
		panic(fmt.Sprintf("navigation: %s needs a label, groups and a pick", component))
	}
	for g, grp := range gs {
		if grp.Title == "" || len(grp.Links) == 0 {
			panic(fmt.Sprintf("navigation: %s group %d needs a title and links", component, g))
		}
	}
	if *picked != [2]int{-1, -1} && (!valid(picked[0], len(gs)) || !valid(picked[1], len(gs[picked[0]].Links))) {
		panic(fmt.Sprintf("navigation: %s pick %v is no link", component, *picked))
	}
}

// MegaMenu opens a wide panel of grouped links.
func MegaMenu(c *ui.Context, label string, gs []LinkGroup, picked *[2]int, open *bool) ui.Element {
	groups("MegaMenu", label, gs, picked)
	if open == nil {
		panic(fmt.Sprintf("navigation: MegaMenu %q needs open", label))
	}
	t := c.Theme()
	b := ui.Box(c).AlignSelf(ui.Start)
	was := *open
	b.Children(func() {
		trigger := opener(c, label, *open, false)
		toggled(trigger, open)
		ui.Popover(c, trigger, open, func() {
			panel := ui.Row(c).Gap(t.Space(8)).Padding(t.Space(2))
			kit.Arrive(panel, ui.Pose{Scale: 0.96, Fade: 1}, 0).Children(func() {
				for g, grp := range gs {
					c.Key(g)
					ui.Column(c).Width(t.Space(40)).Gap(t.Space(0.5)).Children(func() {
						kit.Muted(c, grp.Title).FontSize(t.Rem(0.75)).Bold().Padding(t.Space(1), t.Space(3))
						for i, l := range grp.Links {
							c.Key(i)
							if kit.Clicked(row(c, l), sound.Tap) {
								*picked, *open = [2]int{g, i}, false
								b.SetChanged()
								kit.Log("MegaMenu", "picked", "group", g, "index", i)
							}
						}
					})
				}
			})
		})
	})
	reopened(b, was, *open, "MegaMenu")
	return b
}
