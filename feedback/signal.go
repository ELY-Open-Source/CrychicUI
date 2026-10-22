package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Signal is four rising bars, strength 0 to 4.
func Signal(c *ui.Context, label string, strength int) ui.Element {
	if label == "" || strength < 0 || strength > 4 {
		panic("feedback: Signal needs a label and a strength from 0 to 4")
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(0.5)).Height(t.Space(4)).AlignItems(ui.End).AlignSelf(ui.Start).
		Role(ui.RoleMeter).Label(label).Range(0, 4, float64(strength))
	return root.Children(func() {
		for i := range 4 {
			c.Key(i)
			a := ui.Snappy
			a.Delay = time.Duration(i) * 50 * time.Millisecond
			bar := ui.Box(c).Width(t.Space(1)).Radius(t.Space(0.25)).Background(t.Text.Alpha(0.15))
			bar.Height(t.Space(1) * float32(i+1)).Children(func() {
				lit := ui.Box(c).Absolute().Start(0).End(0).Bottom(0).Radius(t.Space(0.25)).Background(t.Text)
				lit.HeightPercent(lit.AnimateTo("lit", 100*unit(i < strength), a) * lit.Keyframes("rise", grow, ui.KeyframeOptions{Delay: a.Delay}))
			})
		}
	})
}
