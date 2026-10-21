package input

import "github.com/egoist/mygo/ui"

// PinInputOptions tune a PinInput.
type PinInputOptions struct {
	// Label names the hidden input.
	Label string
}

// PinInput is masked digit boxes with popping dots.
func PinInput(c *ui.Context, value *string, length int, o PinInputOptions) ui.Element {
	if o.Label == "" {
		o.Label = "PIN"
	}
	return codeInput(c, o.Label, value, length, true, false)
}
