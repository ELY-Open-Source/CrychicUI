package input

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Address is the value an AddressField edits.
type Address struct{ Street, City, Region, Postal string }

// AddressFieldOptions tune an AddressField.
type AddressFieldOptions struct {
	// Street, City, Region and Postal are row labels.
	Street, City, Region, Postal string
}

// AddressField edits street, city, region and postal code.
func AddressField(c *ui.Context, value *Address, o AddressFieldOptions) ui.Element {
	if value == nil {
		panic("input: AddressField needs a value")
	}
	if o.Street == "" {
		o.Street = "Street"
	}
	if o.City == "" {
		o.City = "City"
	}
	if o.Region == "" {
		o.Region = "Region"
	}
	if o.Postal == "" {
		o.Postal = "Postal code"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	rows := []struct {
		label string
		value *string
	}{{o.Street, &value.Street}, {o.City, &value.City}, {o.Region, &value.Region}, {o.Postal, &value.Postal}}
	return root.Children(func() {
		for i, row := range rows {
			line := ui.Row(c).Key(i).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			kit.Arrive(line, ui.Pose{Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				kit.Muted(c, row.label).Width(t.Rem(6))
				in := ui.TextInput(c, row.value).Grow(1).Placeholder(row.label).Padding(t.Space(2), t.Space(3)).
					Radius(t.Radius).Background(t.Background).Border(1, t.Border)
				changed(root, in)
			})
		}
	})
}
