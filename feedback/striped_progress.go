package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// StripedProgress is a bar whose fill carries moving stripes.
func StripedProgress(c *ui.Context, label string, value float64) ui.Element {
	if label == "" {
		panic("feedback: StripedProgress needs a label")
	}
	v := fraction("StripedProgress", value)
	t := c.Theme()
	h := t.Space(2.5)
	root := ui.Box(c).FillWidth().Height(h).Radius(h/2).Background(t.Border).Clip().
		Role(ui.RoleProgress).Label(label).Range(0, 1, value)
	return root.Children(func() {
		fill := ui.Box(c).Absolute().Start(0).Top(0).Bottom(0).Radius(h / 2).Background(t.Accent).Clip()
		fill.WidthPercent(fill.AnimateTo("w", v*100, ui.Smooth)).Children(func() {
			period := h * 2 * 1.4142
			band := ui.Box(c).Absolute().Top(0).Bottom(0).Left(-period).Right(0).Stripes(t.AccentText.Alpha(0.22), h, h, 45)
			band.Translate(band.Loop("run", 700*time.Millisecond, ui.Linear)*period, 0)
		})
	})
}
