package input

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TemplateFieldOptions tune a TemplateField.
type TemplateFieldOptions struct {
	// Placeholder and Insert are caller copy.
	Placeholder, Insert string
}

// TemplateField inserts variable tokens from a menu.
func TemplateField(c *ui.Context, label string, value *string, variables []string, o TemplateFieldOptions) ui.Element {
	if value == nil || len(variables) == 0 {
		panic("input: TemplateField needs a value and variables")
	}
	if o.Insert == "" {
		o.Insert = "Insert variable"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		in := ui.TextArea(c, value).FillWidth().Lines(3, 8).Placeholder(o.Placeholder).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		changed(root, in)
		menu := ui.MenuButton(c, o.Insert, func(m *ui.Menu) {
			for _, variable := range variables {
				if m.Item(variable).Chosen() {
					*value += "{" + variable + "}"
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("TemplateField", "inserted", "count", len(*value))
				}
			}
		})
		menu.AlignSelf(ui.Start)
	})
}
