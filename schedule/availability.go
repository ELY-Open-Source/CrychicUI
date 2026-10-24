package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AvailabilityOptions word the grid; cells run day by day.
type AvailabilityOptions struct {
	Free, Busy, Mine string
	Days, Hours      []string
}

// Availability picks free cells; busy ones stay shut; cue Toggle.
func Availability(c *ui.Context, cells, selected *[]bool, o AvailabilityOptions) ui.Element {
	if cells == nil || selected == nil || len(*cells) == 0 || len(*cells) != len(*selected) {
		panic("schedule: Availability needs cells and selection")
	}
	headed := len(o.Days) > 0
	if headed && len(o.Days)*len(o.Hours) != len(*cells) {
		panic(fmt.Sprintf("schedule: Availability has %d cells for %d days and %d hours", len(*cells), len(o.Days), len(o.Hours)))
	}
	seen := map[string]bool{}
	for _, name := range append(append([]string(nil), o.Days...), o.Hours...) {
		if name == "" || seen[name] {
			panic(fmt.Sprintf("schedule: Availability heading %q is empty or repeated", name))
		}
		seen[name] = true
	}
	def(&o.Free, "Free")
	def(&o.Busy, "Busy")
	def(&o.Mine, "Your pick")
	t := c.Theme()
	cols := 7
	if headed {
		cols = len(o.Days)
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
			if headed {
				ui.Column(c).Width(t.Space(11)).Shrink(0).Gap(t.Space(1)).Padding(t.Space(6), 0, 0, 0).Children(func() {
					for _, hour := range o.Hours {
						ui.Box(c).Height(t.Space(7)).Center().Children(func() { kit.Muted(c, hour).FontSize(t.Rem(0.65)) })
					}
				})
			}
			grid := ui.Grid(c).Columns(cols).Grow(1).MinWidth(0).Gap(t.Space(1))
			grid.Children(func() {
				if headed {
					for _, day := range o.Days {
						ui.Box(c).Height(t.Space(5)).Center().Children(func() { ui.Text(c, day).FontSize(t.Rem(0.7)).Bold() })
					}
				}
				for k := range *cells {
					i := k
					if headed {
						i = (k%cols)*len(o.Hours) + k/cols
					}
					availCell(c, root, grid, cells, selected, i, k/cols, o)
				}
			})
		})
		ui.Row(c).Gap(t.Space(3)).Children(func() {
			for _, item := range []struct {
				name string
				ink  ui.Color
			}{{o.Free, t.Success.Alpha(0.35)}, {o.Busy, t.Border}, {o.Mine, t.Accent}} {
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(t.Space(0.5)).Background(item.ink)
					kit.Muted(c, item.name).FontSize(t.Rem(0.7)).Shrink(0)
				})
			}
		})
	})
	return root
}

// availCell is cell i, built in grid row row.
func availCell(c *ui.Context, root, grid ui.Element, cells, selected *[]bool, i, row int, o AvailabilityOptions) {
	t := c.Theme()
	free, on := (*cells)[i], (*selected)[i]
	label := o.Busy
	if free {
		label = o.Free
	}
	key := any(i)
	if len(o.Days) > 0 {
		day, hour := o.Days[i/len(o.Hours)], o.Hours[i%len(o.Hours)]
		label, key = day+" "+hour, fmt.Sprintf("%q %q", day, hour)
	}
	c.Key(key)
	b := kit.Press(ui.ButtonBase(c).Label(label).Height(t.Space(7)).Center().Radius(t.Space(1)).Checked(on).Disabled(!free))
	bg := t.SurfaceHover
	if free {
		bg = t.Success.Alpha(0.12).Mix(t.Success.Alpha(0.24), unit(b.Hovered())).Mix(t.Accent, unit(on))
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(t.AccentText)
	kit.Stagger(c, grid, b, row)
	b.Children(func() {
		if on {
			c.Key("check")
			kit.Arrive(glyph(c, "check", 14), ui.Pose{Scale: 0.4, Fade: 1}, 0)
		}
	})
	if b.Clicked() && free {
		(*selected)[i] = !on
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Availability", "toggled", "index", i, "selected", !on)
	}
}
