package media

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VolumeOptions tune a Volume.
type VolumeOptions struct {
	Mute, Unmute string
}

// Volume is a mute speaker and a 0–100 slider.
func Volume(c *ui.Context, label string, level *float64, o VolumeOptions) ui.Element {
	def(&o.Mute, "Mute")
	def(&o.Unmute, "Unmute")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
	keep := ui.Local(root, "keep", func() float64 { return 50 })
	return root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Size(t.Space(9), t.Space(9)).Shrink(0).Center().Radius(999))
		s := choice.ValueSlider(c, label, level, 0, 100, choice.ValueSliderOptions{Step: 1, Format: "%.0f%%"}).Grow(1)
		if b.Clicked() {
			if *level > 0 {
				*keep, *level = *level, 0
			} else {
				*level = *keep
			}
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Volume", "muted", "on", *level == 0)
		}
		if s.Changed() {
			root.SetChanged()
		}
		name, glyph := o.Mute, "volume-2"
		switch {
		case *level == 0:
			name, glyph = o.Unmute, "volume-x"
		case *level < 34:
			glyph = "volume"
		case *level < 67:
			glyph = "volume-1"
		}
		b.Label(name).Tooltip(name).Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).
			Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{}) })
	})
}
