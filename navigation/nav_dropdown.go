package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NavDropdown is a link that opens a panel of sublinks.
func NavDropdown(c *ui.Context, label string, links []string, picked *int, open *bool) ui.Element {
	if label == "" || len(links) == 0 || picked == nil || open == nil || *picked < -1 || *picked >= len(links) {
		panic(fmt.Sprintf("navigation: NavDropdown %q needs links, a pick or -1, and open", label))
	}
	t := c.Theme()
	b := ui.Box(c).AlignSelf(ui.Start)
	was := *open
	b.Children(func() {
		trigger := opener(c, label, *open, true)
		toggled(trigger, open)
		ui.Popover(c, trigger, open, func() {
			ui.Column(c).Width(t.Space(44)).Gap(t.Space(0.5)).Children(func() {
				for i, l := range links {
					c.Key(i)
					if kit.Clicked(row(c, l), sound.Tap) {
						*picked, *open = i, false
						b.SetChanged()
						kit.Log("NavDropdown", "picked", "index", i)
					}
				}
			})
		})
	})
	reopened(b, was, *open, "NavDropdown")
	return b
}
