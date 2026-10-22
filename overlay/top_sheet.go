package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TopSheetOptions tune a TopSheet.
type TopSheetOptions struct {
	// Close names the close button.
	Close string
}

// TopSheet is a panel dropping from the top.
func TopSheet(c *ui.Context, title string, open *bool, o TopSheetOptions, content func()) ui.Element {
	if title == "" || open == nil || content == nil {
		panic(fmt.Sprintf("overlay: TopSheet %q needs open and content", title))
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	place := func(back ui.Element) ui.Element {
		back.AlignItems(ui.Start)
		return ui.Column(c).Grow(1).Padding(t.Space(5)).Gap(t.Space(3))
	}
	if side(c, open, title, place, ui.Pose{Y: -t.Space(24), Fade: 1}, func() {
		if heading(c, title, o.Close).Clicked() {
			*open = false
			root.SetChanged()
		}
		content()
	}) {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "TopSheet"); opened {
		sound.Play(sound.Open)
	}
	return root
}
