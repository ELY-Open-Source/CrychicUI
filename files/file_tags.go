package files

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileTags toggles colored file tags; cue Toggle.
func FileTags(c *ui.Context, label string, tags []string, selected *[]bool) ui.Element {
	if label == "" || len(tags) == 0 || selected == nil || len(*selected) != len(tags) {
		panic("files: FileTags needs tags and selection flags")
	}
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Row(c).Wrap().Gap(t.Space(1.5)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, tag := range tags {
			if tag == "" {
				panic("files: FileTags tag must not be empty")
			}
			c.Key(i)
			ink := inks[i%len(inks)]
			b := kit.Press(ui.ButtonBase(c).Label(tag).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, ink).TextColor(ink).Checked((*selected)[i]))
			b.Background(b.AnimateColor("bg", ink.Alpha(0.12*unit((*selected)[i])), ui.Snappy))
			b.Children(func() { ui.Text(c, tag).FontSize(t.Rem(0.8)).SingleLine() })
			if b.Clicked() {
				(*selected)[i] = !(*selected)[i]
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("FileTags", "toggled", "index", i, "selected", (*selected)[i])
			}
		}
	})
	return root
}
