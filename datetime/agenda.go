package datetime

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Event is a titled span of time.
type Event struct {
	Title      string
	Start, End time.Time
}

// eventKey tells apart identical events by their order.
type eventKey struct {
	e Event
	n int
}

// AgendaOptions tune an Agenda.
type AgendaOptions struct {
	// Empty shows when there are no events.
	Empty string
}

// Agenda lists events under their days; a click picks one.
func Agenda(c *ui.Context, label string, events []Event, picked *int, o AgendaOptions) ui.Element {
	if label == "" || picked == nil || *picked < -1 || *picked >= len(events) {
		panic("datetime: Agenda needs a label and a pick of -1 or an event index")
	}
	for i, e := range events {
		if e.Title == "" || e.End.Before(e.Start) || i > 0 && e.Start.Before(events[i-1].Start) {
			panic(fmt.Sprintf("datetime: Agenda event %d needs a title, an end after its start and sorted order", i))
		}
	}
	def(&o.Empty, "Nothing scheduled")
	t := c.Theme()
	loc := locale(c)
	at := func(x time.Time) string { return i18n.FormatDate(loc, x, i18n.DateOptions{Skeleton: "jm"}) }
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Label(label)
	shown := *picked
	return kit.Arrive(root, ui.Pose{Y: 4, Fade: 1}, 0).Children(func() {
		if len(events) == 0 {
			kit.Muted(c, o.Empty)
			return
		}
		for i, e := range events {
			if i == 0 || !sameDay(e.Start, events[i-1].Start) {
				c.Key(fmt.Sprint("day", e.Start.Format(time.DateOnly)))
				ui.Text(c, i18n.FormatDate(loc, e.Start, i18n.DateOptions{Skeleton: "yMMMEd"})).FontSize(t.Rem(0.85)).
					Bold().TextColor(t.TextMuted).Transition(kit.Enter)
			}
			c.Key(eventKey{e, i - slices.Index(events, e)})
			row := kit.Press(ui.ButtonBase(c).Label(e.Title).FillWidth().Padding(t.Space(2), t.Space(3)).Radius(t.Radius)).
				Checked(shown == i).Transition(kit.Enter)
			bg := t.Accent.Alpha(0)
			switch {
			case shown == i:
				bg = t.Accent.Alpha(0.14)
			case row.Hovered():
				bg = t.SurfaceHover
			}
			row.Background(row.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(1), t.Space(5)).Radius(999).Background(t.Accent)
					kit.Numeric(c, at(e.Start)+" – "+at(e.End)).TextColor(t.TextMuted).Shrink(0)
					ui.Text(c, e.Title)
				})
			})
			if kit.Clicked(row, sound.Tap) {
				*picked = i
				root.SetChanged()
				kit.Log("Agenda", "picked", "index", i)
			}
		}
	})
}
