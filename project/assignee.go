package project

import (
	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AssigneeOptions word the picker.
type AssigneeOptions struct{ Title, None string }

// Assignee selects a person with avatars; cue Select.
func Assignee(c *ui.Context, people []string, selected *int, o AssigneeOptions) ui.Element {
	if len(people) == 0 || selected == nil || *selected < -1 || *selected >= len(people) {
		panic("project: Assignee needs people and selection")
	}
	unique("Assignee", people)
	def(&o.Title, "Assignee")
	def(&o.None, "Unassigned")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "user-round", o.Title, func() {
			name := o.None
			if *selected >= 0 {
				name = people[*selected]
			}
			kit.Chip(c, name, t.Accent)
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1))
		list.Children(func() {
			for i, name := range people {
				on := i == *selected
				c.Key(name)
				b := kit.Press(ui.ButtonBase(c).Label(name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(2)).
					Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(0.7*unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						display.Avatar(c, name, nil).Size(t.Space(7), t.Space(7)).FontSize(t.Rem(0.65))
						ui.Text(c, name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine().
							TextColor(b.AnimateColor("ink", t.Text.Mix(t.Accent, unit(on)), ui.Snappy))
						if on {
							c.Key("check")
							kit.Arrive(glyph(c, "check", 14).TextColor(t.Accent), ui.Pose{Scale: 0.4, Fade: 1}, 0)
						}
					})
				})
				if b.Clicked() && !on {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Assignee", "selected", "index", i)
				}
			}
		})
	})
	return root
}
