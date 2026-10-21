package display

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// LevelOptions tune a Level.
type LevelOptions struct {
	// Label precedes the level number.
	Label string
}

// Level pairs a level number with a progress bar.
func Level(c *ui.Context, level int, progress float64, o LevelOptions) ui.Element {
	if o.Label == "" {
		o.Label = "Level"
	}
	if level < 0 || !finite(progress) || progress < 0 || progress > 1 {
		panic("display: Level needs a nonnegative level and progress")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, fmt.Sprintf("%s %d", o.Label, level)).Bold()
		ui.Row(c).FillWidth().Height(t.Space(2)).Radius(999).Background(t.Surface).Clip().Children(func() {
			fill := ui.Box(c).Height(t.Space(2)).Basis(0).Radius(999).Background(t.Accent)
			shown := fill.AnimateTo("fill", float32(progress), ui.Snappy)
			fill.Grow(shown)
			ui.Box(c).Basis(0).Grow(1 - shown)
		})
	})
}
