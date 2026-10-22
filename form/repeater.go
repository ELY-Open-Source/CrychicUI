package form

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RepeaterOptions tune a Repeater.
type RepeaterOptions struct {
	// Add is button copy; Remove formats the group number in.
	Add, Remove string
}

// Repeater adds and removes groups of fields, one per item.
func Repeater[T comparable](c *ui.Context, items *[]T, group func(i int, item *T), o RepeaterOptions) ui.Element {
	if items == nil || group == nil {
		panic("form: Repeater needs items and a group")
	}
	def(&o.Add, "Add another")
	def(&o.Remove, "Remove group %d")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3))
	ids := ui.Local(root, "rows", func() *rows { return &rows{} })
	keys := (*ids).keys(len(*items))
	return root.Children(func() {
		gone := -1
		for i, key := range keys {
			box := ui.Row(c).Key(key).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border)
			box.Transition(rowMotion(c)).Children(func() {
				was := (*items)[i]
				ui.Column(c).Grow(1).Gap(t.Space(3)).Children(func() { group(i, &(*items)[i]) })
				if (*items)[i] != was {
					root.SetChanged()
				}
				if cross(c, fmt.Sprintf(o.Remove, i+1)).Clicked() {
					gone = i
				}
			})
		}
		if plain(c, o.Add).AlignSelf(ui.Start).Clicked() {
			var zero T
			*items = append(*items, zero)
			(*ids).add()
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("Repeater", "add", "count", len(*items))
		}
		if gone >= 0 {
			*items = append((*items)[:gone], (*items)[gone+1:]...)
			(*ids).drop(gone)
			root.SetChanged()
			sound.Play(sound.Delete)
			kit.Log("Repeater", "remove", "count", len(*items))
		}
	})
}
