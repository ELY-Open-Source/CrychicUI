package content

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ReadingTimeOptions set the pace and copy; zero WPM is 230.
type ReadingTimeOptions struct {
	WPM int
	// Progress is the share read; above 0 shows Left.
	Progress   float32
	Read, Left string
}

// ReadingTime shows the minutes a text takes to read.
func ReadingTime(c *ui.Context, words int, o ReadingTimeOptions) ui.Element {
	def(&o.Read, "%d min read")
	def(&o.Left, "%d min left")
	if o.WPM == 0 {
		o.WPM = 230
	}
	p := float64(o.Progress)
	if words < 0 || words > 10_000_000 || o.WPM < 50 || o.WPM > 1000 || math.IsNaN(p) || p < 0 || p > 1 {
		panic(fmt.Sprintf("content: ReadingTime needs 0 to 10M words, 50 to 1000 wpm and progress 0 to 1, got %d, %d, %v", words, o.WPM, p))
	}
	mins := max(1, (words+o.WPM-1)/o.WPM)
	t := c.Theme()
	root := ui.Row(c).FillWidth().Role(ui.RoleStatus)
	seen := root.LayoutBounds().W
	kit.Settle(c, root, []float32{seen}, func() []float32 { return []float32{root.LayoutBounds().W} })
	// Texts share what the gauge, padding and dot leave.
	room := t.Space(40)
	if seen > 0 {
		room = max(t.Space(6), seen-t.Space(14))
	}
	if o.Progress > 0 {
		room /= 2
	}
	return root.Children(func() {
		pill := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(3), t.Space(1), t.Space(1)).Radius(999).
			Background(t.SurfaceHover).MinWidth(0)
		pill.Children(func() {
			gauge := ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().TextColor(t.Accent)
			share := gauge.AnimateTo("share", o.Progress, ui.Smooth)
			gauge.Draw(func(pt *ui.Painter, r ui.Rect) {
				cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(0.5)
				pt.StrokePath(ring(cx, cy, rad), t.Space(0.5), t.Border)
				if share > 0.001 {
					pt.StrokePath(arc(cx, cy, rad, 0, 360*share), t.Space(0.5), t.Accent)
				}
			})
			gauge.Children(func() { glyph(c, "clock", 12) })
			ui.Text(c, fmt.Sprintf(o.Read, mins)).FontSize(t.Rem(0.76)).Bold().Shrink(0).MaxWidth(room).SingleLine().
				ContentTransition(ui.ContentTransition{})
			if o.Progress > 0 {
				left := int(math.Ceil(float64(mins) * (1 - p)))
				dot(c)
				ui.Text(c, fmt.Sprintf(o.Left, left)).FontSize(t.Rem(0.76)).TextColor(t.TextMuted).Shrink(0).MaxWidth(room).SingleLine().
					ContentTransition(ui.ContentTransition{})
			}
		})
	})
}
