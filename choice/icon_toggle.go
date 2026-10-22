package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// iconToggle is a toggle whose icon swaps and lifts.
func iconToggle(c *ui.Context, component, label, off, on string, value *bool, lift float32) ui.Element {
	if label == "" || value == nil {
		panic("choice: " + component + " needs a label and state")
	}
	glyphs := [2]*ui.SVG{icon.Must(off), icon.Must(on)}
	t := c.Theme()
	b := ui.ToggleBase(c, value).Label(label).Size(t.Space(8), t.Space(8)).Radius(t.Radius)
	shown := *value
	bg, fg := ui.Color{}, t.TextMuted
	switch {
	case shown:
		bg, fg = t.Accent.Alpha(0.12), t.Accent
	case b.Hovered():
		bg, fg = t.SurfaceHover, t.Text
	}
	b.Background(b.AnimateColor("bg", bg, quick)).TextColor(b.AnimateColor("fg", fg, quick))
	kit.Press(b).Children(func() {
		ui.Icon(c, glyphs[int(unit(shown))]).Size(16, 16).ContentTransition(ui.ContentTransition{}).
			AnimatePose(ui.Pose{Y: -lift * unit(shown)}, ui.Bouncy)
	})
	if kit.Changed(b, sound.Toggle) {
		kit.Log(component, "toggled", "on", *value)
	}
	return b
}

// IconToggle is a pressable icon that swaps when toggled.
func IconToggle(c *ui.Context, label, offIcon, onIcon string, on *bool) ui.Element {
	return iconToggle(c, "IconToggle", label, offIcon, onIcon, on, 0)
}
