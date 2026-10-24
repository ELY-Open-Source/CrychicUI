package schedule

import (
	"fmt"
	"sort"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ConflictsOptions word the list; Clear shows when none overlap.
type ConflictsOptions struct{ Title, Format, One, Clear, With string }

// Conflicts flags overlapping events.
func Conflicts(c *ui.Context, events []Event, o ConflictsOptions) ui.Element {
	if len(events) == 0 {
		panic("schedule: Conflicts needs events")
	}
	validEvents("Conflicts", events)
	def(&o.Title, "Conflicts")
	def(&o.Format, "%d conflicts")
	def(&o.One, "1 conflict")
	def(&o.Clear, "No conflicts")
	def(&o.With, "Overlaps %s")
	t := c.Theme()
	list := append([]Event(nil), events...)
	sort.SliceStable(list, func(i, j int) bool { return list[i].Start.Before(list[j].Start) })
	count := 0
	with := map[string]string{}
	for i, a := range list {
		for j := i + 1; j < len(list); j++ {
			b := list[j]
			if a.Start.Before(b.End) && b.Start.Before(a.End) {
				count++
				if with[a.ID] == "" {
					with[a.ID] = b.Title
				}
				if with[b.ID] == "" {
					with[b.ID] = a.Title
				}
			}
		}
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "triangle-alert", o.Title, func() {
			text, ink := o.Clear, t.Success
			if count > 0 {
				text, ink = fmt.Sprintf(o.Format, count), t.Danger
				if count == 1 {
					text = o.One
				}
			}
			kit.Chip(c, text, ink)
		})
		body := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2))
		body.Children(func() {
			for i, e := range list {
				clash := with[e.ID] != ""
				ink := inks(t)[e.Tint]
				c.Key(e.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch).Padding(t.Space(2)).Radius(t.Radius).
					Border(1, t.Border.Mix(t.Danger, unit(clash))).Background(t.Danger.Alpha(0.06 * unit(clash)))
				kit.Stagger(c, body, row, i).Children(func() {
					ui.Box(c).Width(t.Space(0.75)).Shrink(0).Radius(999).Background(ink)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, e.Title).FontSize(t.Rem(0.85)).Bold().SingleLine().MinWidth(0)
						meta(c, "clock", eventText(e, e.Start.Location()))
						if clash {
							ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.Danger).Children(func() {
								glyph(c, "triangle-alert", 12)
								ui.Text(c, fmt.Sprintf(o.With, with[e.ID])).FontSize(t.Rem(0.7)).Bold().Grow(1).MinWidth(0).SingleLine()
							})
						}
					})
				})
			}
		})
	})
	return root
}
