package overlay

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
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
	t := c.Theme()
	root = ui.Box(c)
	draft := ui.Local(root, "draft", func() string { return *value })
	was := ui.Local(root, "was", func() bool { return false })
	if *open && !*was {
		*draft = *value
	}
	*was = *open
	dialog(c, open, t.Space(100), func() {
		ui.Text(c, title).Bold().FontSize(t.Rem(1.1))
		field := ui.TextInput(c, draft).Label(o.Label).FillWidth().AutoFocus()
		ok := strings.TrimSpace(*draft) != ""
		submitted := false
		ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			if button.Ghost(c, o.Cancel).Clicked() {
				*open = false
			}
			send = primary(c, o.Submit, t.Accent, t.AccentText).Disabled(!ok)
			submitted = ok && (send.Clicked() || field.Submitted())
		})
		if submitted {
			*value, *open = *draft, false
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("PromptDialog", "submitted", "length", len(*value))
		}
	})
	moved(root, *open, "PromptDialog")
	return root, send
}
