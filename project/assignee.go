package project

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Assignee selects a person with avatars; cue Select.
func Assignee(c *ui.Context, people []string, selected *int) ui.Element {
	if len(people) == 0 || selected == nil || *selected < -1 || *selected >= len(people) {
		panic("project: Assignee needs people and selection")
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, name := range people {
			if name == "" {
				panic("project: Assignee person must not be empty")
			}
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(name).Row().Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(i == *selected))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			b.Children(func() {
				ui.Box(c).Size(t.Space(5), t.Space(5)).Radius(999).Background(t.Accent.Alpha(0.15)).Center().Children(func() {
					first, _ := utf8.DecodeRuneInString(name)
					ui.Text(c, string(first)).TextColor(t.Accent)
				})
				ui.Text(c, name).FontSize(t.Rem(0.75))
			})
			if b.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Assignee", "selected", "index", i)
			}
		}
	})
	return root
}
