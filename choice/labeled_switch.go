package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LabeledSwitchOptions tune a LabeledSwitch.
type LabeledSwitchOptions struct {
	// On and Off show inside the track; "On", "Off".
	On, Off string
}

// LabeledSwitch is a switch whose inner text cross-fades.
func LabeledSwitch(c *ui.Context, label string, on *bool, o LabeledSwitchOptions) ui.Element {
	if label == "" || on == nil {
		panic("choice: LabeledSwitch needs a label and state")
	}
	if o.On == "" {
		o.On = "On"
	}
	if o.Off == "" {
		o.Off = "Off"
	}
	t := c.Theme()
	root := ui.SwitchBase(c, on).Label(label).Gap(t.Space(2))
	shown := *on
	root.Children(func() {
		w, h, in := t.Space(15), t.Space(6.5), t.Space(0.75)
		tr := ui.Box(c).Size(w, h).Radius(h / 2).Shrink(0).Clip()
		off := t.Border.Mix(t.Text, 0.15)
		tr.Background(tr.AnimateColor("track", off.Mix(t.Accent, unit(shown)), ui.Snappy)).Children(func() {
			text, align := o.Off, ui.End
			if shown {
				text, align = o.On, ui.Start
			}
			ui.Row(c).Absolute().Top(0).Bottom(0).Left(t.Space(2)).Right(t.Space(2)).Justify(align).Children(func() {
				ui.Text(c, text).FontSize(t.Rem(0.75)).Bold().TextColor(t.AccentText).ContentTransition(ui.ContentTransition{})
			})
			knob := ui.Box(c).Absolute().Top(in).Left(in).Size(h-2*in, h-2*in).Radius(h).Background(t.AccentText).
				Shadow(0, 1, 3, 0, kit.Shade(0.25))
			knob.Translate(knob.AnimateTo("x", (w-h)*unit(shown), ui.Snappy), 0)
		})
		ui.Text(c, label)
	})
	if kit.Changed(root, sound.Toggle) {
		kit.Log("LabeledSwitch", "toggled", "on", *on)
	}
	return root
}
