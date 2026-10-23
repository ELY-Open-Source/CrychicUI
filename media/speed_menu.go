package media

import (
	"fmt"
	"math"
	"slices"

	"github.com/egoist/mygo/ui"
)

// SpeedMenuOptions tune a SpeedMenu.
type SpeedMenuOptions struct {
	// Format writes a speed; "%g×" when empty.
	Format string
}

// SpeedMenu picks a playback speed from a list.
func SpeedMenu(c *ui.Context, label string, speeds []float64, speed *float64, open *bool, o SpeedMenuOptions) ui.Element {
	if label == "" || len(speeds) == 0 || speed == nil || open == nil {
		panic(fmt.Sprintf("media: SpeedMenu %q needs a label, speeds and state", label))
	}
	def(&o.Format, "%g×")
	names := make([]string, len(speeds))
	for i, s := range speeds {
		if !(s > 0) || math.IsInf(s, 0) || slices.Index(speeds, s) != i {
			panic(fmt.Sprintf("media: SpeedMenu %q speeds must be distinct, finite and above 0", label))
		}
		names[i] = fmt.Sprintf(o.Format, s)
	}
	at := slices.Index(speeds, *speed)
	if at < 0 {
		panic(fmt.Sprintf("media: SpeedMenu %q speed must be one of its speeds", label))
	}
	b := opener(c, label, names[at], *open)
	menu(c, b, "SpeedMenu", label, names, &at, open)
	*speed = speeds[at]
	return b
}
