package choice

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NumberStepperOptions tune a NumberStepper.
type NumberStepperOptions struct {
	// Format writes the value; "%.0f" when empty.
	Format string
	// Decrease and Increase format the button names; "Decrease %s".
	Decrease, Increase string
}

// NumberStepper is minus and plus around a rolling value.
func NumberStepper(c *ui.Context, label string, value *float64, lo, hi, step float64, o NumberStepperOptions) ui.Element {
	sliderArgs("NumberStepper", label, value, lo, hi, step)
	if step == 0 {
		panic("choice: NumberStepper needs a step")
	}
	if o.Format == "" {
		o.Format = "%.0f"
	}
	if o.Decrease == "" {
		o.Decrease = "Decrease %s"
	}
	if o.Increase == "" {
		o.Increase = "Increase %s"
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).Padding(t.Space(0.5)).Radius(t.Radius).Background(t.Surface).
		Border(1, t.Border).Label(label).FocusGroup(ui.Horizontal)
	down := ui.Local(root, "down", func() bool { return false })
	set := func(v float64) {
		v = snap(min(hi, max(lo, v)), lo, step)
		if v > hi {
			v -= step
		}
		if v != *value {
			*down = v < *value
			*value = v
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("NumberStepper", "changed", "value", v)
		}
	}
	button := func(name, glyph string, off bool) ui.Element {
		b := ui.ButtonBase(c).Label(fmt.Sprintf(name, label)).Disabled(off).Size(t.Space(7), t.Space(7)).Radius(t.Radius - t.Space(0.5))
		bg := ui.Color{}
		if b.Hovered() {
			bg = t.SurfaceHover
		}
		b.Background(b.AnimateColor("bg", bg, quick))
		return kit.Press(b).Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16) })
	}
	return root.Children(func() {
		c.Key("minus")
		if button(o.Decrease, "minus", *value <= lo).Clicked() || root.Shortcut(0, ui.KeyDown) {
			set(*value - step)
		}
		kit.Numeric(c, fmt.Sprintf(o.Format, *value)).ContentTransition(ui.ContentTransition{Numeric: true, Down: *down}).
			MinWidth(t.Space(10)).TextAlign(ui.Center).Bold()
		c.Key("plus")
		if button(o.Increase, "plus", *value >= hi).Clicked() || root.Shortcut(0, ui.KeyUp) {
			set(*value + step)
		}
	})
}
