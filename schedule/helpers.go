// Package schedule holds calendar components.
package schedule

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Event is one calendar event; Tint picks a color, 0-3.
type Event struct {
	ID, Title, Place string
	Start, End       time.Time
	Tint             int
}

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func validEvent(component string, event Event) {
	if event.ID == "" || event.Title == "" || !event.End.After(event.Start) || event.Tint < 0 || event.Tint > 3 {
		panic(fmt.Sprintf("schedule: %s has an invalid event %q", component, event.ID))
	}
}

func validEvents(component string, events []Event) {
	seen := map[string]bool{}
	for _, event := range events {
		validEvent(component, event)
		if seen[event.ID] {
			panic(fmt.Sprintf("schedule: %s has a duplicate event %q", component, event.ID))
		}
		seen[event.ID] = true
	}
}

// eventText is the event's hours on loc's clock.
func eventText(event Event, loc *time.Location) string {
	return event.Start.In(loc).Format("15:04") + "–" + event.End.In(loc).Format("15:04")
}

func sameDate(a, b time.Time) bool {
	a = a.In(b.Location())
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

// inks are the four calendar colors.
func inks(t *ui.Theme) []ui.Color { return []ui.Color{t.Accent, t.Success, t.Warning, t.Danger} }

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	return kit.Arrive(kit.Card(c).FillWidth().Role(ui.RoleGroup), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// pill is a choice chip tinted by ink while on.
func pill(c *ui.Context, label string, on bool, ink ui.Color, content func()) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1), t.Space(2.5)).Radius(999).Shrink(0).Checked(on))
	b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.14), unit(on)), ui.Snappy))
	b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(on)), ui.Snappy))
	b.TextColor(b.AnimateColor("ink", t.Text.Mix(ink, unit(on)), ui.Snappy))
	return b.Children(func() { ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(content) })
}

// glyph is a small icon in the current text color.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// tile is a tinted rounded icon square.
func tile(c *ui.Context, name string, ink ui.Color) ui.Element {
	t := c.Theme()
	return ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).
		TextColor(ink).Children(func() { glyph(c, name, 16) })
}

// span clips an event to day's local midnight bounds.
func span(event Event, day time.Time) (from, to time.Time, ok bool) {
	lo := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	hi := lo.AddDate(0, 0, 1)
	from, to = event.Start, event.End
	if from.Before(lo) {
		from = lo
	}
	if to.After(hi) {
		to = hi
	}
	return from, to, from.Before(to)
}

// hours is v's clock hour on day; 24 past it.
func hours(v, day time.Time) float32 {
	if !sameDate(v, day) {
		return 24
	}
	l := v.In(day.Location())
	return float32(l.Hour()) + float32(l.Minute())/60 + float32(l.Second())/3600
}

// lanes packs events by clock hour; least is a minimum.
func lanes(events []Event, day time.Time, least float32) [][]Event {
	type placed struct {
		e        Event
		top, end float32
	}
	var list []placed
	for _, e := range events {
		if from, to, ok := span(e, day); ok {
			top := hours(from, day)
			list = append(list, placed{e, top, max(hours(to, day), top+least)})
		}
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].top < list[j].top })
	var out [][]Event
	var ends []float32
	for _, p := range list {
		free := -1
		for i, end := range ends {
			if end <= p.top {
				free = i
				break
			}
		}
		if free < 0 {
			out, ends = append(out, nil), append(ends, 0)
			free = len(out) - 1
		}
		out[free], ends[free] = append(out[free], p.e), p.end
	}
	return out
}

// dayColumn lays a day's events by time; block draws each.
func dayColumn(c *ui.Context, day time.Time, events []Event, now time.Time, hourH float32, block func(e Event, top, height float32)) {
	t := c.Theme()
	col := ui.Box(c).Grow(1).Basis(0).MinWidth(0).Height(24*hourH).BorderWidth(0, 0, 0, 1).BorderColor(t.Border)
	col.Draw(func(p *ui.Painter, r ui.Rect) {
		for h := 1; h < 24; h++ {
			y := r.Y + float32(h)*hourH
			p.Line(r.X, y, r.X+r.W, y, 1, t.Border.Alpha(0.6))
		}
	})
	col.Children(func() {
		ui.Row(c).Absolute().Left(0).Right(0).Top(0).Bottom(0).Children(func() {
			for i, lane := range lanes(events, day, t.Space(5)/hourH) {
				c.Key(i)
				ui.Box(c).Grow(1).Basis(0).MinWidth(0).FillHeight().Children(func() {
					for _, e := range lane {
						from, to, _ := span(e, day)
						c.Key(e.ID)
						block(e, hours(from, day)*hourH, max((hours(to, day)-hours(from, day))*hourH, t.Space(5)))
					}
				})
			}
		})
		if sameDate(now, day) {
			c.Key("now")
			y := hours(now.In(day.Location()), day) * hourH
			line := ui.Box(c).Absolute().Left(0).Right(0).Top(y - t.Space(0.75)).Height(t.Space(1.5))
			line.Draw(func(p *ui.Painter, r ui.Rect) {
				mid := r.Y + r.H/2
				p.Fill(ui.Rect{X: r.X, Y: mid - 1, W: r.W, H: 2}, t.Danger, 0)
				p.Fill(ui.Rect{X: r.X - r.H/2, Y: r.Y, W: r.H, H: r.H}, t.Danger, r.H/2)
			})
		}
	})
}

// eventBox is a tinted block with bar, title, time.
func eventBox(c *ui.Context, e Event, loc *time.Location, b ui.Element, on bool) ui.Element {
	t := c.Theme()
	ink := inks(t)[e.Tint]
	b.Padding(t.Space(0.75), t.Space(1), t.Space(0.75), t.Space(1.5)).Radius(t.Space(1)).Clip().
		AlignItems(ui.Stretch).TextColor(ink.Mix(t.Text, 0.35))
	b.Background(b.AnimateColor("bg", ink.Alpha(0.14).Mix(ink.Alpha(0.26), unit(on || b.Hovered())), ui.Snappy))
	b.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y, W: t.Space(0.75), H: r.H}, ink, 0) })
	return b.Children(func() {
		ui.Column(c).FillWidth().MinWidth(0).Children(func() {
			ui.Text(c, e.Title).FontSize(t.Rem(0.7)).Bold().SingleLine().MinWidth(0)
			ui.Text(c, eventText(e, loc)).FontSize(t.Rem(0.6)).SingleLine().MinWidth(0).Opacity(0.8)
		})
	})
}

// instant keys a time over its whole range.
func instant(v time.Time) [2]int64 { return [2]int64{v.Unix(), int64(v.Nanosecond())} }
