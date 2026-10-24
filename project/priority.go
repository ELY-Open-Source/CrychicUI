package project

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
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
	root := ui.Row(c).Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, item := range []struct {
			level PriorityLevel
			name  string
			glyph string
		}{{Urgent, o.Urgent, "flame"}, {High, o.High, "arrow-up"}, {Medium, o.Medium, "minus"}, {Low, o.Low, "arrow-down"}} {
			b := kit.Press(ui.ButtonBase(c).Label(item.name).Row().Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(*level == item.level))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*level == item.level)), ui.Snappy))
			b.Children(func() {
				ui.Icon(c, icon.Must(item.glyph)).Size(16, 16)
				ui.Text(c, item.name).FontSize(t.Rem(0.75))
			})
			if b.Clicked() && *level != item.level {
				*level = item.level
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Priority", "selected", "level", *level)
			}
		}
	})
	return root
}
