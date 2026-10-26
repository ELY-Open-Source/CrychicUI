package home

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HomeStatus is the home at a glance.
type HomeStatus struct {
	Name           string
	Inside, Energy float64
	Lights, On     int
	Locked         bool
	Alerts         int
}

// HomeSummaryOptions word the overview.
type HomeSummaryOptions struct {
	Inside, Lights, Devices, Energy, Degrees, Usage, Locked, Unlocked, Alerts string
}

// HomeSummary is an overview of the home's state.
func HomeSummary(c *ui.Context, s HomeStatus, o HomeSummaryOptions) ui.Element {
	if strings.TrimSpace(s.Name) == "" || !(s.Inside >= -60 && s.Inside <= 60) || !(s.Energy >= 0 && s.Energy <= 1e6) ||
		s.Lights < 0 || s.On < 0 || s.Alerts < 0 {
		panic(fmt.Sprintf("home: HomeSummary %q needs -60 to 60 degrees, 0 to 1e6 kWh and counts of 0 or more", s.Name))
	}
	def(&o.Inside, "Inside")
	def(&o.Lights, "Lights on")
	def(&o.Devices, "Devices on")
	def(&o.Energy, "Energy today")
	def(&o.Degrees, "%s°")
	def(&o.Usage, "%s kWh")
	def(&o.Locked, "All doors locked")
	def(&o.Unlocked, "A door is unlocked")
	def(&o.Alerts, "%d alerts")
	t := c.Theme()
	status, mark, ink := o.Locked, "shield-check", t.Success
	switch {
	case s.Alerts > 0:
		status, mark, ink = fmt.Sprintf(o.Alerts, s.Alerts), "shield-alert", t.Danger
	case !s.Locked:
		status, mark, ink = o.Unlocked, "lock-open", t.Warning
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "house", t.Accent, t.Space(11), true)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, s.Name).FontSize(t.Rem(1)).Bold().FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(ink).Children(func() {
					c.Key(mark)
					kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.6, Fade: 1}, 0)
					ui.Text(c, status).FontSize(t.Rem(0.76)).Bold().Grow(1).MinWidth(0).SingleLine()
				})
			})
		})
		grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2))
		grid.Children(func() {
			stats := [4][3]string{
				{o.Inside, fmt.Sprintf(o.Degrees, strconv.FormatFloat(math.Round(s.Inside*10)/10, 'f', 1, 64)), "thermometer"},
				{o.Lights, number(c, s.Lights), "lightbulb"},
				{o.Devices, number(c, s.On), "plug"},
				{o.Energy, fmt.Sprintf(o.Usage, strconv.FormatFloat(s.Energy, 'f', 1, 64)), "zap"},
			}
			for i, st := range stats {
				cell := ui.Column(c).MinWidth(0).Gap(t.Space(1)).Padding(t.Space(2.5)).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
				kit.Stagger(c, grid, cell, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
						glyph(c, st[2], 14)
						ui.Text(c, st[0]).FontSize(t.Rem(0.68)).Grow(1).MinWidth(0).SingleLine()
					})
					kit.Numeric(c, st[1]).FontSize(t.Rem(1.1)).Bold().FillWidth().MinWidth(0).SingleLine().
						ContentTransition(ui.ContentTransition{Numeric: true})
				})
			}
		})
	})
	return root
}
