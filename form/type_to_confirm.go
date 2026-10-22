package form

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TypeToConfirmOptions tune a TypeToConfirm.
type TypeToConfirmOptions struct {
	// Prompt formats the name in; Action labels the button.
	Prompt, Action string
}

// TypeToConfirm enables a destructive action once name is typed exactly.
func TypeToConfirm(c *ui.Context, name string, typed *string, o TypeToConfirmOptions) (root, confirm ui.Element) {
	if name == "" || typed == nil {
		panic("form: TypeToConfirm needs a name and a value")
	}
	def(&o.Prompt, "Type %s to confirm.")
	def(&o.Action, "Delete")
	t := c.Theme()
	prompt := fmt.Sprintf(o.Prompt, name)
	root = ui.Column(c).FillWidth().Gap(t.Space(3))
	root.Children(func() {
		ui.Text(c, prompt)
		edited(root, ui.TextInput(c, typed).Label(prompt).FillWidth())
		confirm = button.Destructive(c, o.Action).Disabled(*typed != name)
	})
	if confirm.Clicked() {
		kit.Log("TypeToConfirm", "confirmed")
	}
	return root, confirm
}
