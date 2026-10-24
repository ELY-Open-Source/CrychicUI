package finance

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MarketStatusOptions word the card; Next takes the formatted time.
type MarketStatusOptions struct{ Open, Closed, Next, DateFormat, Opens, Closes string }

// MarketStatus shows open or closed with a countdown.
func MarketStatus(c *ui.Context, open bool, next time.Time, o MarketStatusOptions) ui.Element {
	if next.IsZero() {
		panic("finance: MarketStatus needs a next transition")
	}
	def(&o.Open, "Market open")
	def(&o.Closed, "Market closed")
	def(&o.Next, "Next %s")
	def(&o.DateFormat, "Jan 02 15:04")
	def(&o.Opens, "Opens in")
	def(&o.Closes, "Closes in")
	t := c.Theme()
	now := c.Now()
	// Whole seconds from Unix times: Sub saturates past 292 years.
	secs := next.Unix() - now.Unix()
	if next.Nanosecond() > now.Nanosecond() {
		secs++
	}
	if secs > 0 {
		frac := time.Duration((next.Nanosecond() - now.Nanosecond() + 1e9) % 1e9)
		if frac == 0 {
			frac = time.Second
		}
		c.After(frac + time.Millisecond)
	}
	text, glyph, ink, until := o.Closed, "moon", t.TextMuted, o.Opens
	if open {
		text, glyph, ink, until = o.Open, "store", t.Success, o.Closes
	}
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, glyph, ink)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					dot := ui.Box(c).Size(t.Space(4), t.Space(4)).Shrink(0)
					core := dot.AnimateColor("ink", ink, ui.Snappy)
					beat := float32(0)
					if open {
						beat = float32(0.5 + 0.5*math.Sin(2*math.Pi*float64(dot.Loop("pulse", 1600*time.Millisecond, ui.Linear))))
					}
					dot.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
						p.Fill(r, core.Alpha(0.25*beat), 999)
						d := r.W * 0.55
						p.Fill(ui.Rect{X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d}, core, 999)
					}))
					ui.Text(c, text).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
				})
				kit.Muted(c, fmt.Sprintf(o.Next, next.Format(o.DateFormat))).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
			})
		})
		if secs > 0 {
			c.Key("countdown")
			s := secs
			ui.Row(c).FillWidth().AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius).Background(t.Surface).
				Children(func() {
					kit.Muted(c, until).FontSize(t.Rem(0.8)).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%02d:%02d:%02d", s/3600, s/60%60, s%60)).Bold().FontSize(t.Rem(1.3)).Shrink(0).
						ContentTransition(ui.ContentTransition{Numeric: true})
				})
		}
	})
	return root
}
