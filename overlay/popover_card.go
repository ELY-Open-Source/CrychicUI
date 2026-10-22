package overlay

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PopoverCard is an arrowed card scaling from its anchor.
func PopoverCard(c *ui.Context, anchor ui.Element, open *bool, content func()) ui.Element {
	if open == nil || content == nil {
		panic("overlay: PopoverCard needs open and content")
	}
	t := c.Theme()
	root := ui.Box(c)
	if ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.Margin(t.Space(2.5), 0, 0, 0).TransformOrigin(0.5, 0)
		kit.Arrive(panel, ui.Pose{Scale: 0.9, Fade: 1}, 0)
		ui.Column(c).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(3), t.Space(3)).Rotate(45).Margin(0, 0, -t.Space(1.5), 0).Background(t.Background).Border(1, t.Border)
			ui.Column(c).Width(t.Space(72)).Padding(t.Space(4)).Gap(t.Space(2)).Radius(t.Radius).Background(t.Background).
				Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).Children(content)
		})
	}).Dismissed() {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "PopoverCard"); opened {
		sound.Play(sound.Open)
	}
	return root
}
