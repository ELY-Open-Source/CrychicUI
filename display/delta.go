package display

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Delta is a signed change colored up or down.
func Delta(c *ui.Context, value float64) ui.Element {
	if !finite(value) {
		panic("display: Delta value must be finite")
	}
	t := c.Theme()
	color := t.TextMuted
	if value > 0 {
		color = t.Success
	} else if value < 0 {
		color = t.Danger
	}
	return kit.Numeric(c, fmt.Sprintf("%+.1f", value)).TextColor(color)
}
