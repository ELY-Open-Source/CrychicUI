package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/egoist/mygo/ui"
)

// DockBar is a bottom dock whose icons magnify by pointer.
func DockBar(c *ui.Context, label string, apps []NavItem, opened *int) ui.Element {
	if label == "" || len(apps) == 0 || opened == nil || *opened < -1 || *opened >= len(apps) {
		panic(fmt.Sprintf("navigation: DockBar %q needs apps and an opened app or -1", label))
	}
	t := c.Theme()
	items := make([]motion.MagnifyItem, len(apps))
	for i, a := range apps {
		items[i] = motion.MagnifyItem{Icon: a.Icon, Label: a.Label}
	}
	root := ui.Row(c).FillWidth().Justify(ui.Center).Padding(t.Space(2)).FocusGroup(ui.Horizontal).Label(label)
	return root.Children(func() {
		if motion.Magnify(c, items, opened).Changed() {
			root.SetChanged()
			kit.Log("DockBar", "opened", "index", *opened)
		}
	})
}
