package geo

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Step is one route leg; Turn names an icon.
type Step struct {
	Text, Turn string
	Meters     float64
	Minutes    int
}

// DirectionsOptions word the route and its arrival.
type DirectionsOptions struct {
	Title, Time, Hours, Arrive string
	Imperial                   bool
	Units                      Units
}

// Directions lists a route's steps on a rail, with totals.
func Directions(c *ui.Context, steps []Step, o DirectionsOptions) ui.Element {
	def(&o.Title, "Directions")
	def(&o.Time, "%d min")
	def(&o.Hours, "%d h %d min")
	def(&o.Arrive, "You have arrived")
	o.Units.fill()
	if len(steps) == 0 {
		panic(fmt.Sprintf("geo: Directions %q needs steps", o.Title))
	}
	total, minutes := 0.0, 0
	for _, s := range steps {
		if strings.TrimSpace(s.Text) == "" || strings.TrimSpace(s.Turn) == "" || s.Minutes < 0 || s.Minutes > 100000 {
			panic(fmt.Sprintf("geo: Directions step %q needs text, a turn icon and 0 to 100,000 minutes", s.Text))
		}
		meters("Directions", s.Text, s.Meters)
		total += s.Meters
		minutes += s.Minutes
	}
	meters("Directions", o.Title, total)
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "route", t.Accent, t.Space(9))
			ui.Text(c, o.Title).FontSize(t.Rem(0.92)).Bold().Grow(1).MinWidth(0).Wrap()
		})
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
			chip := func(mark, text string) {
				ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.SurfaceHover).Children(func() {
					glyph(c, mark, 12).TextColor(t.TextMuted)
					ui.Text(c, text).FontSize(t.Rem(0.72)).Bold().Shrink(0).MaxWidth(t.Space(56)).SingleLine()
				})
			}
			chip("navigation", span(total, o.Imperial, o.Units))
			chip("clock", duration(minutes, o.Time, o.Hours))
		})
		list := ui.Column(c).FillWidth().Role(ui.RoleList)
		list.Children(func() {
			row := func(i int, mark, text, meta string, ink ui.Color, last bool) {
				c.Key(i)
				r := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
				kit.Stagger(c, list, r, i).Children(func() {
					ui.Column(c).Width(t.Space(8)).Shrink(0).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(t.Space(8), t.Space(8)).Center().Radius(999).Background(ink.Alpha(0.12)).TextColor(ink).Children(func() {
							glyph(c, mark, 16)
						})
						if !last {
							ui.Box(c).Width(t.Space(0.5)).Grow(1).MinHeight(t.Space(2)).Background(t.Border)
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(1.5), 0, t.Space(3), 0).Children(func() {
						ui.Text(c, text).FontSize(t.Rem(0.84)).FillWidth().Wrap()
						if meta != "" {
							kit.Muted(c, meta).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
						}
					})
				})
			}
			for i, s := range steps {
				meta := ""
				if s.Meters > 0 {
					meta = span(s.Meters, o.Imperial, o.Units)
				}
				ink := t.Accent
				if i > 0 {
					ink = t.TextMuted
				}
				row(i, s.Turn, s.Text, meta, ink, false)
			}
			row(len(steps), "flag", o.Arrive, "", t.Danger, true)
		})
	})
}
