package typography

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Duration is a human duration with rolling digits.
func Duration(c *ui.Context, d time.Duration) ui.Element {
	if d < 0 {
		panic("typography: Duration must be non-negative")
	}
	var text string
	switch {
	case d >= time.Hour:
		text = fmt.Sprintf("%dh %dm", d/time.Hour, (d%time.Hour)/time.Minute)
	case d >= time.Minute:
		text = fmt.Sprintf("%dm %ds", d/time.Minute, (d%time.Minute)/time.Second)
	default:
		text = fmt.Sprintf("%ds", d/time.Second)
	}
	return kit.Numeric(c, text)
}
