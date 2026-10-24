package project

import (
	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TaskRow checks a task; the title strikes; cue Toggle.
func TaskRow(c *ui.Context, task Task, done *bool) ui.Element {
	validTask("TaskRow", task)
	if done == nil {
		panic("project: TaskRow needs done")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Role(ui.RoleGroup)
	root.Background(root.AnimateColor("bg", t.SurfaceHover.Alpha(0.6*unit(root.Hovered())), ui.Snappy))
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		c.Key("check")
		check := kit.Press(ui.ButtonBase(c).Label(task.Title).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Radius(999).Checked(*done))
		check.Border(1, check.AnimateColor("line", t.Border.Mix(t.TextMuted, unit(check.Hovered())).Mix(t.Success, unit(*done)), ui.Snappy))
		check.Background(check.AnimateColor("bg", t.Success.Alpha(unit(*done)), ui.Snappy)).TextColor(t.AccentText)
		check.Children(func() {
			if *done {
				c.Key("mark")
				kit.Arrive(glyph(c, "check", 12), ui.Pose{Scale: 0.4, Fade: 1}, 0)
			}
		})
		if check.Clicked() {
			*done = !*done
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("TaskRow", "toggled", "done", *done)
		}
		title := ui.Text(c, task.Title).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine().
			TextColor(root.AnimateColor("ink", t.Text.Mix(t.TextMuted, unit(*done)), ui.Snappy))
		if *done {
			title.Strikethrough()
		}
		if task.Due != "" {
			ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
				glyph(c, "calendar", 12)
				ui.Text(c, task.Due).FontSize(t.Rem(0.7)).Shrink(0)
			})
		}
		if task.Assignee != "" {
			display.Avatar(c, task.Assignee, nil).Size(t.Space(6), t.Space(6)).FontSize(t.Rem(0.6)).Tooltip(task.Assignee)
		}
	})
	return root
}
