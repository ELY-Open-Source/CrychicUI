package input

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// IBANFieldOptions tune an IBANField.
type IBANFieldOptions struct {
	// Label and Error are caller copy.
	Label, Error string
}

// IBANField groups by four and validates the checksum.
func IBANField(c *ui.Context, value *string, o IBANFieldOptions) ui.Element {
	if value == nil {
		panic("input: IBANField needs a value")
	}
	if o.Label == "" {
		o.Label = "IBAN"
	}
	if o.Error == "" {
		o.Error = "Check the IBAN checksum."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		in := ui.TextInput(c, value).FillWidth().Placeholder(o.Label).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border)
		if changed(root, in) {
			*value = formatIBAN(*value)
		}
		if ibanComplete(*value) && !ibanValid(*value) {
			kit.Arrive(ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
		}
	})
}

// formatIBAN uppercases and groups an IBAN by four.
func formatIBAN(s string) string {
	var clean strings.Builder
	for _, r := range strings.ToUpper(s) {
		if r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			clean.WriteRune(r)
		}
	}
	d := clean.String()
	if len(d) > 34 {
		d = d[:34]
	}
	var parts []string
	for len(d) > 4 {
		parts = append(parts, d[:4])
		d = d[4:]
	}
	if d != "" {
		parts = append(parts, d)
	}
	return strings.Join(parts, " ")
}

func ibanClean(s string) string {
	var out strings.Builder
	for _, r := range strings.ToUpper(s) {
		if r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			out.WriteRune(r)
		}
	}
	return out.String()
}

func ibanComplete(s string) bool { return len(ibanClean(s)) >= 8 }

// ibanValid checks an IBAN's mod-97 checksum.
func ibanValid(s string) bool {
	clean := ibanClean(s)
	if len(clean) < 4 || len(clean) > 34 {
		return false
	}
	for i, r := range clean[:2] {
		if r < 'A' || r > 'Z' || clean[i+2] < '0' || clean[i+2] > '9' {
			return false
		}
	}
	rearranged := clean[4:] + clean[:4]
	mod := 0
	for _, r := range rearranged {
		switch {
		case r >= '0' && r <= '9':
			mod = (mod*10 + int(r-'0')) % 97
		case r >= 'A' && r <= 'Z':
			mod = (mod*100 + int(r-'A'+10)) % 97
		default:
			return false
		}
	}
	return mod == 1
}
