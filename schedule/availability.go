package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AvailabilityOptions are the availability labels.
type AvailabilityOptions struct{ Free, Busy string }

// Availability toggles free and busy cells; cue Toggle.
func Availability(c *ui.Context, cells, selected *[]bool, o AvailabilityOptions) ui.Element {
	if cells == nil || selected == nil || len(*cells) == 0 || len(*cells) != len(*selected) {
		panic("schedule: Availability needs cells and selection")
	}
	def(&o.Free, "Free")
	def(&o.Busy, "Busy")
	t := c.Theme()
	root := ui.Grid(c).Columns(7).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, free := range *cells {
			c.Key(i)
			text, ink := o.Busy, t.Danger
			if free {
				text, ink = o.Free, t.Success
			}
			chosen := (*selected)[i]
			b := kit.Press(ui.ButtonBase(c).Label(text).Padding(t.Space(1)).Radius(t.Space(0.5)).
				Border(1, ink).TextColor(ink).Checked(chosen))
			b.Background(b.AnimateColor("bg", ink.Alpha(0.12*unit(chosen)), ui.Snappy))
			b.Children(func() { ui.Text(c, text).FontSize(t.Rem(0.65)) })
			if b.Clicked() {
				(*selected)[i] = !(*selected)[i]
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Availability", "toggled", "index", i, "selected", (*selected)[i])
			}
		}
	})
	return root
}
