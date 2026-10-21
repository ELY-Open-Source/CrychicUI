package input

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// KeyValueRow is one key and value input row.
type KeyValueRow struct{ Key, Value string }

// KeyValueFieldOptions tune a KeyValueField.
type KeyValueFieldOptions struct {
	// Add, Key and Value are caller copy.
	Add, Key, Value string
}

// KeyValueField edits rows of keys and values.
func KeyValueField(c *ui.Context, rows *[]KeyValueRow, o KeyValueFieldOptions) ui.Element {
	if rows == nil {
		panic("input: KeyValueField needs rows")
	}
	if o.Add == "" {
		o.Add = "Add row"
	}
	if o.Key == "" {
		o.Key = "Key"
	}
	if o.Value == "" {
		o.Value = "Value"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		removed := false
		for i := range *rows {
			if removed {
				break
			}
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).Transition(ui.ElementTransition{Position: true, Size: true})
			row.Children(func() {
				key := ui.TextInput(c, &(*rows)[i].Key).Grow(1).Basis(0).Placeholder(o.Key).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
				changed(root, key)
				value := ui.TextInput(c, &(*rows)[i].Value).Grow(1).Basis(0).Placeholder(o.Value).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
				changed(root, value)
				remove := ui.ButtonBase(c).Label("Remove row").Padding(t.Space(2)).TextColor(t.TextMuted)
				kit.Press(remove).Children(func() { ui.Icon(c, icon.Must("x")).Size(16, 16) })
				if kit.Clicked(remove, sound.Delete) {
					*rows = append((*rows)[:i], (*rows)[i+1:]...)
					root.SetChanged()
					kit.Log("KeyValueField", "removed", "count", len(*rows))
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
			*rows = append(*rows, KeyValueRow{})
			root.SetChanged()
			kit.Log("KeyValueField", "added", "count", len(*rows))
		}
	})
}
