package form

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AddressFormOptions tune an AddressForm.
type AddressFormOptions struct {
	// Country labels the select; Address labels the rows.
	Country string
	Address input.AddressFieldOptions
}

// AddressForm edits a country and a street address.
func AddressForm(c *ui.Context, value *input.Address, country *string, countries []string, o AddressFormOptions) ui.Element {
	if value == nil || country == nil {
		panic("form: AddressForm needs an address and a country")
	}
	def(&o.Country, "Country")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4))
	return root.Children(func() {
		FieldGroup(c, o.Country, FieldGroupOptions{}, func() {
			if choice.Dropdown(c, o.Country, countries, country).Changed() {
				root.SetChanged()
				kit.Log("AddressForm", "country")
			}
		})
		edited(root, input.AddressField(c, value, o.Address))
	})
}
