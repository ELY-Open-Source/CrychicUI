package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StandUpPerson is one stand-up report; Blockers is optional.
type StandUpPerson struct {
	Name, Yesterday, Today, Blockers string
}

// StandUpOptions word the sections and status chips.
type StandUpOptions struct{ Yesterday, Today, Blockers, None, Blocked, OnTrack string }

// StandUp shows yesterday, today, and blockers per person.
func StandUp(c *ui.Context, people []StandUpPerson, o StandUpOptions) ui.Element {
	if len(people) == 0 {
		panic("project: StandUp needs people")
	}
	names := make([]string, len(people))
	for i, person := range people {
		if person.Yesterday == "" || person.Today == "" {
			panic(fmt.Sprintf("project: StandUp report for %q needs yesterday and today", person.Name))
		}
		names[i] = person.Name
	}
	unique("StandUp", names)
	def(&o.Yesterday, "Yesterday")
	def(&o.Today, "Today")
	def(&o.Blockers, "Blockers")
	def(&o.None, "No blockers")
	def(&o.Blocked, "Blocked")
	def(&o.OnTrack, "On track")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, person := range people {
			c.Key(person.Name)
			report := kit.Card(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3))
			kit.Stagger(c, root, report, i).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
					display.Avatar(c, person.Name, nil).Size(t.Space(8), t.Space(8)).FontSize(t.Rem(0.65))
					ui.Text(c, person.Name).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
					if person.Blockers != "" {
						kit.Chip(c, o.Blocked, t.Danger)
					} else {
						kit.Chip(c, o.OnTrack, t.Success)
					}
				})
				section := func(mark, label, text string, ink ui.Color, alert bool) {
					row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)
					if alert {
						row.Background(t.Danger.Alpha(0.08))
					}
					row.Children(func() {
						glyph(c, mark, 14).TextColor(ink)
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							kit.Muted(c, label).FontSize(t.Rem(0.65)).Bold()
							ui.Text(c, text).FontSize(t.Rem(0.8)).Wrap().FillWidth()
						})
					})
				}
				section("history", o.Yesterday, person.Yesterday, t.TextMuted, false)
				section("calendar-check", o.Today, person.Today, t.Accent, false)
				if person.Blockers != "" {
					section("octagon-alert", o.Blockers, person.Blockers, t.Danger, true)
				} else {
					ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, t.Space(2)).TextColor(t.Success).Children(func() {
						glyph(c, "circle-check", 14)
						ui.Text(c, o.None).FontSize(t.Rem(0.75)).Shrink(0)
					})
				}
			})
		}
	})
	return root
}
