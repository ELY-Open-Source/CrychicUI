package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// WaffleOptions tune a Waffle.
type WaffleOptions = PieOptions

// Waffle fills a 10 by 10 grid by share.
func Waffle(c *ui.Context, label string, names []string, vs []float64, at *int, o WaffleOptions) ui.Element {
	values("Waffle", label, names, vs)
	sum := 0.0
	for _, v := range vs {
		if v < 0 {
			panic(fmt.Sprintf("chart: Waffle %q values must not be negative", label))
		}
		sum += v
	}
	if sum == 0 {
		panic(fmt.Sprintf("chart: Waffle %q needs a value above zero", label))
	}
	chosen("Waffle", label, at, len(names))
	def(&o.Format, "%g")
	t := c.Theme()
	size := side("Waffle", label, o.Size, t.Space(56))
	gap := t.Space(0.75)
	cell := (size - 9*gap) / 10
	want := shares(vs, 100)
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		var owner []int
		for i, k := range shares(springs(root, "v", vs), 100) {
			for range k {
				owner = append(owner, i)
			}
		}
		filled := func(k int, g float32) bool { return float32(k) < g*100 }
		wheel{label: label, comp: "Waffle", n: len(names), size: size,
			hit: func(x, y, g float32) int {
				col, row := int((x+size/2)/(cell+gap)), int((y+size/2)/(cell+gap))
				if x+size/2 < 0 || y+size/2 < 0 || col > 9 || row > 9 || !filled(row*10+col, g) {
					return -1
				}
				return owner[row*10+col]
			},
			describe: func(i int) string {
				return fmt.Sprintf("%s: "+o.Format+" (%d%%)", names[i], vs[i], want[i])
			},
			draw: func(p *ui.Painter, r ui.Rect, hot int, g float32) {
				for k, i := range owner {
					box := ui.Rect{X: r.X + float32(k%10)*(cell+gap), Y: r.Y + float32(k/10)*(cell+gap), W: cell, H: cell}
					ink := t.Border.Alpha(0.5)
					if filled(k, g) {
						ink = tone(t, i)
						if *at >= 0 && *at != i {
							ink = ink.Alpha(0.35)
						}
						if i == hot {
							ink = ink.Mix(t.Text, 0.15)
						}
					}
					p.Fill(box, ink, cell/4)
				}
			},
		}.build(c, root, at, nil)
		legend(c, names)
	})
}
