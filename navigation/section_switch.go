package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SectionSwitch is segments over sections; the section slides.
func SectionSwitch(c *ui.Context, label string, names []string, active *int, section func(i int)) ui.Element {
	if section == nil {
		panic(fmt.Sprintf("navigation: SectionSwitch %q needs sections", label))
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Stretch)
	return root.Children(func() {
		if choice.Segments(c, label, names, active).Changed() {
			root.SetChanged()
			kit.Log("SectionSwitch", "switched", "index", *active)
		}
		kit.Slide(c, *active, func() { section(*active) })
	})
}
