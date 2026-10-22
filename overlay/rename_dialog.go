package overlay

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// RenameDialogOptions tune a RenameDialog.
type RenameDialogOptions struct {
	Label, Cancel, Save string
}

// RenameDialog edits a name, its stem selected.
func RenameDialog(c *ui.Context, title string, name *string, open *bool, o RenameDialogOptions) (root, send ui.Element) {
	if title == "" || name == nil || open == nil {
		panic(fmt.Sprintf("overlay: RenameDialog %q needs a name and open", title))
	}
	def(&o.Label, "Name")
	def(&o.Cancel, "Cancel")
	def(&o.Save, "Rename")
	return ask(c, asking{component: "RenameDialog", title: title, label: o.Label, cancel: o.Cancel, submit: o.Save, rename: true}, name, open)
}
