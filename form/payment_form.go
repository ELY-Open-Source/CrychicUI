package form

import (
	"regexp"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var expiryPattern = regexp.MustCompile(`^(0[1-9]|1[0-2])/\d\d$`)

// Payment is the value a PaymentForm edits.
type Payment struct{ Number, Expiry, CVC, Name string }

// PaymentFormOptions tune a PaymentForm.
type PaymentFormOptions struct {
	// Name and Pay are labels; the rest tune card fields.
	Name, Pay string
	Card      input.CardNumberFieldOptions
	Expiry    input.ExpiryFieldOptions
	CVC       input.CVCFieldOptions
}

// PaymentForm takes card, expiry, CVC and name, then pays.
func PaymentForm(c *ui.Context, value *Payment, o PaymentFormOptions) (root, pay ui.Element) {
	if value == nil {
		panic("form: PaymentForm needs a value")
	}
	def(&o.Name, "Name on card")
	def(&o.Pay, "Pay")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(4))
	root.Children(func() {
		edited(root, input.CardNumberField(c, &value.Number, o.Card))
		FormRow(c, []func(){
			func() { edited(root, input.ExpiryField(c, &value.Expiry, o.Expiry)) },
			func() { edited(root, input.CVCField(c, &value.CVC, o.CVC)) },
		})
		FieldGroup(c, o.Name, FieldGroupOptions{}, func() {
			edited(root, ui.TextInput(c, &value.Name).FillWidth())
		})
		pay = primary(c, o.Pay).Disabled(!payable(value))
	})
	if kit.Clicked(pay, sound.Success) {
		kit.Log("PaymentForm", "pay")
	}
	return root, pay
}

// payable reports whether every card field is complete.
func payable(v *Payment) bool {
	n, cvc := digits(v.Number), digits(v.CVC)
	return len(n) >= 13 && len(n) <= 19 && expiryPattern.MatchString(v.Expiry) && len(cvc) >= 3 && len(cvc) <= 4 &&
		len(cvc) == len(v.CVC) && !blank(v.Name)
}

// digits keeps the ASCII digits of s.
func digits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}
