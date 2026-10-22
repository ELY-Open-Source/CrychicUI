package data

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Density is how tall table rows are.
type Density uint8

// The three densities; zero is the default.
const (
	DensityDefault Density = iota
	DensityCompact
	DensityRoomy
)

// Pad is a row's padding at this density.
func (d Density) Pad(t *ui.Theme) float32 {
	return t.Space([]float32{2, 1, 3.5}[d])
}

// DensitySwitchOptions tune a DensitySwitch.
type DensitySwitchOptions struct {
	// Names label compact, default and roomy.
	Names [3]string
}

// DensitySwitch picks compact, default or roomy rows.
func DensitySwitch(c *ui.Context, label string, d *Density, o DensitySwitchOptions) ui.Element {
	if label == "" || d == nil || *d > DensityRoomy {
		panic("data: DensitySwitch needs a label and a known density")
	}
	for i, n := range [3]string{"Compact", "Default", "Roomy"} {
		def(&o.Names[i], n)
	}
	order := []Density{DensityCompact, DensityDefault, DensityRoomy}
	at := slices.Index(order, *d)
	root := choice.Segments(c, label, o.Names[:], &at)
	if root.Changed() {
		*d = order[at]
		kit.Log("DensitySwitch", "chose", "density", fmt.Sprint(at))
	}
	return root
}
