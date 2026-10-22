package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VirtualListOptions tune a VirtualList.
type VirtualListOptions struct {
	// Height of the list; zero gives the default.
	Height float32
}

// VirtualList builds rows in view; chosen may be nil.
func VirtualList(c *ui.Context, label string, n int, chosen *int, row func(i int), o VirtualListOptions) ui.Element {
	if label == "" || n < 0 || row == nil || (chosen != nil && (*chosen < -1 || *chosen >= n)) {
		panic(fmt.Sprintf("data: VirtualList %q needs a row count, rows and a choice in range", label))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Height(height(c, o.Height))
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	list.Selected, list.RoundRows = chosen, true
	root.Children(func() {
		l := ui.List(c, list, n, func(i int) {
			r := ui.Row(c).FillWidth().Padding(t.Space(2), t.Space(3))
			kit.Arrive(r, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() { row(i) })
		}).Label(label).Grow(1)
		if l.Changed() {
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("VirtualList", "chose", "index", *chosen)
		}
	})
	return root
}
