package overlay

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CalloutPopover is a popover edged in a tone.
func CalloutPopover(c *ui.Context, anchor ui.Element, open *bool, tone feedback.Tone, title, text string) ui.Element {
	if open == nil || title == "" {
		panic(fmt.Sprintf("overlay: CalloutPopover %q needs open", title))
	}
	t := c.Theme()
	ink, glyph := toned(t, tone)
	root := ui.Box(c)
	if ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.Margin(t.Space(2), 0, 0, 0).TransformOrigin(0.5, 0)
		kit.Arrive(panel, ui.Pose{Scale: 0.92, Fade: 1}, 0)
		ui.Row(c).Width(t.Space(76)).Radius(t.Radius).Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).
			Clip().Role(ui.RoleGroup).Label(title).Children(func() {
			edge := ui.Box(c).Width(t.Space(1)).FillHeight().Background(ink)
			kit.Arrive(edge, ui.Pose{Fade: 1}, 120*time.Millisecond)
			ui.Column(c).Grow(1).Padding(t.Space(3), t.Space(4)).Gap(t.Space(1.5)).Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(ink)
					ui.Text(c, title).Bold()
				})
				if text != "" {
					kit.Muted(c, text)
				}
			})
		})
	}).Dismissed() {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "CalloutPopover"); opened {
		sound.Play(sound.Open)
	}
	return root
}
