package display

import (
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"
)

// Counter is a big number counting toward its target.
func Counter(c *ui.Context, target int) ui.Element {
	if target < 0 {
		panic("display: Counter target must be non-negative")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	current := ui.Local(root, "current", func() int { return 0 })
	targetValue := ui.Local(root, "target", func() int { return target })
	started := ui.Local(root, "started", func() bool { return false })
	if !*started {
		*started = true
		root.Jump("count", 0)
	}
	p := root.AnimateTo("count", 1, ui.Smooth)
	shown := *current + int(math.Round(float64(*targetValue-*current)*float64(p)))
	if *targetValue != target {
		*current = shown
		*targetValue = target
		root.Jump("count", 0)
		p = root.AnimateTo("count", 1, ui.Smooth)
	}
	shown = *current + int(math.Round(float64(*targetValue-*current)*float64(p)))
	if p == 1 {
		shown = *targetValue
		*current = *targetValue
	}
	return root.Children(func() {
		ui.Text(c, strconv.Itoa(shown)).Bold().FontSize(t.Rem(2)).FontFeatures("tnum")
	})
}
