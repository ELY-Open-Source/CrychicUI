package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// ViolinOptions tune a Violin.
type ViolinOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
	// Tip words a group: name, median, count.
	Tip string
}

// Violin mirrors each group's density.
func Violin(c *ui.Context, label string, groups []Group, at *int, o ViolinOptions) ui.Element {
	names, all := sample("Violin", label, groups)
	chosen("Violin", label, at, len(groups))
	def(&o.Format, "%g")
	def(&o.Tip, "%s: median %s, %d values")
	t := c.Theme()
	lo, hi := mc.Extent(all...)
	if hi == lo {
		lo, hi = lo-1, hi+1
	}
	const steps = 32
	grid32 := make([]float64, steps)
	for i := range grid32 {
		grid32[i] = lo + (hi-lo)*float64(i)/(steps-1)
	}
	dens := make([][]float64, len(groups))
	meds := make([]float64, len(groups))
	top := 0.0
	for i, g := range groups {
		dens[i] = density(g.Values, grid32)
		meds[i] = quartiles(g.Values)[2]
		for _, d := range dens[i] {
			top = max(top, d)
		}
	}
	for _, ds := range dens {
		for k := range ds {
			ds[k] /= top
		}
	}
	x := mc.NewBand(names...).Padding(0.2, 0.1)
	y := mc.NewLinear(lo, hi).Nice(4)
	ticks := y.Ticks(4)
	xs := make([]float32, len(groups))
	for i := range groups {
		xs[i] = x.At(i)
	}
	describe := func(i int) string {
		return fmt.Sprintf(o.Tip, names[i], fmt.Sprintf(o.Format, meds[i]), len(groups[i].Values))
	}
	root := ui.Column(c).FillWidth()
	hover := ui.Local(root, "hover", func() int { return -1 })
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := grow(area)
			shown := make([][]float64, len(dens))
			for i, ds := range dens {
				shown[i] = springs(area, fmt.Sprint("d", i), ds)
			}
			mids := springs(area, "m", meds)
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for i, ds := range shown {
					mid := mc.X(r, x.At(i))
					half := (mc.X(r, x.Band(i).Hi) - mc.X(r, x.Band(i).Lo)) / 2 * g
					var path ui.Path
					for k, d := range ds {
						px, py := mid+half*float32(d), mc.Y(r, y.At(grid32[k]))
						if k == 0 {
							path.MoveTo(px, py)
						} else {
							path.LineTo(px, py)
						}
					}
					for k := len(ds) - 1; k >= 0; k-- {
						path.LineTo(mid-half*float32(ds[k]), mc.Y(r, y.At(grid32[k])))
					}
					path.Close()
					ink := tone(t, i)
					if *at >= 0 && *at != i {
						ink = ink.Alpha(0.4)
					}
					p.FillPath(&path, ink.Alpha(0.7))
					med := mc.Y(r, y.At(mids[i]))
					p.Line(mid-half/3, med, mid+half/3, med, 2, t.Background)
				}
			}))
			follow(c, root, area, at, hover, xs, 1, describe, "Violin")
		})
	})
}
