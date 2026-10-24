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
	root := ui.Row(c).Wrap().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, value := range []TaskStatus{Todo, Doing, Done} {
			name := statusText(value, o.Todo, o.Doing, o.Done)
			icon, ink := statusLook(t, value)
			on := *status == value
			c.Key(int(value))
			b := pill(c, name, on, ink, func() {
				c.Key(on)
				kit.Arrive(glyph(c, icon, 14).TextColor(ink), ui.Pose{Scale: 0.6, Fade: 1}, 0)
				ui.Text(c, name).FontSize(t.Rem(0.75)).Bold().Shrink(0)
			})
			kit.Stagger(c, root, b, int(value))
			if b.Clicked() && !on {
				*status = value
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Status", "selected", "status", value)
			}
		}
	})
	return root
}
