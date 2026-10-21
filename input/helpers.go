package input

import (
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// finite reports whether v is finite.
func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// asciiDigits keeps only ASCII digits.
func asciiDigits(s string) string {
	var out strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			out.WriteRune(r)
		}
	}
	return out.String()
}

// changed propagates an input's Changed to its wrapper.
func changed(root, input ui.Element) bool {
	if !input.Changed() {
		return false
	}
	root.SetChanged()
	sound.Play(sound.Tick)
	return true
}
