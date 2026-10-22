package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ScrollProgress is a bar that fills as state scrolls down.
func ScrollProgress(c *ui.Context, label string, state *ui.ScrollState) ui.Element {
	if label == "" || state == nil {
		panic("motion: ScrollProgress needs a label and a scroll state")
	}
	var f float32
	if state.MaxY > 0 {
		f = min(max(state.Y/state.MaxY, 0), 1)
	}
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(t.Space(1)).Background(t.Border).Role(ui.RoleProgress).Label(label).
		Range(0, 1, float64(f))
	kit.Settle(c, root, []float32{state.Y, state.MaxY}, func() []float32 { return []float32{state.Y, state.MaxY} })
	return root.Children(func() {
		bar := ui.Box(c).Absolute().Top(0).Bottom(0).Start(0).Background(t.Accent)
		bar.WidthPercent(bar.AnimateTo("f", 100*f, ui.Snappy))
	})
}
