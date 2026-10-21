package input

import (
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// CardNumberFieldOptions tune a CardNumberField.
type CardNumberFieldOptions struct {
	// Label is caller copy.
	Label string
}

// CardNumberField groups card digits and swaps its brand icon.
func CardNumberField(c *ui.Context, value *string, o CardNumberFieldOptions) ui.Element {
	if value == nil {
		panic("input: CardNumberField needs a value")
	}
	if o.Label == "" {
		o.Label = "Card number"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, o.Label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			in := ui.TextInput(c, value).Grow(1).Placeholder(o.Label).Padding(t.Space(2), t.Space(3))
			if changed(root, in) {
				*value = groupCardDigits(*value)
			}
			ui.Icon(c, icon.Must(cardBrandIcon(*value))).Size(16, 16).TextColor(t.TextMuted).
				ContentTransition(ui.ContentTransition{})
		})
	})
}

// groupCardDigits formats card digits in fours.
func groupCardDigits(s string) string {
	d := asciiDigits(s)
	if len(d) > 19 {
		d = d[:19]
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

// cardBrandIcon picks an icon from a card prefix.
func cardBrandIcon(s string) string {
	digits := asciiDigits(s)
	if strings.HasPrefix(digits, "4") {
		return "credit-card-check"
	}
	if strings.HasPrefix(digits, "34") || strings.HasPrefix(digits, "37") {
		return "landmark"
	}
	if strings.HasPrefix(digits, "5") {
		return "wallet-cards"
	}
	if len(digits) >= 4 {
		if n, err := strconv.Atoi(digits[:4]); err == nil && n >= 2221 && n <= 2720 {
			return "wallet-cards"
		}
	}
	return "credit-card"
}
