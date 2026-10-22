package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CompareOptions tune a Compare.
type CompareOptions struct {
	Before, After string
}

// Compare splits before and after at a dragged divider.
func Compare(c *ui.Context, label string, before, after *ui.Bitmap, split *float64, o CompareOptions) ui.Element {
	pictures("Compare", label, []*ui.Bitmap{before, after})
	if split == nil || !(*split >= 0 && *split <= 1) {
		panic(fmt.Sprintf("media: Compare %q needs a split from 0 to 1", label))
	}
	def(&o.Before, "Before")
	def(&o.After, "After")
	t := c.Theme()
	w, h := before.Size()
	s := ui.SliderBase(c, split, 0, 1).Step(0.01).Label(label).FillWidth().AspectRatio(float32(w) / float32(h)).
		Radius(t.Radius).Clip()
	f := s.AnimateTo("split", float32(*split), ui.Snappy)
	_, _, held := s.Dragged()
	grip := s.AnimateTo("grip", unit(held || s.Hovered()), ui.Snappy)
	rtl := c.Direction() == ui.RightToLeft
	s.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		x := r.X + r.W*f
		cut := ui.Rect{X: x, Y: r.Y, W: r.X + r.W - x, H: r.H}
		if rtl {
			x = r.X + r.W*(1-f)
			cut = ui.Rect{X: r.X, Y: r.Y, W: x - r.X, H: r.H}
		}
		p.Image(before, r, ui.Cover)
		p.Clip(cut, 0, func() { p.Image(after, r, ui.Cover) })
		line := t.Space(0.5)
		p.Fill(ui.Rect{X: x - line/2, Y: r.Y, W: line, H: r.H}, t.Background, 0)
		d := t.Space(8) * (1 + 0.15*grip)
		p.Fill(ui.Rect{X: x - d/2, Y: r.Y + r.H/2 - d/2, W: d, H: d}, t.Border, d/2)
		p.Fill(ui.Rect{X: x - d/2 + 1, Y: r.Y + r.H/2 - d/2 + 1, W: d - 2, H: d - 2}, t.Background, d/2-1)
	}))
	s.Children(func() {
		for k, name := range []string{o.Before, o.After} {
			c.Key(k)
			tag := ui.Text(c, name).SingleLine().Absolute().Top(t.Space(2)).FontSize(t.Rem(0.75)).Bold().
				Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.Background.Alpha(0.85)).TextColor(t.Text)
			if k == 0 {
				tag.Start(t.Space(2)).Opacity(tag.AnimateTo("seen", unit(*split > 0.15), ui.Snappy))
			} else {
				tag.End(t.Space(2)).Opacity(tag.AnimateTo("seen", unit(*split < 0.85), ui.Snappy))
			}
		}
	})
	released(s, "Compare", *split)
	return s
}
