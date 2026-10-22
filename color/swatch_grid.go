package color

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SwatchGridOptions tune a SwatchGrid.
type SwatchGridOptions struct {
	// Columns per row; 8 when zero.
	Columns int
	// Names name the swatches; their hex when empty.
	Names []string
}

// SwatchGrid picks one of its colors; the chosen one rings.
func SwatchGrid(c *ui.Context, label string, colors []ui.Color, value *ui.Color, o SwatchGridOptions) ui.Element {
	if o.Columns == 0 {
		o.Columns = 8
	}
	if label == "" || len(colors) == 0 || value == nil || o.Columns < 1 || o.Names != nil && len(o.Names) != len(colors) {
		panic(fmt.Sprintf("color: SwatchGrid %q needs colors, a color, Columns >= 1 and a name per color", label))
	}
	t := c.Theme()
	size, gap := t.Space(8), t.Space(1.5)
	pick := slices.Index(colors, *value)
	shown, changed := pick, false
	root := ui.RadioGroup(c, func() {
		for i, col := range colors {
			name := hexOf(col)
			if o.Names != nil {
				name = o.Names[i]
			}
			c.Key(i)
			b := kit.Press(ui.RadioBase(c, &pick, i).Label(name).Tooltip(name).Size(size, size).Padding(t.Space(0.5)).Radius(t.Radius))
			ring := t.Accent.Alpha(0)
			if i == shown {
				ring = t.Accent
			} else if b.Hovered() {
				ring = t.Border
			}
			b.Border(2, b.AnimateColor("ring", ring, ui.Snappy)).Children(func() {
				swatch(c, col, size-t.Space(1)-4, size-t.Space(1)-4)
			})
			changed = kit.Changed(b, sound.Select) || changed
		}
	}).Label(label).Row().Wrap().Gap(gap).MaxWidth(float32(o.Columns)*(size+gap) - gap)
	if changed {
		*value = colors[pick]
		root.SetChanged()
		kit.Log("SwatchGrid", "picked", "index", pick)
	}
	return root
}
