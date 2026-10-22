package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FullscreenModalOptions tune a FullscreenModal.
type FullscreenModalOptions struct {
	// Close names the close button.
	Close string
}

// FullscreenModal fills the window and slides up.
func FullscreenModal(c *ui.Context, title string, open *bool, o FullscreenModalOptions, content func()) ui.Element {
	if title == "" || open == nil || content == nil {
		panic(fmt.Sprintf("overlay: FullscreenModal %q needs open and content", title))
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	place := func(back ui.Element) ui.Element {
		return ui.Column(c).Grow(1).FillHeight().Padding(t.Space(6)).Gap(t.Space(4))
	}
	if side(c, open, title, place, ui.Pose{Y: t.Space(24), Fade: 1}, func() {
		if heading(c, title, o.Close).Clicked() {
			*open = false
			root.SetChanged()
		}
		ui.Column(c).FillWidth().Grow(1).Children(content)
	}) {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "FullscreenModal"); opened {
		sound.Play(sound.Open)
	}
	return root
}
