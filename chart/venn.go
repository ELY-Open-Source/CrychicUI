package chart

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// VennOptions tune a Venn.
type VennOptions struct {
	// Size is the chart's width; zero gives the default.
	Size float32
	// Format prints a count; "%g" when empty.
	Format string
	// Join joins set names; " & " default.
	Join string
}

// Venn overlaps 2 or 3 sets with region counts.
func Venn(c *ui.Context, label string, sets []string, counts []float64, at *int, o VennOptions) ui.Element {
	named("Venn", label, sets)
	masks := map[int][]int{2: {1, 2, 3}, 3: {1, 2, 4, 3, 5, 6, 7}}[len(sets)]
	if masks == nil || len(counts) != len(masks) {
		panic(fmt.Sprintf("chart: Venn %q needs 2 or 3 sets; counts run A, B, AB or A, B, C, AB, AC, BC, ABC", label))
	}
	finite("Venn", label, counts...)
	for _, v := range counts {
		if v < 0 {
			panic(fmt.Sprintf("chart: Venn %q counts must not be negative", label))
		}
	}
	chosen("Venn", label, at, len(masks))
	def(&o.Format, "%g")
	def(&o.Join, " & ")
	t := c.Theme()
	size := side("Venn", label, o.Size, t.Space(56))
	r := size * 0.27
	cx, cy := make([]float32, len(sets)), make([]float32, len(sets))
	for k := range sets {
		a := 2 * math.Pi * float64(k) / float64(len(sets))
		if len(sets) == 2 {
			a += math.Pi / 2
		}
		cx[k], cy[k] = mc.Polar(0, 0, r*0.6, a)
	}
	region := func(m int) string {
		var in []string
		for k, s := range sets {
			if m&(1<<k) != 0 {
				in = append(in, s)
			}
		}
		return strings.Join(in, o.Join)
	}
	spot := func(m int) (float32, float32) {
		var ix, iy, ox, oy, ni, no float32
		for k := range sets {
			if m&(1<<k) != 0 {
				ix, iy, ni = ix+cx[k], iy+cy[k], ni+1
			} else {
				ox, oy, no = ox+cx[k], oy+cy[k], no+1
			}
		}
		ix, iy = ix/ni, iy/ni
		if no > 0 {
			ix, iy = ix+0.5*(ix-ox/no), iy+0.5*(iy-oy/no)
		}
		return ix, iy
	}
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		wheel{label: label, comp: "Venn", n: len(masks), size: size,
			hit: func(x, y, g float32) int {
				m := 0
				for k := range sets {
					if radius(x-cx[k], y-cy[k]) <= r*g {
						m |= 1 << k
					}
				}
				for i, want := range masks {
					if want == m {
						return i
					}
				}
				return -1
			},
			describe: func(i int) string { return region(masks[i]) + ": " + fmt.Sprintf(o.Format, counts[i]) },
			draw: func(p *ui.Painter, b ui.Rect, hot int, g float32) {
				ox, oy := b.X+b.W/2, b.Y+b.H/2
				for k := range sets {
					ink := tone(t, k)
					if *at >= 0 && masks[*at]&(1<<k) == 0 {
						ink = ink.Alpha(0.4)
					}
					circle := new(ui.Path).Circle(ox+cx[k], oy+cy[k], r*g)
					p.FillPath(circle, ink.Alpha(0.3))
					width := float32(1.5)
					if hot >= 0 && masks[hot]&(1<<k) != 0 {
						width = 3
					}
					p.StrokePath(circle, width, ink)
				}
			},
		}.build(c, root, at, func() {
			g := grow(root)
			for i, m := range masks {
				x, y := spot(m)
				c.Key(i)
				ui.Row(c).Absolute().Left(size/2+x).Top(size/2+y).Size(0, 0).AlignItems(ui.Center).Children(func() {
					ui.Column(c).Width(0).AlignItems(ui.Center).Children(func() {
						kit.Numeric(c, fmt.Sprintf(o.Format, counts[i])).NoWrap().Bold().FontSize(t.Rem(0.85)).Opacity(g).
							ContentTransition(ui.ContentTransition{Numeric: true})
					})
				})
			}
		})
		legend(c, sets)

	})
}
