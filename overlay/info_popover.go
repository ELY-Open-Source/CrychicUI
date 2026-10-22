package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// InfoPopover is an "i" button opening a short text.
func InfoPopover(c *ui.Context, label, text string, open *bool) ui.Element {
	if label == "" || text == "" || open == nil {
		panic(fmt.Sprintf("overlay: InfoPopover %q needs text and open", label))
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Size(t.Space(7), t.Space(7)).Center().Radius(999)).Expanded(*open)
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || b.Hovered())), ui.Snappy)).Children(func() {
			ui.Icon(c, icon.Must("info")).Size(16, 16).TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(*open)), ui.Snappy))
		})
		if b.Clicked() {
			*open = !*open
			root.SetChanged()
		}
		if ui.PopoverBase(c, b, open, func(panel ui.Element) {
			panel.Margin(t.Space(1.5), 0, 0, 0).MaxWidth(t.Space(64)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).
				Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Fade: 1}, 0)
			ui.Text(c, text)
		}).Dismissed() {
			root.SetChanged()
		}
	})
	if opened, _ := moved(root, *open, "InfoPopover"); opened {
		sound.Play(sound.Open)
	}
	return root
}
