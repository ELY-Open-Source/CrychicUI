package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Template is one project template; Icon is optional.
type Template struct {
	Title, Description, Icon string
}

// Templates shows template cards; cue Select.
func Templates(c *ui.Context, templates []Template, selected *int) ui.Element {
	if len(templates) == 0 || selected == nil || *selected < -1 || *selected >= len(templates) {
		panic("project: Templates needs templates and a selection in range")
	}
	titles := make([]string, len(templates))
	for i, template := range templates {
		if template.Description == "" {
			panic(fmt.Sprintf("project: Templates %q needs a description", template.Title))
		}
		titles[i] = template.Title
	}
	unique("Templates", titles)
	t := c.Theme()
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, template := range templates {
			c.Key(template.Title)
			on := i == *selected
			mark := template.Icon
			if mark == "" {
				mark = "layout-template"
			}
			b := kit.Press(ui.ButtonBase(c).Label(template.Title).Column().AlignItems(ui.Stretch).Gap(t.Space(2)).
				Padding(t.Space(3)).Radius(t.Radius).Checked(on))
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, 0.5*unit(b.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			kit.Stagger(c, root, b, i).Children(func() {
				ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
					tile(c, mark, t.Accent)
					ui.Box(c).Grow(1)
					c.Key(on)
					if on {
						badge := ui.Box(c).Size(t.Space(5), t.Space(5)).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText)
						kit.Arrive(badge, ui.Pose{Scale: 0.4, Fade: 1}, 0).Children(func() {
							glyph(c, "check", 12)
						})
					}
				})
				ui.Text(c, template.Title).FontSize(t.Rem(0.85)).Bold().Wrap().FillWidth()
				kit.Muted(c, template.Description).FontSize(t.Rem(0.72)).Wrap().FillWidth()
			})
			if b.Clicked() && !on {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Templates", "selected", "index", i)
			}
		}
	})
	return root
}
