package input

import "github.com/egoist/mygo/ui"

// AdornedFieldOptions tune an AdornedField.
type AdornedFieldOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// AdornedField has prefix and suffix slots with focus tint.
func AdornedField(c *ui.Context, label string, value *string, prefix, suffix func(), o AdornedFieldOptions) ui.Element {
	if value == nil {
		panic("input: AdornedField needs a value")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2)).Padding(0, t.Space(2)).Radius(t.Radius).Border(1, t.Border)
		bg := t.Background
		if row.FocusWithin() {
			bg = t.Accent.Alpha(0.08)
		}
		row.Background(row.AnimateColor("focus", bg, ui.Snappy)).Children(func() {
			if prefix != nil {
				ui.Row(c).Children(prefix)
			}
			in := ui.TextInput(c, value).Grow(1).Placeholder(o.Placeholder).Padding(t.Space(2), 0)
			changed(root, in)
			if suffix != nil {
				ui.Row(c).Children(suffix)
			}
		})
	})
}
