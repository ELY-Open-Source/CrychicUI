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
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Row(c).Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, label := range labels {
			if label == "" {
				panic("project: Labels label must not be empty")
			}
			c.Key(i)
			ink := inks[i%len(inks)]
			b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, ink).TextColor(ink).Checked((*selected)[i]))
			b.Background(b.AnimateColor("bg", ink.Alpha(0.12*unit((*selected)[i])), ui.Snappy))
			b.Children(func() { ui.Text(c, label).FontSize(t.Rem(0.75)) })
			if b.Clicked() {
				(*selected)[i] = !(*selected)[i]
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Labels", "toggled", "index", i, "selected", (*selected)[i])
			}
		}
	})
	return root
}
