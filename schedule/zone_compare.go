package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Zone is one timezone and hour offset from UTC.
type Zone struct {
	Name   string
	Offset int
}

// ZoneCompareOptions word the day shifts.
type ZoneCompareOptions struct{ Next, Prev string }

// ZoneCompare shows one UTC hour across zones, with work hours.
func ZoneCompare(c *ui.Context, hour int, zones []Zone, o ZoneCompareOptions) ui.Element {
	if hour < 0 || hour > 23 || len(zones) == 0 {
		panic("schedule: ZoneCompare has invalid data")
	}
	def(&o.Next, "+1 day")
	def(&o.Prev, "−1 day")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, zone := range zones {
			if zone.Name == "" || zone.Offset < -23 || zone.Offset > 23 {
				panic(fmt.Sprintf("schedule: ZoneCompare has an invalid zone %q", zone.Name))
			}
			raw := hour + zone.Offset
			local := (raw + 24) % 24
			work := local >= 9 && local < 17
			name, ink := "moon", t.TextMuted
			if local >= 6 && local < 18 {
				name, ink = "sun", t.Warning
			}
			c.Key(i)
			row := kit.Card(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2.5), t.Space(3))
			kit.Stagger(c, root, row, i).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
					tile(c, name, ink)
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, zone.Name).FontSize(t.Rem(0.85)).Bold().SingleLine().MinWidth(0)
						kit.Muted(c, fmt.Sprintf("UTC%+d", zone.Offset)).FontSize(t.Rem(0.7))
					})
					if raw >= 24 || raw < 0 {
						text := o.Next
						if raw < 0 {
							text = o.Prev
						}
						kit.Chip(c, text, t.TextMuted)
					}
					kit.Numeric(c, time.Date(2000, 1, 1, local, 0, 0, 0, time.UTC).Format("15:04")).FontSize(t.Rem(1.2)).Bold().
						TextColor(t.Text.Mix(t.Success, unit(work)))
				})
				strip := ui.Box(c).FillWidth().Height(t.Space(2)).Radius(999).Clip()
				strip.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					step := r.W / 24
					for h := range 24 {
						fill := t.SurfaceHover
						if h >= 9 && h < 17 {
							fill = t.Success.Alpha(0.3)
						}
						if h == local {
							fill = t.Accent
						}
						p.Fill(ui.Rect{X: r.X + float32(h)*step + 0.5, Y: r.Y, W: step - 1, H: r.H}, fill, 0)
					}
				}))
			})
		}
	})
	return root
}
