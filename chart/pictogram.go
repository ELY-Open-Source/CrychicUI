package chart

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PictogramOptions tune a Pictogram.
type PictogramOptions struct {
	// Per is one icon's value; zero gives 1.
	Per float64
	// Format prints a value; "%g" when empty.
	Format string
}

// Pictogram fills a row of icons per name.
func Pictogram(c *ui.Context, label string, names []string, vs []float64, glyph string, at *int, o PictogramOptions) ui.Element {
	values("Pictogram", label, names, vs)
	finite("Pictogram", label, o.Per)
	if o.Per < 0 {
		panic(fmt.Sprintf("chart: Pictogram %q per must not be negative", label))
	}
	if o.Per == 0 {
		o.Per = 1
	}
	most := 0
	for _, v := range vs {
		if v < 0 {
			panic(fmt.Sprintf("chart: Pictogram %q values must not be negative", label))
		}
		most = max(most, int(math.Ceil(v/o.Per)))
	}
	if most > 100 {
		panic(fmt.Sprintf("chart: Pictogram %q needs 100 icons a row at most; raise Per", label))
	}
	chosen("Pictogram", label, at, len(names))
	def(&o.Format, "%g")
	svg := icon.Must(glyph)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, v := range vs {
			c.Key(names[i])
			full := int(math.Round(v / o.Per))
			tip := names[i] + ": " + fmt.Sprintf(o.Format, v)
			row := kit.Press(ui.ButtonBase(c).Label(tip).Tooltip(tip).Checked(*at == i).FillWidth().Justify(ui.Start).Gap(t.Space(3)).
				AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
			hover := t.SurfaceHover.Alpha(unit(row.Hovered() || *at == i))
			row.Background(row.AnimateColor("bg", hover, ui.Snappy)).Children(func() {
				kit.Muted(c, names[i]).Width(t.Space(20)).Shrink(0).NoWrap()
				ui.Row(c).Wrap().Gap(t.Space(0.5)).Children(func() {
					for k := range most {
						c.Key(k)
						ink := t.Border
						if k < full {
							ink = tone(t, i)
						}
						if *at >= 0 && *at != i && k < full {
							ink = ink.Alpha(0.4)
						}
						glyphBox := ui.Icon(c, svg).Size(t.Space(4), t.Space(4))
						glyphBox.TextColor(glyphBox.AnimateColor("ink", ink, ui.Snappy))
						kit.Arrive(glyphBox, ui.Pose{Scale: 0.4, Fade: 1}, time.Duration(k)*20*time.Millisecond)
					}
				})
			})
			if row.Clicked() {
				pick(root, at, i, "Pictogram")
			}
		}
	})
}
