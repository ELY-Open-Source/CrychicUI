package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChipButton is a filter chip toggling with a growing check.
func ChipButton(c *ui.Context, label string, on *bool) ui.Element {
	if on == nil {
		panic("button: ChipButton needs state")
	}
	t := c.Theme()
	root := textBase(c, label).Radius(999).Checked(*on).Border(1, t.Border)
	bg, fg := t.Surface, t.Text
	if *on {
		bg, fg = t.Accent.Alpha(0.16), t.Accent
	}
	root.Background(root.AnimateColor("bg", bg, ui.Snappy)).TextColor(fg).Children(func() {
		check := ui.Icon(c, icon.Must("check")).Size(16, 16)
		scale := check.AnimateTo("show", unit(*on), ui.Snappy)
		check.Scale(scale, scale).Opacity(scale)
		ui.Text(c, label).SingleLine()
	})
	if root.Clicked() {
		*on = !*on
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("ChipButton", "toggled", "on", *on)
	}
	return root
}
