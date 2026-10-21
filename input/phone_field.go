package input

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PhoneFieldOptions tune a PhoneField.
type PhoneFieldOptions struct {
	// CountryCode and Placeholder are caller copy.
	CountryCode, Placeholder string
}

// PhoneField groups digits as they are typed.
func PhoneField(c *ui.Context, label string, value *string, o PhoneFieldOptions) ui.Element {
	if value == nil {
		panic("input: PhoneField needs a value")
	}
	if o.CountryCode == "" {
		o.CountryCode = "+1"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			kit.Muted(c, o.CountryCode).Padding(t.Space(2), 0, t.Space(2), t.Space(3))
			in := ui.TextInput(c, value).Grow(1).Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3), t.Space(2), t.Space(1))
			if changed(root, in) {
				*value = groupDigits(*value, []int{3, 3, 4})
			}
		})
	})
}

// groupDigits groups digits by the given widths.
func groupDigits(s string, widths []int) string {
	d := asciiDigits(s)
	var parts []string
	for _, w := range widths {
		if len(d) <= w {
			break
		}
		parts = append(parts, d[:w])
		d = d[w:]
	}
	if d != "" {
		parts = append(parts, d)
	}
	return strings.Join(parts, " ")
}
