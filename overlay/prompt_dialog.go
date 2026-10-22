package overlay

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// PromptDialogOptions tune a PromptDialog.
type PromptDialogOptions struct {
	Label, Cancel, Submit string
}

// PromptDialog asks for one line; value changes only on submit.
func PromptDialog(c *ui.Context, title string, value *string, open *bool, o PromptDialogOptions) (root, send ui.Element) {
	if title == "" || value == nil || open == nil {
		panic(fmt.Sprintf("overlay: PromptDialog %q needs a value and open", title))
	}
	def(&o.Label, title)
	def(&o.Cancel, "Cancel")
	def(&o.Submit, "Submit")
	return ask(c, asking{component: "PromptDialog", title: title, label: o.Label, cancel: o.Cancel, submit: o.Submit}, value, open)
}
