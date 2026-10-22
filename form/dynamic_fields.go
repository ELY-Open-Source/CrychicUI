package form

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DynamicField is one added field: its kind and value.
type DynamicField struct{ Kind, Value string }

// DynamicFieldsOptions tune a DynamicFields.
type DynamicFieldsOptions struct {
	// Kind labels the type select; Add and Remove label buttons.
	Kind, Add, Remove string
}

// DynamicFields adds and removes fields of chosen kinds.
func DynamicFields(c *ui.Context, fields *[]DynamicField, kinds []string, o DynamicFieldsOptions) ui.Element {
	if fields == nil || len(kinds) == 0 {
		panic("form: DynamicFields needs fields and kinds")
	}
	for _, f := range *fields {
		if !slices.Contains(kinds, f.Kind) {
			panic(fmt.Sprintf("form: DynamicFields kind %q is not offered", f.Kind))
		}
	}
	def(&o.Kind, "Field type")
	def(&o.Add, "Add field")
	def(&o.Remove, "Remove field")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3))
	ids := ui.Local(root, "rows", func() *rows { return &rows{} })
	kind := ui.Local(root, "kind", func() string { return kinds[0] })
	keys := (*ids).keys(len(*fields))
	return root.Children(func() {
		gone := -1
		for i, key := range keys {
			f := &(*fields)[i]
			ui.Row(c).Key(key).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Transition(rowMotion(c)).Children(func() {
				kit.Muted(c, f.Kind).Width(t.Rem(6))
				edited(root, ui.TextInput(c, &f.Value).Label(f.Kind).Grow(1))
				if cross(c, o.Remove).Clicked() {
					gone = i
				}
			})
		}
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			choice.Dropdown(c, o.Kind, kinds, kind)
			if plain(c, o.Add).Clicked() {
				*fields = append(*fields, DynamicField{Kind: *kind})
				(*ids).add()
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("DynamicFields", "add", "count", len(*fields))
			}
		})
		if gone >= 0 {
			*fields = append((*fields)[:gone], (*fields)[gone+1:]...)
			(*ids).drop(gone)
			root.SetChanged()
			sound.Play(sound.Delete)
			kit.Log("DynamicFields", "remove", "count", len(*fields))
		}
	})
}
