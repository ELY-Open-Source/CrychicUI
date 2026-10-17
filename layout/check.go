package layout

import "math"

// finite reports v is neither NaN nor infinite.
func finite(v float32) bool { return !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) }
