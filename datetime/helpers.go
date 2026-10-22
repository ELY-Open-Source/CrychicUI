// Package datetime holds date, time and clock components.
package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// locale is the locale MyGo's own calendar formats with.
func locale(c *ui.Context) string { return c.Region().FormatLocale(c.Locale()) }

// sameDay reports whether a and b share a calendar day.
func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// monthOf is the first day of t's month, at midnight.
func monthOf(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// onDay keeps clock's time of day on day's date.
func onDay(day, clock time.Time) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), clock.Hour(), clock.Minute(), clock.Second(),
		clock.Nanosecond(), clock.Location())
}

// addMonths moves t n months, clamping to the last day.
func addMonths(t time.Time, n int) time.Time {
	first := time.Date(t.Year(), t.Month()+time.Month(n), 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	last := first.AddDate(0, 1, -1).Day()
	return first.AddDate(0, 0, min(t.Day(), last)-1)
}

// look is how one day cell shows.
type look struct{ chosen, band, dot bool }

// grid lays out month m's six weeks; style styles days.
func grid(c *ui.Context, m time.Time, cursor *time.Time, compact bool, style func(time.Time) look) (picked, hovered time.Time) {
	t := c.Theme()
	loc := locale(c)
	first := c.Region().FirstWeekday
	w, h, names := t.Space(8), t.Space(7), i18n.Weekdays(loc, i18n.Short, true)
	if compact {
		w, h, names = t.Space(6), t.Space(5.5), i18n.Weekdays(loc, i18n.Narrow, true)
	}
	g := ui.Column(c).Gap(t.Space(0.5)).Shrink(0).Focusable().Role(ui.RoleTable).
		Label(i18n.FormatDate(loc, m, i18n.DateOptions{Skeleton: "yMMMM"}))
	back, forth := ui.KeyLeft, ui.KeyRight
	if c.Direction() == ui.RightToLeft {
		back, forth = forth, back
	}
	cur := *cursor
	for _, k := range []struct {
		key ui.Key
		to  time.Time
	}{{back, cur.AddDate(0, 0, -1)}, {forth, cur.AddDate(0, 0, 1)}, {ui.KeyUp, cur.AddDate(0, 0, -7)},
		{ui.KeyDown, cur.AddDate(0, 0, 7)}, {ui.KeyPageUp, addMonths(cur, -1)}, {ui.KeyPageDown, addMonths(cur, 1)}} {
		if g.Shortcut(0, k.key) {
			*cursor = k.to
		}
	}
	if g.Shortcut(0, ui.KeyEnter) || g.Shortcut(0, ui.KeySpace) {
		picked = cur
	}
	focused := g.Focused()
	g.Children(func() {
		ui.Row(c).Children(func() {
			for i := range 7 {
				ui.Text(c, names[(int(first)+i)%7]).Width(w).TextAlign(ui.Center).FontSize(t.Rem(0.75)).
					TextColor(t.TextMuted).Role(ui.RoleColumnHeader)
			}
		})
		start := m.AddDate(0, 0, -((int(m.Weekday()) - int(first) + 7) % 7))
		today := c.Now()
		slide(c, monthIndex(m), func() { weeks(c, start, m, w, h, cursor, focused, today, style, &picked, &hovered) })
	})
	return picked, hovered
}

// weeks lays out six rows of days from start.
func weeks(c *ui.Context, start, m time.Time, w, h float32, cursor *time.Time, focused bool, today time.Time,
	style func(time.Time) look, picked, hovered *time.Time) {
	cur := *cursor
	for row := range 6 {
		ui.Row(c).Role(ui.RoleRow).Children(func() {
			for col := range 7 {
				day := start.AddDate(0, 0, row*7+col)
				if d, over := cell(c, day, m, w, h, style(day), focused && sameDay(day, cur), sameDay(day, today)); d {
					*picked, *cursor = day, day
				} else if over {
					*hovered = day
				}
			}
		})
	}
}

// cell is one day, reporting a click and a hover.
func cell(c *ui.Context, day, m time.Time, w, h float32, lk look, atCursor, today bool) (clicked, hovered bool) {
	t := c.Theme()
	loc := locale(c)
	b := ui.Box(c).Size(w, h).Center().Radius(t.Radius).Role(ui.RoleButton).
		Label(i18n.FormatDate(loc, day, i18n.DateOptions{Skeleton: "yMMMMd"})).Checked(lk.chosen)
	bg, fg := t.Accent.Alpha(0), t.Text
	if day.Month() != m.Month() {
		fg = t.TextMuted
	}
	switch {
	case lk.chosen:
		bg, fg = t.Accent, t.AccentText
	case lk.band:
		bg = t.Accent.Alpha(0.14)
	case b.Hovered() || atCursor:
		bg = t.SurfaceHover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy))
	if today && !lk.chosen {
		b.Border(1, t.Accent)
	}
	b.Children(func() {
		ui.Column(c).AlignItems(ui.Center).Children(func() {
			ui.Text(c, i18n.FormatDate(loc, day, i18n.DateOptions{Skeleton: "d"}))
			if lk.dot {
				dot := t.Accent
				if lk.chosen {
					dot = t.AccentText
				}
				ui.Box(c).Size(t.Space(1), t.Space(1)).Radius(999).Background(dot)
			}
		})
	})
	return b.Clicked(), b.Hovered()
}

