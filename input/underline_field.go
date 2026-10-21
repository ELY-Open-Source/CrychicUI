package input

import "github.com/egoist/mygo/ui"

// UnderlineFieldOptions tune an UnderlineField.
type UnderlineFieldOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// UnderlineField grows an underline from its center on focus.
func UnderlineField(c *ui.Context, label string, value *string, o UnderlineFieldOptions) ui.Element {
	if value == nil {
		panic("input: UnderlineField needs a value")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		box := ui.Box(c).FillWidth()
		box.Children(func() {
			in := ui.TextInputBase(c, value).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), 0)
			changed(root, in)
			ui.Box(c).Absolute().Bottom(0).Left(0).Right(0).Height(1).Background(t.Border)
			line := ui.Box(c).Absolute().Bottom(0).Left(0).Right(0).Height(t.Space(0.5)).Background(t.Accent).TransformOrigin(0.5, 0.5)
			target := float32(0)
			if in.Focused() {
				target = 1
			}
			line.Scale(line.AnimateTo("focus", target, ui.Snappy), 1)
		})
	})
}
