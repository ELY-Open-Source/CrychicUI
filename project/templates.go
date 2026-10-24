package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Template is one project template.
type Template struct {
	Title, Description string
}

// Templates shows template cards; cue Select.
func Templates(c *ui.Context, templates []Template, selected *int) ui.Element {
	if len(templates) == 0 || selected == nil || *selected < -1 || *selected >= len(templates) {
		panic("project: Templates needs templates and selection")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, template := range templates {
			if template.Title == "" || template.Description == "" {
				panic("project: Templates has an invalid template")
			}
			c.Key(i)
			card := kit.Press(ui.ButtonBase(c).Label(template.Title).Column().Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Checked(i == *selected))
			card.Background(card.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			card.Children(func() {
				ui.Text(c, template.Title).Bold()
				ui.Text(c, template.Description).FontSize(t.Rem(0.75)).Wrap()
			})
			if card.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Templates", "selected", "index", i)
			}
		}
	})
	return root
}
