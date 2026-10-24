package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ImageGenerating is a sheen placeholder with a percent.
func ImageGenerating(c *ui.Context, prompt string, progress float64) ui.Element {
	if prompt == "" || !(progress >= 0 && progress <= 1) {
		panic(fmt.Sprintf("agent: ImageGenerating %q needs a prompt and progress from 0 to 1", prompt))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius*2).Clip().Background(t.Surface).Border(1, t.Border).
		Role(ui.RoleProgress).Label(prompt).Range(0, 1, progress)
	return root.Children(func() {
		canvas := ui.Box(c).FillWidth().AspectRatio(1.5).Center()
		if progress < 1 && !c.Preferences().ReduceMotion {
			ph := float64(canvas.Loop("sheen", 1600*time.Millisecond, ui.Linear))
			canvas.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				band := r.W / 3
				x := r.X - band + float32(ph)*(r.W+2*band)
				for i, a := range []float32{0.03, 0.06, 0.03} {
					p.Fill(ui.Rect{X: x + float32(i)*band/3, Y: r.Y, W: band / 3, H: r.H}, t.Text.Alpha(a), 0)
				}
			}))
		}
		canvas.Children(func() {
			ui.Column(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must("image")).Size(t.Space(8), t.Space(8)).TextColor(t.TextMuted)
				kit.Numeric(c, fmt.Sprintf("%.0f%%", progress*100)).Bold().FontSize(t.Rem(1.3))
			})
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2), t.Space(3), t.Space(3)).Background(t.Background).Children(func() {
			kit.Muted(c, prompt).FontSize(t.Rem(0.8)).Wrap()
			bar(c, "progress", progress, t.Accent)
		})
	})
}
