package finance

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WatchItem is one watched symbol and its recent prices.
type WatchItem struct {
	Symbol, Name  string
	Price, Change float64
	Points        []float64
}

// WatchlistOptions word the card.
type WatchlistOptions struct{ Title string }

// Watchlist lists symbols with sparklines and change; cue Select.
func Watchlist(c *ui.Context, items []WatchItem, currency string, selected *int, o WatchlistOptions) ui.Element {
	if len(items) == 0 || currency == "" || selected == nil || *selected < -1 || *selected >= len(items) {
		panic("finance: Watchlist needs items, a currency and a selection or -1")
	}
	for i, item := range items {
		if item.Symbol == "" || item.Price < 0 || !finite(item.Price) || !finite(item.Change) || len(item.Points) < 2 ||
			slices.IndexFunc(items, func(w WatchItem) bool { return w.Symbol == item.Symbol }) != i {
			panic(fmt.Sprintf("finance: Watchlist item %q needs a unique symbol, finite values and two points", item.Symbol))
		}
	}
	def(&o.Title, "Watchlist")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "activity", o.Title, func() { kit.Chip(c, fmt.Sprint(len(items)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5)).FocusGroup(ui.Vertical)
		list.Children(func() {
			for i, item := range items {
				c.Key(item.Symbol)
				on := i == *selected
				row := kit.Press(ui.ButtonBase(c).Label(item.Symbol).Row().FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).
					Padding(t.Space(2)).Radius(t.Radius).Checked(on))
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					symbolTile(c, item.Symbol)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, item.Symbol).Bold().FontSize(t.Rem(0.9)).FillWidth().SingleLine()
						if item.Name != "" {
							kit.Muted(c, item.Name).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
						}
					})
					sparkline(c, item.Points, changeInk(t, item.Change))
					ui.Column(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
						money(c, item.Price, currency).Bold().FontSize(t.Rem(0.85))
						changePill(c, item.Change)
					})
				})
				if row.Clicked() && !on {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Watchlist", "selected", "index", i)
				}
			}
		})
	})
	return root
}

// sparkline draws points as a line over a faint area.
func sparkline(c *ui.Context, points []float64, ink ui.Color) ui.Element {
	lo, hi := points[0], points[0]
	for _, point := range points {
		if !finite(point) {
			panic("finance: sparkline points must be finite")
		}
		lo, hi = min(lo, point), max(hi, point)
	}
	span := hi - lo
	if !finite(span) {
		panic("finance: sparkline range overflows")
	}
	t := c.Theme()
	box := ui.Box(c).Size(t.Space(18), t.Space(8)).Shrink(0)
	draw := box.Keyframes("crychic.draw", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	return box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		var line, area ui.Path
		pad := t.Space(0.5)
		area.MoveTo(r.X, r.Y+r.H)
		for i, point := range points {
			x := r.X + r.W*float32(i)/float32(len(points)-1)
			y := r.Y + r.H/2
			if hi > lo {
				y = r.Y + pad + (r.H-2*pad)*(1-float32((point-lo)/span))
			}
			if i == 0 {
				line.MoveTo(x, y)
			} else {
				line.LineTo(x, y)
			}
			area.LineTo(x, y)
		}
		area.LineTo(r.X+r.W, r.Y+r.H).Close()
		p.FillPath(&area, ink.Alpha(0.12*draw))
		p.StrokePath(line.Trim(0, draw), t.Space(0.4), ink)
	}))
}
