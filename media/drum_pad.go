package media

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DrumPadOptions tune a DrumPad.
type DrumPadOptions struct {
	// Names label the pads; the piece names when empty.
	Names []string
	// Columns per row; zero gives 3.
	Columns int
}

// DrumPad pads play drums and flash; hit is the last.
func DrumPad(c *ui.Context, label string, pads []sound.Piece, hit *int, o DrumPadOptions) ui.Element {
	if label == "" || len(pads) == 0 || len(o.Names) != 0 && len(o.Names) != len(pads) || o.Columns < 0 {
		panic(fmt.Sprintf("media: DrumPad %q needs a label, pads, a name each or none and columns not negative", label))
	}
	for _, p := range pads {
		if int(p) >= len(sound.Pieces()) {
			panic(fmt.Sprintf("media: DrumPad %q has unknown %v", label, p))
		}
	}
	chosen("DrumPad", label, hit, -1, len(pads))
	if len(o.Names) == 0 {
		for _, p := range pads {
			s := p.String()
			o.Names = append(o.Names, strings.ToUpper(s[:1])+s[1:])
		}
	}
	if o.Columns == 0 {
		o.Columns = 3
	}
	t := c.Theme()
	tones := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Label(label)
	padEls := make([]ui.Element, len(pads))
	root.Children(func() {
		for r := 0; r < len(pads); r += o.Columns {
			c.Key(r)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				for i := r; i < r+o.Columns; i++ {
					c.Key(i)
					if i >= len(pads) {
						ui.Box(c).Grow(1).Basis(0)
						continue
					}
					padEls[i] = kit.Press(ui.ButtonBase(c).Label(o.Names[i]).Grow(1).Basis(0).AspectRatio(1.4).Center().
						Radius(t.Radius * 1.5))
				}
			})
		}
	})
	for i, b := range padEls {
		run := ui.Local(b, "run", func() int { return 0 })
		if strike(b) {
			*run++
			*hit = i
			root.SetChanged()
			sound.PlayDrum(pads[i])
			kit.Log("DrumPad", "hit", "pad", i)
		}
		tint := tones[i%len(tones)]
		lit := flash(b, *run)
		b.Background(tint.Alpha(0.14+0.6*lit)).Border(1, tint.Alpha(0.4+0.4*lit)).
			Children(func() { ui.Text(c, o.Names[i]).SingleLine().Bold() })
	}
	return root
}
