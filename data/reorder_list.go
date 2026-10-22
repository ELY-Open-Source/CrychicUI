package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReorderListOptions tune a ReorderList.
type ReorderListOptions struct {
	// Height of the list; zero gives the default.
	Height float32
}

// ReorderList rows drag, or move by Alt+arrows.
func ReorderList(c *ui.Context, label string, items *[]string, o ReorderListOptions) ui.Element {
	if label == "" || items == nil {
		panic(fmt.Sprintf("data: ReorderList %q needs items", label))
	}
	unique("ReorderList", label, *items)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Height(height(c, o.Height))
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	sel := ui.Local(root, "chosen", func() int { return -1 })
	if *sel >= len(*items) {
		*sel = -1
	}
	place := func(rows []int, to int) {
		*items = move(*items, rows, to)
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("ReorderList", "moved", "rows", len(rows), "to", to)
	}
	list.Selected, list.RoundRows, list.Reorder = sel, true, place
	list.Key = func(i int) any { return (*items)[i] }
	root.Children(func() {
		l := ui.List(c, list, len(*items), func(i int) {
			slides(c, ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)), i).Children(func() {
				ui.Icon(c, icon.Must("grip-vertical")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, (*items)[i]).SingleLine()
			})
		}).Label(label).Grow(1)
		switch at := *sel; {
		case at > 0 && l.Shortcut(ui.Alt, ui.KeyUp):
			place([]int{at}, at-1)
		case at >= 0 && at < len(*items)-1 && l.Shortcut(ui.Alt, ui.KeyDown):
			place([]int{at}, at+2)
		}
	})
	return root
}
