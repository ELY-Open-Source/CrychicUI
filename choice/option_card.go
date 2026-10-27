package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// OptionCard is a radio card with an animated border.
func OptionCard(c *ui.Context, title, body, value string, selected *string) ui.Element {
	if selected == nil || title == "" {
		panic("choice: " + who("OptionCard", title) + " needs a title and selection")
	}
	t := c.Theme()
	root := ui.RadioBase(c, selected, value).Label(title).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Radius(t.Radius)
	on := *selected == value
	edge, fill := t.Border, t.Background
	if on {
		edge, fill = t.Accent, t.Accent.Alpha(0.06)
	} else if root.Hovered() {
		edge = t.Border.Mix(t.Text, 0.25)
	}
	root.Background(root.AnimateColor("fill", fill, ui.Snappy)).Border(1, root.AnimateColor("edge", edge, ui.Snappy))
	root.AlignItems(ui.Start).Children(func() {
		dot(c, root, on)
		ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, title).Bold()
			if body != "" {
				kit.Muted(c, body)
			}
		})
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("OptionCard", "selected")
	}
	return root
}