// head is a title between previous and next buttons.
func head(c *ui.Context, title, prev, next string) (back, ahead ui.Element) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
		back = button.IconButton(c, prev, "chevron-left")
		ui.Text(c, title).Bold().Grow(1).TextAlign(ui.Center).ContentTransition(ui.ContentTransition{})
		ahead = button.IconButton(c, next, "chevron-right")
	})
	return back, ahead
}

// slide keys content, sliding in toward the side moved to.
func slide(c *ui.Context, key int, content func()) {
	t := c.Theme()
	box := ui.Box(c).Clip()
	last := ui.Local(box, "key", func() int { return key })
	dir := ui.Local(box, "dir", func() float32 { return 1 })
	if key != *last {
		*dir = 1
		if key < *last {
			*dir = -1
		}
		*last = key
	}
	shift := *dir * t.Space(8)
	box.Children(func() {
		ui.Column(c).Key(key).Transition(ui.ElementTransition{Enter: &ui.Motion{X: shift}, Exit: &ui.Motion{X: -shift}}).
			Children(content)
	})
}

// monthIndex numbers months so later months sort higher.
func monthIndex(t time.Time) int { return t.Year()*12 + int(t.Month()) }

func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// def sets an empty copy field to its English default.
func def(field *string, english string) {
	if *field == "" {
		*field = english
	}
}

// dayBefore reports whether a's day comes before b's.
func dayBefore(a, b time.Time) bool { return !sameDay(a, b) && a.Before(b) }

// withYearMonth moves t to y and mo, clamping the day.
func withYearMonth(t time.Time, y int, mo time.Month) time.Time {
	first := time.Date(y, mo, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	return first.AddDate(0, 0, min(t.Day(), first.AddDate(0, 1, -1).Day())-1)
}

// tile is a grid choice tinted when chosen.
func tile(c *ui.Context, label string, chosen, muted bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Height(t.Space(9)).Grow(1).Basis(0).Radius(t.Radius)).Checked(chosen)
	bg, fg := t.Accent.Alpha(0), t.Text
	switch {
	case chosen:
		bg, fg = t.Accent, t.AccentText
	case b.Hovered():
		bg = t.SurfaceHover
	case muted:
		fg = t.TextMuted
	}
	return b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy)).
		Children(func() { ui.Text(c, label) })
}

// tiles lays labels out in rows of per.
func tiles(c *ui.Context, labels []string, per int, each func(i int, label string)) {
	t := c.Theme()
	ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
		for r := 0; r < len(labels); r += per {
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				for i := r; i < min(r+per, len(labels)); i++ {
					each(i, labels[i])
				}
			})
		}
	})
}

// segment strokes a straight line at any angle.
func segment(p *ui.Painter, x0, y0, x1, y1, width float32, col ui.Color) {
	var l ui.Path
	p.StrokePath(l.MoveTo(x0, y0).LineTo(x1, y1), width, col)
}

// act is a button whose cue its component plays.
func act(c *ui.Context, label string, solid bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(3.5)).Radius(t.Radius))
	bg, fg := t.Surface, t.Text
	if solid {
		bg, fg = t.Accent, t.AccentText
	}
	if b.Hovered() {
		bg = bg.Mix(t.Text, 0.08)
	}
	return b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(fg).Children(func() { ui.Text(c, label).SingleLine() })
}

// ring draws frac of a circle around content.
func ring(c *ui.Context, frac float32, col ui.Color, content func()) ui.Element {
	t := c.Theme()
	size, w := t.Space(36), t.Space(2)
	box := ui.Box(c).Size(size, size).Shrink(0).Center()
	shown := box.AnimateTo("frac", frac, ui.Smooth)
	return box.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-w/2
		var track ui.Path
		p.StrokePath(track.Circle(cx, cy, rad), w, t.Border)
		if shown > 0 {
			p.Transform(ui.Rotation(-90), cx, cy, func() {
				var arc ui.Path
				p.StrokePath(arc.Circle(cx, cy, rad).Trim(0, shown), w, col)
			})
		}
	}).Children(content)
}

// twelve reports whether times here use AM and PM.
func twelve(c *ui.Context) bool {
	if hc := c.Region().HourCycle; hc != 0 {
		return hc == 12
	}
	return i18n.HourCycle(locale(c)) == 12
}
