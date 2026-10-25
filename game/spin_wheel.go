package game

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SpinWheelOptions word the wheel; Result takes the slice.
type SpinWheelOptions struct {
	Spin, Result string
	Rand         func(n int) int
}

// SpinWheel spins to a random slice; cue Tick per slice.
func SpinWheel(c *ui.Context, slices []string, picked *int, o SpinWheelOptions) ui.Element {
	def(&o.Spin, "Spin")
	n := len(slices)
	if picked == nil || n < 2 || n > 12 || *picked < -1 || *picked >= n {
		panic(fmt.Sprintf("game: SpinWheel %q needs 2 to 12 slices and a pick in range", o.Spin))
	}
	for _, s := range slices {
		if strings.TrimSpace(s) == "" {
			panic(fmt.Sprintf("game: SpinWheel %q has a blank slice", o.Spin))
		}
	}
	def(&o.Result, "Landed on %s")
	if o.Rand == nil {
		o.Rand = pick
	}
	t := c.Theme()
	step := 2 * math.Pi / float64(n)
	land := func(i int) float64 { return mod(-(float64(i)+0.5)*step, 2*math.Pi) }
	under := func(phi float64) int { return int(mod(-phi, 2*math.Pi)/step) % n }
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger, t.Accent.Mix(t.Danger, 0.5), t.Success.Mix(t.Accent, 0.5)}
	ink := func(i int) ui.Color {
		if i == n-1 && n%len(inks) == 1 {
			return inks[2]
		}
		return inks[i%len(inks)]
	}
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	spins := ui.Local(root, "spins", func() int { return 0 })
	angle := ui.Local(root, "angle", func() float64 {
		if *picked >= 0 {
			return land(*picked)
		}
		return 0
	})
	from := ui.Local(root, "from", func() float64 { return *angle })
	last := ui.Local(root, "last", func() int { return under(*angle) })
	seen := ui.Local(root, "seen", func() [2]int { return [2]int{*picked, n} })
	// A new pick or count turns the wheel there.
	if *seen != [2]int{*picked, n} {
		if *picked >= 0 {
			*angle = land(*picked)
		}
		*from, *last, *seen = *angle, under(*angle), [2]int{*picked, n}
	}
	spinning := false
	root.Children(func() {
		size := t.Space(60)
		c.Key(*spins)
		wheel := ui.Box(c).Size(size, size).Shrink(0)
		phi := *angle
		if *spins > 0 {
			u := float64(wheel.Keyframes("spin", []ui.Keyframe{{Value: 0}, {At: 3 * time.Second, Value: 1}}))
			phi = *from + (*angle-*from)*(1-math.Pow(1-u, 3))
			spinning = u < 1 && *from != *angle
			if at := under(phi); spinning && at != *last {
				*last = at
				sound.Play(sound.Tick)
			}
		}
		wheel.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2-t.Space(1)
			for i := range n {
				a := float64(i)*step + phi
				p.FillPath(wedge(cx, cy, rr, a, a+step), ink(i).Alpha(0.85))
			}
			p.StrokePath(ring(cx, cy, rr), t.Space(0.75), t.Background)
			p.FillPath(ring(cx, cy, t.Space(4)), t.Background)
			p.FillPath(ring(cx, cy, t.Space(2.5)), t.Text)
			var tip ui.Path
			tip.MoveTo(cx-t.Space(2.5), r.Y).LineTo(cx+t.Space(2.5), r.Y).LineTo(cx, r.Y+t.Space(5))
			p.FillPath(tip.Close(), t.Text)
		}))
		wheel.Children(func() {
			for i, s := range slices {
				mid := (float64(i)+0.5)*step + phi
				w := t.Space(18)
				x := size/2 + 0.6*size/2*float32(math.Sin(mid))
				y := size/2 - 0.6*size/2*float32(math.Cos(mid))
				ui.Text(c, s).FontSize(t.Rem(0.68)).Bold().TextColor(t.AccentText).Absolute().Left(x - w/2).Top(y - t.Space(2)).
					Width(w).TextAlign(ui.Center).SingleLine()
			}
		})
		c.Key("result")
		text := ""
		if *picked >= 0 {
			text = fmt.Sprintf(o.Result, slices[*picked])
		}
		res := ui.Text(c, text).FontSize(t.Rem(0.95)).Bold().Shrink(0)
		res.Opacity(res.AnimateTo("shown", unit(*picked >= 0 && !spinning), ui.Smooth))
		b := kit.Press(ui.ButtonBase(c).Label(o.Spin).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0).Disabled(spinning))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "rotate-ccw", 16)
				ui.Text(c, o.Spin).FontSize(t.Rem(0.85)).Bold().Shrink(0)
			})
		})
		if b.Clicked() && !spinning {
			i := o.Rand(n)
			if i < 0 || i >= n {
				panic(fmt.Sprintf("game: SpinWheel %q Rand gave %d for %d slices", o.Spin, i, n))
			}
			*from = mod(*angle, 2*math.Pi)
			to := land(i)
			for to < *from+8*math.Pi {
				to += 2 * math.Pi
			}
			*angle, *picked, *last, *seen = to, i, under(*from), [2]int{i, n}
			*spins++
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("SpinWheel", "spun", "slice", i)
		}
	})
	return root
}

// mod is a non-negative remainder.
func mod(a, b float64) float64 {
	r := math.Mod(a, b)
	if r < 0 {
		r += b
	}
	return r
}
