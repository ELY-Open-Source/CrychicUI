package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConfirmDialogOptions tune a ConfirmDialog.
type ConfirmDialogOptions struct {
	Cancel, Confirm string
}

// ConfirmDialog asks with a title, text, Cancel and Confirm.
func ConfirmDialog(c *ui.Context, title, text string, open *bool, o ConfirmDialogOptions) (root, confirm ui.Element) {
	if title == "" || open == nil {
		panic(fmt.Sprintf("overlay: ConfirmDialog %q needs open", title))
	}
	def(&o.Cancel, "Cancel")
	def(&o.Confirm, "Confirm")
	t := c.Theme()
	root = ui.Box(c)
	if dialog(c, open, t.Space(100), func() {
		ui.Text(c, title).Bold().FontSize(t.Rem(1.1))
		kit.Muted(c, text)
		ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			if button.Ghost(c, o.Cancel).Clicked() {
				*open = false
				root.SetChanged()
			}
			confirm = button.Solid(c, o.Confirm)
			if confirm.Clicked() {
				*open = false
				kit.Log("ConfirmDialog", "confirmed")
			}
		})
	}).Dismissed() {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "ConfirmDialog"); opened {
		sound.Play(sound.Open)
	}
	return root, confirm
}
