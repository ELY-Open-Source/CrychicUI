package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PriorityOptions are the priority labels.
type PriorityOptions struct{ Urgent, High, Medium, Low string }

// Priority selects urgent to low; cue Select.
func Priority(c *ui.Context, level *PriorityLevel, o PriorityOptions) ui.Element {
	if level == nil || *level > Low {
		panic("project: Priority needs a level")
	}
	def(&o.Urgent, "Urgent")
	def(&o.High, "High")
	def(&o.Medium, "Medium")
	def(&o.Low, "Low")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, value := range []PriorityLevel{Urgent, High, Medium, Low} {
			name := priorityText(value, o.Urgent, o.High, o.Medium, o.Low)
			icon, ink := priorityLook(t, value)
			on := *level == value
			c.Key(int(value))
			b := pill(c, name, on, ink, func() {
				glyph(c, icon, 14).TextColor(ink)
				ui.Text(c, name).FontSize(t.Rem(0.75)).Bold().Shrink(0)
			})
			kit.Stagger(c, root, b, int(value))
			if b.Clicked() && !on {
				*level = value
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Priority", "selected", "level", value)
			}
		}
	})
	return root
}
