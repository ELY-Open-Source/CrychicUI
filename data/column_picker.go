package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ColumnPicker shows, hides and reorders columns.
func ColumnPicker(c *ui.Context, label string, order *[]string, hidden *ui.Selection[string]) ui.Element {
	if label == "" || order == nil || hidden == nil || len(*order) == 0 {
		panic(fmt.Sprintf("data: ColumnPicker %q needs columns and hidden ones", label))
	}
	unique("ColumnPicker", label, *order)
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(64))
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	sel := ui.Local(root, "chosen", func() int { return -1 })
	if *sel >= len(*order) {
		*sel = -1
	}
	place := func(rows []int, to int) {
		*order = move(*order, rows, to)
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("ColumnPicker", "moved", "to", to)
	}
	list.Selected, list.RoundRows, list.Reorder = sel, true, place
	list.Key = func(i int) any { return (*order)[i] }
	root.Children(func() {
		l := ui.List(c, list, len(*order), func(i int) {
			name := (*order)[i]
			slides(c, ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)), i).Children(func() {
				ui.Icon(c, icon.Must("grip-vertical")).Size(16, 16).TextColor(t.TextMuted)
				on := !hidden.Has(name)
				if checkbox(c, name, &on, false).Changed() {
					toggle(hidden, name)
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("ColumnPicker", "shown", "index", i, "on", on)
				}
				ui.Text(c, name).Grow(1).SingleLine().Opacity(0.5 + 0.5*unit(on))
			})
		}).Label(label)
		switch at := *sel; {
		case at > 0 && l.Shortcut(ui.Alt, ui.KeyUp):
			place([]int{at}, at-1)
		case at >= 0 && at < len(*order)-1 && l.Shortcut(ui.Alt, ui.KeyDown):
			place([]int{at}, at+2)
		}
	})
	return root
}
