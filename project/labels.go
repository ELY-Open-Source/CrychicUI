package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Labels toggles colored label picker; cue Toggle.
func Labels(c *ui.Context, labels []string, selected *[]bool) ui.Element {
	if len(labels) == 0 || selected == nil || len(*selected) != len(labels) {
		panic("project: Labels needs labels and selection")
	}
	unique("Labels", labels)
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, label := range labels {
			ink, on := labelInk(t, label), (*selected)[i]
			c.Key(label)
			b := pill(c, label, on, ink, func() {
				ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(ink)
				ui.Text(c, label).FontSize(t.Rem(0.75)).Bold().Shrink(0).MaxWidth(t.Space(40)).SingleLine()
				if on {
					c.Key("check")
					kit.Arrive(glyph(c, "check", 12), ui.Pose{Scale: 0.4, Fade: 1}, 0)
				}
			})
			kit.Stagger(c, root, b, i)
			if b.Clicked() {
				(*selected)[i] = !on
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Labels", "toggled", "index", i, "selected", !on)
			}
		}
	})
	return root
}
