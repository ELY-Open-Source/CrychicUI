// Package chat holds message bubbles, lists and composers.
package chat

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Message is one line of a conversation.
type Message struct {
	Sender, Text string
	Mine         bool
	At           time.Time
}

// sprintf writes a formatted string.
func sprintf(format string, args ...any) string { return fmt.Sprintf(format, args...) }

// def fills an empty string with its default.
func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

// finite reports whether v is a usable size.
func finite(v float32) bool { return !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) }

// size panics on a bad size; zero becomes d.
func size(component string, v *float32, d float32) {
	if *v < 0 || !finite(*v) {
		panic(fmt.Sprintf("chat: %s sizes must be finite and not negative", component))
	}
	if *v == 0 {
		*v = d
	}
}

// messages panics on one lacking sender, text or time.
func messages(component string, list []Message) {
	for _, m := range list {
		if m.Sender == "" || m.Text == "" || m.At.IsZero() {
			panic(fmt.Sprintf("chat: %s messages need a sender, text and a time", component))
		}
	}
}

// shell is a bubble surface, accent for mine.
func shell(c *ui.Context, mine bool) ui.Element {
	t := c.Theme()
	r := t.Radius * 2
	b := ui.Column(c).MaxWidth(t.Space(80)).Padding(t.Space(2), t.Space(3)).Radius(r).Gap(t.Space(1))
	if mine {
		return b.AlignSelf(ui.End).Background(t.Accent).TextColor(t.AccentText)
	}
	return b.AlignSelf(ui.Start).Background(t.Surface).Border(1, t.Border)
}

// clock writes d as m:ss.
func clock(d time.Duration) string {
	s := int64(d / time.Second)
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// when writes a time as a short local clock.
func when(c *ui.Context, at time.Time) string {
	return i18n.FormatDate(c.Region().FormatLocale(c.Locale()), at, i18n.DateOptions{Time: i18n.ShortStyle})
}

// round is an icon button; the caller plays its cue.
func round(c *ui.Context, name, glyph string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999))
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
	return b.Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
	})
}
