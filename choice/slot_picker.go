package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SlotPicker is time slot chips; taken ones dim.
func SlotPicker(c *ui.Context, label string, slots, taken []string, selected *string) ui.Element {
	distinct("SlotPicker", slots)
	for _, s := range taken {
		if !slices.Contains(slots, s) {
			panic("choice: SlotPicker taken slot is unknown")
		}
	}
	if label == "" || selected == nil || *selected != "" && (!slices.Contains(slots, *selected) || slices.Contains(taken, *selected)) {
		panic("choice: SlotPicker needs a label and an empty or free selection")
	}
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for _, s := range slots {
			busy := slices.Contains(taken, s)
			c.Key(s)
			b := ui.RadioBase(c, selected, s).Label(s).Disabled(busy).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)
			bg, edge, fg := t.Background, t.Border, t.Text
			switch {
			case *selected == s:
				bg, edge, fg = t.Accent, t.Accent, t.AccentText
			case !busy && b.Hovered():
				bg = t.SurfaceHover
			}
			b.Background(b.AnimateColor("bg", bg, quick)).Border(1, b.AnimateColor("edge", edge, quick)).
				TextColor(b.AnimateColor("fg", fg, quick))
			kit.Press(b).Children(func() {
				text := ui.Text(c, s).FontFeatures("tnum")
				if busy {
					text.Strikethrough()
				}
			})
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().Wrap().Gap(t.Space(2))
	if changed {
		root.SetChanged()
		kit.Log("SlotPicker", "selected", "index", slices.Index(slots, *selected))
	}
	return root
}
