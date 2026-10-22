package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ImageItem is a labeled picture tile.
type ImageItem struct {
	Label string
	Image ui.ImageSource
}

// ImageChoice is picture tiles; the chosen one rings.
func ImageChoice(c *ui.Context, label string, items []ImageItem, selected *int) ui.Element {
	labels := make([]string, len(items))
	for i, it := range items {
		if it.Image == nil {
			panic("choice: ImageChoice items need images")
		}
		labels[i] = it.Label
	}
	distinct("ImageChoice", labels)
	if label == "" || selected == nil || *selected < 0 || *selected >= len(items) {
		panic("choice: ImageChoice needs a label and a valid selection")
	}
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for i, it := range items {
			c.Key(it.Label)
			b := ui.RadioBase(c, selected, i).Label(it.Label).Column().Gap(t.Space(1.5)).AlignItems(ui.Center)
			ring := t.Accent.Alpha(0)
			if *selected == i {
				ring = t.Accent
			} else if b.Hovered() {
				ring = t.Border
			}
			kit.Press(b).Children(func() {
				frame := ui.Box(c).Padding(t.Space(0.5)).Radius(t.Radius + t.Space(1))
				frame.Border(2, frame.AnimateColor("ring", ring, quick)).Children(func() {
					ui.Box(c).Size(t.Space(24), t.Space(16)).Radius(t.Radius).Clip().Background(t.Surface).Children(func() {
						ui.Image(c, it.Image).Fill().Fit(ui.Cover)
					})
				})
				ui.Text(c, it.Label).FontSize(t.Rem(0.85))
			})
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().Wrap().Gap(t.Space(3))
	if changed {
		root.SetChanged()
		kit.Log("ImageChoice", "selected", "index", *selected)
	}
	return root
}
