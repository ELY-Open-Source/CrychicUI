package desktop

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// hud shows a level in segments; it fades after changes.
func hud(c *ui.Context, label, mark string, level float32) ui.Element {
	t := c.Theme()
	root := ui.Column(c).Width(t.Space(48)).Gap(t.Space(3)).Padding(t.Space(4)).AlignItems(ui.Center).Radius(t.Radius*2.5).
		Background(t.SurfaceHover.Alpha(0.94)).Border(1, t.Border).Shadow(0, 8, 24, 0, t.Text.Alpha(0.12)).Role(ui.RoleStatus).Label(label)
	seen := ui.Local(root, "seen", func() float32 { return -1 })
	at := ui.Local(root, "at", func() time.Time { return time.Time{} })
	if level != *seen {
		*seen, *at = level, c.Now()
	}
	left := 1500*time.Millisecond - c.Now().Sub(*at)
	if left > 0 {
		c.After(left)
	}
	root.Opacity(root.AnimateTo("shown", unit(left > 0), ui.Smooth))
	n := int(math.Round(float64(level) * 16))
	return root.Children(func() {
		glyph(c, mark, 32).TextColor(t.Text).ContentTransition(ui.ContentTransition{})
		ui.Row(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
			for i := range 16 {
				seg := ui.Box(c).Grow(1).Basis(0).Height(t.Space(1.5)).Radius(999)
				seg.Background(seg.AnimateColor("fill", t.Border.Mix(t.Text, unit(i < n)), ui.Snappy))
			}
		})
	})
}
