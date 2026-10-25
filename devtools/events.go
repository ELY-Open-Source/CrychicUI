package devtools

import (
	"fmt"
	"hash/fnv"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Event is one app event; Detail is optional.
type Event struct {
	ID, Time, Name, Detail string
}

// EventsOptions word the panel.
type EventsOptions struct{ Title, Filter, Empty string }

// Events is an event stream with a filter; cue Tick.
func Events(c *ui.Context, events []Event, filter *string, o EventsOptions) ui.Element {
	if filter == nil {
		panic("devtools: Events needs a filter")
	}
	ids := make([]string, len(events))
	for i, e := range events {
		if e.Time == "" || e.Name == "" {
			panic(fmt.Sprintf("devtools: Events %q needs a time and name", e.ID))
		}
		ids[i] = e.ID
	}
	if len(events) > 0 {
		unique("Events", ids)
	}
	def(&o.Title, "Events")
	def(&o.Filter, "Filter events")
	def(&o.Empty, "No matching events")
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	var shown []Event
	for _, e := range events {
		if lo, _ := kit.IndexFold(e.Name, *filter); *filter == "" || lo >= 0 {
			shown = append(shown, e)
		}
	}
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "radio", o.Title, func() { kit.Chip(c, fmt.Sprint(len(shown)), t.Accent) })
		frame := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
		frame.Children(func() {
			glyph(c, "filter", 14).TextColor(t.TextMuted)
			in := ui.TextInputBase(c, filter).Label(o.Filter).Placeholder(o.Filter).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0)
			if in.Changed() {
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("Events", "filtered", "runes", utf8.RuneCountInString(*filter))
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(shown) == 0 {
				empty(c, "radio", o.Empty)
			}
			for i, e := range shown {
				h := fnv.New32a()
				h.Write([]byte(e.Name))
				ink := inks[h.Sum32()%uint32(len(inks))]
				c.Key(e.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					kit.Muted(c, e.Time).Font("monospace").FontSize(t.Rem(0.7)).Shrink(0)
					ui.Box(c).Padding(t.Space(0.5), t.Space(1.5)).Radius(t.Radius * 0.6).Background(ink.Alpha(0.12)).TextColor(ink).
						Shrink(0).Children(func() {
						ui.Text(c, e.Name).Font("monospace").FontSize(t.Rem(0.72)).Bold().Shrink(0).MaxWidth(t.Space(36)).SingleLine()
					})
					if e.Detail != "" {
						kit.Muted(c, e.Detail).Font("monospace").FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
					}
				})
			}
		})
	})
	return root
}
