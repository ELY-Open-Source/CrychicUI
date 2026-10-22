package form

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Shipping is what a ShippingForm edits; Speed names one.
type Shipping struct {
	Address input.Address
	Speed   string
}

// Speed is a delivery option; Detail holds price and days.
type Speed struct{ Name, Detail string }

// ShippingFormOptions tune a ShippingForm.
type ShippingFormOptions struct {
	// Delivery titles the speed cards; Address labels the rows.
	Delivery string
	Address  input.AddressFieldOptions
}

// ShippingForm edits an address and picks a delivery speed card.
func ShippingForm(c *ui.Context, value *Shipping, speeds []Speed, o ShippingFormOptions) ui.Element {
	if value == nil || len(speeds) == 0 {
		panic("form: ShippingForm needs a value and speeds")
	}
	def(&o.Delivery, "Delivery")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4))
	return root.Children(func() {
		edited(root, input.AddressField(c, &value.Address, o.Address))
		ui.Text(c, o.Delivery).FontWeight(600)
		for _, s := range speeds {
			if choice.OptionCard(c, s.Name, s.Detail, s.Name, &value.Speed).Changed() {
				root.SetChanged()
				kit.Log("ShippingForm", "speed")
			}
		}
	})
}
