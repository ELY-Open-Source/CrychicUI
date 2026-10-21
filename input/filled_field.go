package input

import "github.com/egoist/mygo/ui"

// FilledFieldOptions tune a FilledField.
type FilledFieldOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// FilledField is a filled field whose color animates.
func FilledField(c *ui.Context, label string, value *string, o FilledFieldOptions) ui.Element {
	if value == nil {
		panic("input: FilledField needs a value")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		in := ui.TextInputBase(c, value).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius)
		bg := t.Surface
		if in.Focused() {
			bg = t.SurfaceHover
		}
		in.Background(in.AnimateColor("fill", bg, ui.Snappy))
		changed(root, in)
	})
}
