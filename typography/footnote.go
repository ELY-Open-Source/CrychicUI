package typography

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Footnote is a superscript number with a popover note.
func Footnote(c *ui.Context, number int, note string, open *bool) ui.Element {
	if number <= 0 || open == nil {
		panic("typography: Footnote needs a number and open state")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start)
	return root.Children(func() {
		anchor := ui.ButtonBase(c).Label("Footnote " + strconv.Itoa(number)).TextColor(t.Accent)
		kit.Press(anchor).Children(func() {
			ui.Text(c, strconv.Itoa(number)).FontSize(t.Rem(0.75)).Translate(0, -4)
		})
		if kit.Clicked(anchor, sound.Open) {
			*open = !*open
			root.SetChanged()
			kit.Log("Footnote", "toggled", "open", *open)
		}
		panel := ui.Popover(c, anchor, open, func() {
			ui.Text(c, note).MaxWidth(t.Rem(24))
		})
		if panel.Dismissed() {
			root.SetChanged()
			kit.Log("Footnote", "closed")
		}
	})
}
