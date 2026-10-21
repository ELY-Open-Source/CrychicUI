package display

import (
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Streak is a flickering flame beside a rolling day count.
func Streak(c *ui.Context, days int) ui.Element {
	if days < 0 {
		panic("display: Streak days must be non-negative")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center)
	return root.Children(func() {
		flame := ui.Icon(c, icon.Must("flame")).Size(16, 16).TextColor(t.Warning)
		phase := root.Loop("flicker", 900*time.Millisecond, ui.Linear)
		pulse := phase * 2
		if pulse > 1 {
			pulse = 2 - pulse
		}
		scale := 1 + 0.08*pulse
		flame.Scale(scale, scale)
		kit.Numeric(c, strconv.Itoa(days)).Bold().FontSize(t.Rem(1.25))
	})
}
