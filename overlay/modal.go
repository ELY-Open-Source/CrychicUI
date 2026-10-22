package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ModalOptions tune a Modal.
type ModalOptions struct {
	// Close names the close button.
	Close string
}

// Modal is a titled dialog scaling in over a dim.
func Modal(c *ui.Context, title string, open *bool, o ModalOptions, content func()) ui.Element {
	if title == "" || open == nil || content == nil {
		panic(fmt.Sprintf("overlay: Modal %q needs open and content", title))
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	if dialog(c, open, t.Space(110), func() {
		if heading(c, title, o.Close).Clicked() {
			*open = false
			root.SetChanged()
		}
		content()
	}).Dismissed() {
		root.SetChanged()
	}
	switch opened, closed := moved(root, *open, "Modal"); {
	case opened:
		sound.Play(sound.Open)
	case closed:
		sound.Play(sound.Close)
	}
	return root
}
