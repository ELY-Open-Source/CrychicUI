package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuickViewOptions tune a QuickView.
type QuickViewOptions struct {
	// Close names the close button.
	Close string
}

// QuickView is a preview card scaling up over its anchor.
func QuickView(c *ui.Context, anchor ui.Element, title string, open *bool, o QuickViewOptions, content func()) ui.Element {
	if title == "" || open == nil || content == nil {
		panic(fmt.Sprintf("overlay: QuickView %q needs open and content", title))
	}
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	if ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.AttachTo(anchor, ui.AnchorCenter, ui.AnchorCenter).Width(t.Space(90)).Padding(t.Space(4)).Radius(t.Radius*1.5).
			Background(t.Background).Border(1, t.Border).Shadow(0, 12, 32, 0, kit.Shade(0.2))
		kit.Arrive(panel, ui.Pose{Scale: 0.85, Fade: 1}, 0)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
			if heading(c, title, o.Close).Clicked() {
				*open = false
				root.SetChanged()
			}
			content()
		})
	}).Dismissed() {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "QuickView"); opened {
		sound.Play(sound.Open)
	}
	return root
}
