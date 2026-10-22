package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DiffItem is an item and how it changed.
type DiffItem struct {
	Change Change
	Text   string
}

// DiffListOptions tune a DiffList.
type DiffListOptions struct {
	// Added, Removed and Edited name the marks.
	Added, Removed, Edited string
}

// DiffList marks added, removed and changed items.
func DiffList(c *ui.Context, label string, items []DiffItem, o DiffListOptions) ui.Element {
	if label == "" || len(items) == 0 {
		panic(fmt.Sprintf("data: DiffList %q needs items", label))
	}
	for _, it := range items {
		if it.Change > ChangeEdited || it.Text == "" {
			panic(fmt.Sprintf("data: DiffList %q items need text and a known change", label))
		}
	}
	def(&o.Added, "Added")
	def(&o.Removed, "Removed")
	def(&o.Edited, "Changed")
	t := c.Theme()
	ink := map[Change]ui.Color{ChangeSame: t.TextMuted, ChangeAdded: t.Success, ChangeRemoved: t.Danger, ChangeEdited: t.Warning}
	mark := map[Change]string{ChangeAdded: "+", ChangeRemoved: "−", ChangeEdited: "~"}
	name := map[Change]string{ChangeAdded: o.Added, ChangeRemoved: o.Removed, ChangeEdited: o.Edited}
	root := frame(c, label).Role(ui.RoleList)
	return root.Children(func() {
		for i, it := range items {
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3)).
				Background(ink[it.Change].Alpha(0.1 * unit(it.Change != ChangeSame)))
			kit.Arrive(row, ui.Pose{X: t.Space(-2), Fade: 1}, time.Duration(i)*30*time.Millisecond).Children(func() {
				ui.Text(c, mark[it.Change]).Bold().TextColor(ink[it.Change]).Width(t.Space(4)).Shrink(0).Label(name[it.Change])
				text := ui.Text(c, it.Text).Grow(1)
				if it.Change == ChangeRemoved {
					text.Strikethrough().TextColor(t.TextMuted)
				}
			})
		}
	})
}
