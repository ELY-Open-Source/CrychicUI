package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// SegmentProgress is n segments filling in order.
func SegmentProgress(c *ui.Context, label string, n, done int) ui.Element {
	if label == "" || n < 1 || done < 0 || done > n {
		panic("feedback: SegmentProgress needs a label, n >= 1 and 0 <= done <= n")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleProgress).Label(label).Range(0, float64(n), float64(done))
	return root.Children(func() {
		for i := range n {
			c.Key(i)
			seg := ui.Box(c).Grow(1).Basis(0).Height(t.Space(1.5)).Radius(t.Space(0.75)).Background(t.Border).Clip()
			seg.Children(func() {
				a := ui.Smooth
				a.Delay = time.Duration(i) * 60 * time.Millisecond
				fill := ui.Box(c).Absolute().Start(0).Top(0).Bottom(0).Background(t.Accent)
				fill.WidthPercent(fill.AnimateTo("w", 100*unit(i < done), a))
			})
		}
	})
}
