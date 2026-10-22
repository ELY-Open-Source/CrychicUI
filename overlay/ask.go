package overlay

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// asking is how a one-line dialog reads and sends.
type asking struct {
	component, title, label, cancel, submit string
	// rename selects the name's stem and needs a change.
	rename bool
}

// ask is a one-line dialog submitting a draft.
func ask(c *ui.Context, a asking, value *string, open *bool) (root, send ui.Element) {
	t := c.Theme()
	root = ui.Box(c)
	draft := ui.Local(root, "draft", func() string { return *value })
	was := ui.Local(root, "was", func() bool { return false })
	opening := *open && !*was
	if opening {
		*draft = *value
	}
	*was = *open
	if dialog(c, open, t.Space(100), func() {
		ui.Text(c, a.title).Bold().FontSize(t.Rem(1.1))
		field := ui.TextInput(c, draft).Label(a.label).FillWidth().AutoFocus()
		if opening && a.rename {
			end := utf8.RuneCountInString(*draft)
			if dot := strings.LastIndex(*draft, "."); dot > 0 {
				end = utf8.RuneCountInString((*draft)[:dot])
			}
			field.SetTextSelection(0, end)
		}
		ok := strings.TrimSpace(*draft) != "" && (!a.rename || *draft != *value)
		submitted := false
		ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			if button.Ghost(c, a.cancel).Clicked() {
				*open = false
				root.SetChanged()
			}
			send = primary(c, a.submit, t.Accent, t.AccentText).Disabled(!ok)
			submitted = ok && (send.Clicked() || field.Submitted())
		})
		if submitted {
			*value, *open = *draft, false
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log(a.component, "submitted", "length", len(*value))
		}
	}).Dismissed() {
		root.SetChanged()
	}
	moved(root, *open, a.component)
	return root, send
}
