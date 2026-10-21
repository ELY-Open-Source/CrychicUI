package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// fabBase returns a floating action without its cue.
func fabBase(c *ui.Context, label, iconName string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(14), t.Space(14)).Center().Radius(999).
		TextColor(t.AccentText).Shadow(0, 2, 6, 0, kit.Shade(0.18)))
	bg := t.Accent
	if b.Pressed() {
		bg = t.AccentPressed
	} else if b.Hovered() {
		bg = t.AccentHover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
		ui.Icon(c, icon.Must(iconName)).Size(16, 16)
	})
	scale := b.Keyframes("arrive", []ui.Keyframe{
		{Value: 0.5},
		{At: 160 * time.Millisecond, Value: 1.08, Ease: ui.EaseOut},
		{At: 240 * time.Millisecond, Value: 1, Ease: ui.EaseOut},
	})
	return b.Scale(scale, scale)
}

// Fab bounces in with a pop cue.
func Fab(c *ui.Context, label, iconName string) ui.Element {
	b := fabBase(c, label, iconName)
	kit.Clicked(b, sound.Pop)
	return b
}
