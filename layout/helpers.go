package layout

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// finite reports v is neither NaN nor infinite.
func finite(v float32) bool { return !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) }

// unit is 1 for true, 0 for false.
func unit(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

// grow eases 0 to 1 after first build.
func grow(e ui.Element) float32 {
	return e.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
}

var collapseSpring = ui.Snappy

// collapseMotion springs a box's height in and out.
var collapseMotion = ui.ElementTransition{
	Animation: &collapseSpring,
	Enter:     &ui.Motion{Collapse: true, Opacity: 1},
	Exit:      &ui.Motion{Collapse: true, Opacity: 1},
}
