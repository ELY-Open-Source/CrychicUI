package input

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CurrencyFieldOptions tune a CurrencyField.
type CurrencyFieldOptions struct {
	// Currency is the visible symbol or code.
	Currency string
}

// CurrencyField prefixes currency and rounds to cents on blur.
func CurrencyField(c *ui.Context, label string, value *float64, o CurrencyFieldOptions) ui.Element {
	if value == nil || !finite(*value) {
		panic("input: CurrencyField needs a finite value")
	}
	if o.Currency == "" {
		o.Currency = "$"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	wasFocused := ui.Local(root, "focused", func() bool { return false })
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		row := ui.Row(c).FillWidth().AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		row.Children(func() {
			kit.Muted(c, o.Currency).Padding(t.Space(2), 0, t.Space(2), t.Space(3))
			in := ui.NumberInput(c, value, 0, 1000000000, 0.01).Grow(1).Padding(t.Space(2), t.Space(3), t.Space(2), t.Space(1))
			changed(root, in)
			if *wasFocused && !in.Focused() {
				*value = math.Round(*value*100) / 100
			}
			*wasFocused = in.Focused()
		})
	})
}
