package input

import "github.com/egoist/mygo/ui"

// FloatingLabelOptions tune a FloatingLabel.
type FloatingLabelOptions struct {
	// Placeholder is used while the value is empty.
	Placeholder string
}

// FloatingLabel floats its label on focus or value.
func FloatingLabel(c *ui.Context, label string, value *string, o FloatingLabelOptions) ui.Element {
	if value == nil {
		panic("input: FloatingLabel needs a value")
	}
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(t.Space(14)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
	return root.Children(func() {
		active := root.FocusWithin() || *value != ""
		y := t.Space(4)
		if active {
			y = t.Space(2)
		}
		size := t.Rem(1)
		if active {
			size = t.Rem(0.85)
		}
		floating := ui.Text(c, label).Absolute().Left(t.Space(3)).Top(0).TextColor(t.TextMuted).PassThrough()
		floating.Translate(0, floating.AnimateTo("y", y, ui.Snappy)).FontSize(floating.AnimateTo("size", size, ui.Snappy))
		in := ui.TextInput(c, value).Fill().Label(label).Placeholder(o.Placeholder).Padding(t.Space(5), t.Space(3), t.Space(2), t.Space(3))
		changed(root, in)
	})
}
