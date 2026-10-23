package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WatchItem is one watched symbol and sparkline.
type WatchItem struct {
	Symbol string
	Price  float64
	Change float64
	Points []float64
}

// Watchlist shows symbols with sparklines; cue Select.
func Watchlist(c *ui.Context, items []WatchItem, currency string, selected *int) ui.Element {
	if len(items) == 0 || currency == "" || selected == nil || *selected < -1 || *selected >= len(items) {
		panic("finance: Watchlist needs items and selection")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, item := range items {
			if item.Symbol == "" || item.Price < 0 || !finite(item.Price) || !finite(item.Change) || len(item.Points) < 2 {
				panic("finance: Watchlist has an invalid item")
			}
			c.Key(i)
			row := kit.Press(ui.ButtonBase(c).Label(item.Symbol).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(2)).Radius(t.Radius).Checked(i == *selected))
			row.Background(row.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			row.Children(func() {
				ui.Text(c, item.Symbol).Bold().Grow(1)
				sparkline(c, item.Points, t.Accent)
				money(c, item.Price, currency)
				kit.Numeric(c, fmt.Sprintf("%+.2f%%", item.Change))
			})
			if row.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Watchlist", "selected", "index", i)
			}
		}
	})
	return root
}

func sparkline(c *ui.Context, points []float64, ink ui.Color) ui.Element {
	t := c.Theme()
	lo, hi := points[0], points[0]
	for _, point := range points {
		if !finite(point) {
			panic("finance: sparkline point must be finite")
		}
		lo, hi = min(lo, point), max(hi, point)
	}
	span := hi - lo
	if !finite(span) {
		panic("finance: sparkline range overflow")
	}
	box := ui.Box(c).Width(t.Space(20)).Height(t.Space(6))
	box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		var path ui.Path
		for i, point := range points {
			x := r.X + r.W*float32(i)/float32(len(points)-1)
			y := r.Y + r.H*0.5
			if hi > lo {
				y = r.Y + r.H*(1-float32((point-lo)/span))
			}
			if i == 0 {
				path.MoveTo(x, y)
			} else {
				path.LineTo(x, y)
			}
		}
		p.StrokePath(&path, 1.5, ink)
	}))
	return box
}
