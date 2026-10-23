package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ImageGenerating is a sheen placeholder with a percent.
func ImageGenerating(c *ui.Context, label string, progress float64) ui.Element {
	if label == "" || !(progress >= 0 && progress <= 1) {
		panic(fmt.Sprintf("agent: ImageGenerating %q needs a label and progress from 0 to 1", label))
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Width(t.Space(90)).AspectRatio(1.5).Radius(t.Radius*2).Clip().
		Background(t.Surface).Center().Role(ui.RoleProgress).Label(label).Range(0, 1, progress)
	if progress < 1 && !c.Preferences().ReduceMotion {
		ph := float64(root.Loop("sheen", 1600*time.Millisecond, ui.Linear))
		root.Draw(kit.Inside(t.Radius*2, func(p *ui.Painter, r ui.Rect) {
			band := r.W / 3
			x := r.X - band + float32(ph)*(r.W+2*band)
			p.Fill(ui.Rect{X: x, Y: r.Y, W: band, H: r.H}, t.Text.Alpha(0.06), 0)
		}))
	}
	return root.Children(func() {
		kit.Numeric(c, fmt.Sprintf("%.0f%%", progress*100)).Bold().FontSize(t.Rem(1.2)).TextColor(t.TextMuted)
	})
}
