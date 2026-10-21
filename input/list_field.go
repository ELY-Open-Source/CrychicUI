package input

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ListFieldOptions tune a ListField.
type ListFieldOptions struct {
	// Add and Item are caller copy.
	Add, Item string
}

// ListField edits rows of strings with add and remove.
func ListField(c *ui.Context, label string, items *[]string, o ListFieldOptions) ui.Element {
	if items == nil {
		panic("input: ListField needs items")
	}
	if o.Add == "" {
		o.Add = "Add item"
	}
	if o.Item == "" {
		o.Item = "Item"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		removed := false
		for i := range *items {
			if removed {
				break
			}
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).Transition(ui.ElementTransition{Position: true, Size: true})
			row.Children(func() {
				in := ui.TextInput(c, &(*items)[i]).Grow(1).Placeholder(o.Item).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
				changed(root, in)
				remove := ui.ButtonBase(c).Label("Remove " + o.Item).Padding(t.Space(2)).TextColor(t.TextMuted)
				kit.Press(remove).Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
				if kit.Clicked(remove, sound.Delete) {
					*items = append((*items)[:i], (*items)[i+1:]...)
					root.SetChanged()
					kit.Log("ListField", "removed", "count", len(*items))
					removed = true
				}
			})
		}
		add := ui.ButtonBase(c).Label(o.Add).Gap(t.Space(1)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
		kit.Press(add).Children(func() {
			ui.Icon(c, icon.Must("plus")).Size(16, 16)
			ui.Text(c, o.Add)
		})
		if kit.Clicked(add, sound.Pop) {
			*items = append(*items, "")
			root.SetChanged()
			kit.Log("ListField", "added", "count", len(*items))
		}
	})
}
