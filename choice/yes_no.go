package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// YesNoOptions tune a YesNo.
type YesNoOptions struct {
	// Yes and No label the buttons; "Yes", "No".
	Yes, No string
}

// YesNo is joined buttons; -1 unanswered, 0 no, 1 yes.
func YesNo(c *ui.Context, label string, value *int, o YesNoOptions) ui.Element {
	if label == "" || value == nil || *value < -1 || *value > 1 {
		panic("choice: YesNo needs a label and a value from -1 to 1")
	}
	if o.Yes == "" {
		o.Yes = "Yes"
	}
	if o.No == "" {
		o.No = "No"
	}
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for _, v := range []int{1, 0} {
			name, tint := o.Yes, t.Success
			if v == 0 {
				name, tint = o.No, t.Danger
			}
			c.Key(v)
			b := ui.RadioBase(c, value, v).Label(name).Padding(t.Space(1.5), t.Space(5)).Radius(t.Radius - t.Space(0.5))
			bg, fg := ui.Color{}, t.Text
			if *value == v {
				bg, fg = tint.Alpha(0.16), tint
			} else if b.Hovered() {
				bg = t.SurfaceHover
			}
			b.Background(b.AnimateColor("bg", bg, quick)).TextColor(b.AnimateColor("fg", fg, quick))
			kit.Press(b).Children(func() { ui.Text(c, name).Bold() })
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().AlignSelf(ui.Start).Gap(t.Space(0.5)).Padding(t.Space(0.5)).Radius(t.Radius).
		Background(t.Surface).Border(1, t.Border)
	if changed {
		root.SetChanged()
		kit.Log("YesNo", "answered", "value", *value)
	}
	return root
}
