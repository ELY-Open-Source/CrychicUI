package input

import "github.com/egoist/mygo/ui"

// NameFieldOptions tune a NameField.
type NameFieldOptions struct {
	// Given and Family are placeholder copy.
	Given, Family string
}

// NameField edits given and family names in one row.
func NameField(c *ui.Context, given, family *string, o NameFieldOptions) ui.Element {
	if given == nil || family == nil {
		panic("input: NameField needs both values")
	}
	if o.Given == "" {
		o.Given = "Given name"
	}
	if o.Family == "" {
		o.Family = "Family name"
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		first := ui.TextInput(c, given).Grow(1).Basis(0).Placeholder(o.Given).Padding(t.Space(2), t.Space(3)).
			Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		changed(root, first)
		last := ui.TextInput(c, family).Grow(1).Basis(0).Placeholder(o.Family).Padding(t.Space(2), t.Space(3)).
			Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		changed(root, last)
	})
}
