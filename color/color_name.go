package color

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NamedColor is a color and what people call it.
type NamedColor struct {
	Name  string
	Color ui.Color
}

// ColorName shows the nearest of names for value; it cross-fades.
func ColorName(c *ui.Context, value ui.Color, names []NamedColor) ui.Element {
	if len(names) == 0 {
		panic("color: ColorName needs names")
	}
	best, near := "", math.Inf(1)
	l, a, b := oklab(value)
	for i, n := range names {
		if n.Name == "" {
			panic(fmt.Sprintf("color: ColorName entry %d has no name", i))
		}
		nl, na, nb := oklab(n.Color)
		if d := math.Hypot(math.Hypot(l-nl, a-na), b-nb); d < near {
			best, near = n.Name, d
		}
	}
	t := c.Theme()
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		swatch(c, value, t.Space(8), t.Space(8))
		ui.Column(c).Children(func() {
			ui.Text(c, best).Bold().ContentTransition(ui.ContentTransition{})
			kit.Muted(c, hexOf(value)).ContentTransition(ui.ContentTransition{})
		})
	})
}
