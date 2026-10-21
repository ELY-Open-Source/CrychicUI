package typography

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// LoadingDots cycles three dots beside a label.
func LoadingDots(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(0.5))
	phase := root.Loop("dots", 900*time.Millisecond, ui.Linear)
	active := int(phase*3) % 3
	return root.Children(func() {
		ui.Text(c, label)
		for i := range 3 {
			dot := ui.Text(c, ".")
			target := float32(0.25)
			if i == active {
				target = 1
			}
			dot.Opacity(dot.AnimateTo("alpha", target, ui.Snappy))
		}
	})
}
