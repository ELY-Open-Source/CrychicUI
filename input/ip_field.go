package input

import (
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// IPFieldOptions tune an IPField.
type IPFieldOptions struct {
	// Label, Placeholder and Error are caller copy.
	Label, Placeholder, Error string
}

// IPField edits four octets and validates their range.
func IPField(c *ui.Context, value *string, o IPFieldOptions) ui.Element {
	if value == nil {
		panic("input: IPField needs a value")
	}
	if o.Label == "" {
		o.Label = "IP address"
	}
	if o.Placeholder == "" {
		o.Placeholder = "192.168.0.1"
	}
	if o.Error == "" {
		o.Error = "Each octet must be 0 through 255."
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	pendingCaret := ui.Local(root, "pendingCaret", func() int { return -1 })
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		in := ui.TextInput(c, value).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		if *pendingCaret >= 0 {
			in.SetTextSelection(*pendingCaret, *pendingCaret)
			*pendingCaret = -1
		}
		if changed(root, in) {
			formatted := formatIPv4(*value)
			if formatted != *value {
				*value = formatted
				*pendingCaret = len(formatted)
			}
		}
		if *value != "" && !validIPv4(*value) {
			kit.Arrive(ui.Text(c, o.Error).TextColor(t.Danger).FontSize(t.Rem(0.85)), ui.Pose{Fade: 1}, 0)
		}
	})
}

// formatIPv4 groups digits and preserves typed octet boundaries.
func formatIPv4(s string) string {
	if strings.Contains(s, ".") {
		parts := strings.Split(s, ".")
		var octets []string
		for _, part := range parts {
			digits := asciiDigits(part)
			for len(digits) > 3 && len(octets) < 3 {
				octets = append(octets, digits[:3])
				digits = digits[3:]
			}
			octets = append(octets, digits)
		}
		if len(octets) > 4 {
			octets = octets[:4]
		}
		return strings.Join(octets, ".")
	}
	digits := asciiDigits(s)
	if len(digits) > 12 {
		digits = digits[:12]
	}
	var parts []string
	for len(digits) > 3 && len(parts) < 3 {
		parts = append(parts, digits[:3])
		digits = digits[3:]
	}
	if digits != "" {
		parts = append(parts, digits)
	}
	return strings.Join(parts, ".")
}

// validIPv4 checks all four octets' range.
func validIPv4(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 || n > 255 {
			return false
		}
	}
	return true
}
