package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Swatch is a named color.
type Swatch struct {
	Name  string
	Color ui.Color
}

// ColorChoice is swatch radios; the chosen one rings.
func ColorChoice(c *ui.Context, label string, swatches []Swatch, selected *int) ui.Element {
	names := make([]string, len(swatches))
	for i, s := range swatches {
		names[i] = s.Name
	}
	distinct("ColorChoice", names)
	if label == "" || selected == nil || *selected < 0 || *selected >= len(swatches) {
		panic("choice: ColorChoice needs a label and a valid selection")
	}
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for i, s := range swatches {
			c.Key(s.Name)
			b := ui.RadioBase(c, selected, i).Label(s.Name).Tooltip(s.Name).Size(t.Space(8), t.Space(8)).
				Padding(t.Space(0.75)).Radius(999)
			ring := t.Accent.Alpha(0)
			if *selected == i {
				ring = t.Accent
			} else if b.Hovered() {
				ring = t.Border
			}
			b.Border(2, b.AnimateColor("ring", ring, quick))
			kit.Press(b).Children(func() {
				ui.Box(c).Fill().Radius(999).Background(s.Color).Border(1, kit.Shade(0.1))
			})
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().Wrap().Gap(t.Space(2))
	if changed {
		root.SetChanged()
		kit.Log("ColorChoice", "selected", "index", *selected)
	}
	return root
}
