package input

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

// MaskedFieldOptions tune a MaskedField.
type MaskedFieldOptions struct {
	// Label and Placeholder are caller copy.
	Label, Placeholder string
}

// MaskedField fills a mask such as ###-##-####.
func MaskedField(c *ui.Context, label string, value *string, mask string, o MaskedFieldOptions) ui.Element {
	if value == nil || !strings.Contains(mask, "#") {
		panic("input: MaskedField needs a value and digit slots")
	}
	if o.Label == "" {
		o.Label = label
	}
	if o.Placeholder == "" {
		o.Placeholder = mask
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		in := ui.TextInput(c, value).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		if changed(root, in) {
			*value = applyDigitMask(*value, mask)
		}
	})
}

// applyDigitMask fills # slots with typed digits.
func applyDigitMask(text, mask string) string {
	digits := asciiDigits(text)
	var out strings.Builder
	at := 0
	for _, r := range mask {
		if r == '#' {
			if at >= len(digits) {
				return out.String()
			}
			out.WriteByte(digits[at])
			at++
		} else if at > 0 {
			out.WriteRune(r)
		}
	}
	return out.String()
}
