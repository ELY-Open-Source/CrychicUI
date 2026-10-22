package form

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HelpPopover is a "?" button opening help in a popover.
func HelpPopover(c *ui.Context, label, help string) ui.Element {
	if label == "" || help == "" {
		panic("form: HelpPopover needs a label and help")
	}
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999))
	open := ui.Local(b, "open", func() bool { return false })
	ink := t.TextMuted
	if *open || b.Hovered() {
		ink = t.Accent
	}
	b.TextColor(b.AnimateColor("ink", ink, ui.Snappy)).Expanded(*open).Children(func() {
		ui.Icon(c, icon.Must("circle-question-mark")).Size(16, 16)
	})
	if b.Clicked() {
		if *open = !*open; *open {
			sound.Play(sound.Open)
			kit.Log("HelpPopover", "open")
		}
	}
	ui.Popover(c, b, open, func() {
		ui.Text(c, help).MaxWidth(t.Rem(18))
	})
	return b
}
