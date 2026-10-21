package input

import (
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ExpiryFieldOptions tune an ExpiryField.
type ExpiryFieldOptions struct {
	// Label and Error are caller copy.
	Label, Error string
}

// ExpiryField formats MM/YY and validates its month.
func ExpiryField(c *ui.Context, value *string, o ExpiryFieldOptions) ui.Element {
	if value == nil {
		panic("input: ExpiryField needs a value")
	}
	if o.Label == "" {
		o.Label = "Expiry"
	}
	if o.Error == "" {
		o.Error = "Use a month from 01 to 12."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		in := ui.TextInput(c, value).FillWidth().Placeholder("MM/YY").Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border)
		if changed(root, in) {
			*value = formatExpiry(*value)
		}
		if len(*value) == 5 && !validExpiry(*value) {
			kit.Arrive(ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
		}
	})
}

// formatExpiry inserts the MM/YY separator.
func formatExpiry(s string) string {
	d := asciiDigits(s)
	if len(d) > 4 {
		d = d[:4]
	}
	if len(d) > 2 {
		return d[:2] + "/" + d[2:]
	}
	return d
}

// validExpiry checks the month of a complete MM/YY value.
func validExpiry(s string) bool {
	month, err := strconv.Atoi(strings.TrimSuffix(s[:2], "/"))
	return err == nil && month >= 1 && month <= 12
}
