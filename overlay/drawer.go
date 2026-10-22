package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DrawerOptions tune a Drawer.
type DrawerOptions struct {
	// Close names the close button.
	Close string
}

// Drawer is a side panel sliding in from the end.
func Drawer(c *ui.Context, title string, open *bool, o DrawerOptions, content func()) ui.Element {
	if title == "" || open == nil || content == nil {
		panic(fmt.Sprintf("overlay: Drawer %q needs open and content", title))
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	w, from := t.Space(80), t.Space(80)
	if c.Direction() == ui.RightToLeft {
		from = -from
	}
	place := func(back ui.Element) ui.Element {
		back.Justify(ui.End)
		return ui.Column(c).Width(w).FillHeight().Padding(t.Space(5)).Gap(t.Space(3))
	}
	if side(c, open, title, place, ui.Pose{X: from}, func() {
		if heading(c, title, o.Close).Clicked() {
			*open = false
			root.SetChanged()
		}
		content()
	}) {
		root.SetChanged()
	}
	switch opened, closed := moved(root, *open, "Drawer"); {
	case opened:
		sound.Play(sound.Open)
	case closed:
		sound.Play(sound.Close)
	}
	return root
}
