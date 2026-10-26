package desktop

import (
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// Process is a running process; CPU is a percent.
type Process struct {
	Name, Icon string
	CPU        float64
	Memory     int64
}

// TaskManagerOptions word the table; CPU formats a percent.
type TaskManagerOptions struct{ Title, CPU string }

// TaskManager lists processes busiest first, with CPU and memory.
func TaskManager(c *ui.Context, procs []Process, o TaskManagerOptions) ui.Element {
	def(&o.Title, "Processes")
	def(&o.CPU, "%.1f%%")
	names := make([]string, len(procs))
	for i, p := range procs {
		if p.Icon == "" || math.IsNaN(p.CPU) || p.CPU < 0 || p.CPU > 100 || p.Memory < 0 {
			panic(fmt.Sprintf("desktop: TaskManager process %q needs an icon, CPU 0 to 100 and memory of 0 or more", p.Name))
		}
		names[i] = p.Name
	}
	unique("TaskManager", names)
	sorted := slices.Clone(procs)
	slices.SortStableFunc(sorted, func(a, b Process) int {
		switch {
		case a.CPU > b.CPU:
			return -1
		case a.CPU < b.CPU:
			return 1
		}
		return 0
	})
	t := c.Theme()
	root := card(c).Role(ui.RoleTable).Label(o.Title)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(4)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			tile(c, "cpu", t.Accent, t.Space(8))
			ui.Text(c, o.Title).FontSize(t.Rem(0.92)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Chip(c, strconv.Itoa(len(procs)), t.TextMuted)
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1), 0)
		list.Children(func() {
			for i, p := range sorted {
				c.Key(strconv.Quote(p.Name))
				ink := t.Accent
				switch {
				case p.CPU >= 50:
					ink = t.Danger
				case p.CPU >= 20:
					ink = t.Warning
				}
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(4)).Role(ui.RoleRow)
				kit.Stagger(c, list, row, i).Children(func() {
					glyph(c, p.Icon, 16).TextColor(t.TextMuted)
					ui.Text(c, p.Name).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
					ui.Column(c).Width(t.Space(18)).Shrink(0).Gap(t.Space(0.5)).Children(func() {
						kit.Numeric(c, fmt.Sprintf(o.CPU, p.CPU)).Font("monospace").FontSize(t.Rem(0.72)).TextColor(ink).FillWidth().TextAlign(ui.End).SingleLine()
						bar := ui.Box(c).FillWidth().Height(t.Space(0.75)).Radius(999).Background(t.SurfaceHover)
						fill := bar.AnimateTo("fill", float32(p.CPU/100), ui.Smooth)
						bar.Draw(func(pt *ui.Painter, r ui.Rect) {
							if fill > 0 {
								pt.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, ink, r.H/2)
							}
						})
					})
					typography.Bytes(c, p.Memory).Font("monospace").FontSize(t.Rem(0.72)).TextColor(t.TextMuted).Width(t.Space(18)).Shrink(0).TextAlign(ui.End).SingleLine()
				})
			}
		})
	})
}
