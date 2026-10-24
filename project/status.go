package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StatusOptions are the status labels.
type StatusOptions struct{ Todo, Doing, Done string }

// Status selects status with colored dots; cue Select.
func Status(c *ui.Context, status *TaskStatus, o StatusOptions) ui.Element {
	if status == nil || *status > Done {
		panic("project: Status needs a state")
	}
	def(&o.Todo, "To do")
	def(&o.Doing, "Doing")
	def(&o.Done, "Done")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, item := range []struct {
			value TaskStatus
			name  string
			ink   ui.Color
		}{{Todo, o.Todo, t.TextMuted}, {Doing, o.Doing, t.Accent}, {Done, o.Done, t.Success}} {
			b := kit.Press(ui.ButtonBase(c).Label(item.name).Row().Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(*status == item.value))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*status == item.value)), ui.Snappy))
			b.Children(func() {
				ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(item.ink)
				ui.Text(c, item.name).FontSize(t.Rem(0.75))
			})
			if b.Clicked() && *status != item.value {
				*status = item.value
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Status", "selected", "status", *status)
			}
		}
	})
	return root
}
