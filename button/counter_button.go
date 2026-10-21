package button

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CounterButton increments a rolling count badge.
func CounterButton(c *ui.Context, label string, count *int) ui.Element {
	if count == nil || *count < 0 {
		panic("button: CounterButton needs a non-negative count")
	}
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	root.Children(func() {
		ui.Text(c, label).SingleLine()
		kit.Numeric(c, strconv.Itoa(*count)).Padding(t.Space(0.5), t.Space(1)).Radius(999).Background(t.Accent.Alpha(0.12))
	})
	if root.Clicked() {
		(*count)++
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("CounterButton", "incremented", "count", *count)
	}
	return root
}
