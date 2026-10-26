package desktop

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BatteryItemOptions mark charging and format the percent.
type BatteryItemOptions struct {
	Charging bool
	Percent  string
}

// BatteryItem is a menu bar battery with its charge.
func BatteryItem(c *ui.Context, level float32, o BatteryItemOptions) ui.Element {
	def(&o.Percent, "%d%%")
	share("BatteryItem", o.Percent, level)
	// Never 100 before full, never 0 above empty.
	p := int(math.Round(float64(level) * 100))
	if level < 1 {
		p = min(p, 99)
	}
	if level > 0 {
		p = max(p, 1)
	}
	t := c.Theme()
	ink := t.Text
	switch {
	case o.Charging:
		ink = t.Success
	case level < 0.2:
		ink = t.Danger
	}
	root := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Role(ui.RoleMeter)
	kit.Arrive(root, ui.Pose{Y: -t.Space(1), Fade: 1}, 0)
	return root.Children(func() {
		ui.Text(c, fmt.Sprintf(o.Percent, p)).FontSize(t.Rem(0.76)).Shrink(0).MaxWidth(t.Space(20)).SingleLine().ContentTransition(ui.ContentTransition{})
		cell := ui.Box(c).Size(t.Space(7), t.Space(3.5)).Shrink(0).Center()
		fill := cell.AnimateTo("fill", level, ui.Smooth)
		paint := cell.AnimateColor("ink", ink, ui.Snappy)
		cell.Draw(func(pt *ui.Painter, r ui.Rect) {
			body := ui.Rect{X: r.X, Y: r.Y, W: r.W - t.Space(0.75), H: r.H}
			pt.Fill(body, t.TextMuted.Alpha(0.5), t.Space(0.75))
			pt.Fill(ui.Rect{X: body.X + 1, Y: body.Y + 1, W: body.W - 2, H: body.H - 2}, t.Background, t.Space(0.5))
			pt.Fill(ui.Rect{X: body.X + 2, Y: body.Y + 2, W: (body.W - 4) * fill, H: body.H - 4}, paint, t.Space(0.25))
			pt.Fill(ui.Rect{X: r.X + r.W - t.Space(0.5), Y: r.Y + r.H*0.3, W: t.Space(0.5), H: r.H * 0.4}, t.TextMuted.Alpha(0.5), t.Space(0.25))
		})
		cell.Children(func() {
			if o.Charging {
				kit.Arrive(glyph(c, "zap", t.Space(2.5)).TextColor(t.Text), ui.Pose{Scale: 0.5, Fade: 1}, 0)
			}
		})
	})
}
