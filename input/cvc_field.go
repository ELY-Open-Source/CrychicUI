package input

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CVCFieldOptions tune a CVCField.
type CVCFieldOptions struct {
	// Label and Error are caller copy.
	Label, Error string
}

// CVCField is three or four masked digits.
func CVCField(c *ui.Context, value *string, o CVCFieldOptions) ui.Element {
	if value == nil {
		panic("input: CVCField needs a value")
	}
	if o.Label == "" {
		o.Label = "CVC"
	}
	if o.Error == "" {
		o.Error = "Use three or four digits."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		in := ui.TextInput(c, value).FillWidth().Placeholder(o.Label).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border)
		in.Password()
		if changed(root, in) {
			digits := asciiDigits(*value)
			if len(digits) > 4 {
				digits = digits[:4]
			}
			*value = digits
		}
		if len(*value) > 0 && len(*value) < 3 {
			kit.Arrive(ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
		}
	})
}
