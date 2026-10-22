package navigation

import "github.com/egoist/mygo/ui"

// SideNav is a column of destinations; the highlight slides.
func SideNav(c *ui.Context, label string, list []NavItem, active *int) ui.Element {
	items("SideNav", label, list, active)
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(56)).Gap(t.Space(0.5)).Padding(t.Space(2)).FocusGroup(ui.Vertical).Label(label)
	id := root.ID()
	root.Children(func() {
		if navList(c, "SideNav", id, list, active, 0) {
			root.SetChanged()
		}
	})
	return root
}
