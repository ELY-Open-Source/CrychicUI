package media

import (
	"cmp"
	"fmt"
	"math"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CoverFlowOptions tune a CoverFlow.
type CoverFlowOptions struct {
	// Size of a cover; zero gives the default.
	Size float32
	// Item names a cover by number; "Cover %d" by default.
	Item string
}

// CoverFlow turns covers in 3D around the chosen one.
func CoverFlow(c *ui.Context, label string, covers []*ui.Bitmap, at *int, o CoverFlowOptions) ui.Element {
	pictures("CoverFlow", label, covers)
	chosen("CoverFlow", label, at, 0, len(covers))
	def(&o.Item, "Cover %d")
	t := c.Theme()
	size("CoverFlow", label, &o.Size, t.Space(36))
	sz, n := o.Size, len(covers)
	root := ui.Box(c).FillWidth().Height(sz * 1.3).Clip().Focusable().Label(label)
	b := root.LayoutBounds()
	kit.Settle(c, root, []float32{b.W}, func() []float32 { return []float32{root.LayoutBounds().W} })
	pos := root.AnimateTo("pos", float32(*at), ui.Snappy)
	rtl := c.Direction() == ui.RightToLeft
	place := func(d float32) (x, turn, scale float32) {
		a, sign := min(abs(d), 1), float32(1)
		if d < 0 {
			sign = -1
		}
		x = sign * (a*sz*0.62 + max(abs(d)-1, 0)*sz*0.3)
		turn, scale = -sign*a*55, 1-0.18*a
		if rtl {
			x, turn = -x, -turn
		}
		return x, turn, scale
	}
	fade := func(d float32) float32 { return max(0, min(1, 3-abs(d))) }
	set := func(to int) {
		to = max(0, min(n-1, to))
		if to == *at {
			return
		}
		*at = to
		root.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("CoverFlow", "shown", "index", to)
	}
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(i, j int) int { return cmp.Compare(abs(float32(j)-pos), abs(float32(i)-pos)) })
	_, _, prev, next := ends(c)
	switch {
	case root.Shortcut(0, ui.KeyEnter), root.Shortcut(0, ui.KeySpace):
		// Keys never pick under the pointer.
	case root.Shortcut(0, prev):
		set(*at - 1)
	case root.Shortcut(0, next):
		set(*at + 1)
	case root.Shortcut(0, ui.KeyHome):
		set(0)
	case root.Shortcut(0, ui.KeyEnd):
		set(n - 1)
	case root.Clicked():
		px, _, _ := root.PointerPosition()
		for k := len(order) - 1; k >= 0; k-- {
			d := float32(order[k]) - pos
			x, turn, s := place(d)
			if fade(d) > 0 && abs(px-(b.W/2+x)) <= sz*s*float32(math.Cos(float64(turn)*math.Pi/180))/2 {
				set(order[k])
				break
			}
		}
	}
	return root.Children(func() {
		for _, i := range order {
			c.Key(i)
			d := float32(i) - pos
			x, turn, s := place(d)
			card := ui.Box(c).Absolute().Left(b.W/2+x-sz/2).Top(sz*0.15).Size(sz, sz).Radius(t.Radius).Clip().Background(t.Surface).
				Perspective(sz*3).RotateY(turn).Scale(s, s).Opacity(fade(d)).
				Role(ui.RoleImage).Label(fmt.Sprintf(o.Item, i+1))
			card.Children(func() { ui.Image(c, covers[i]).Fill().Fit(ui.Cover) })
		}
	})
}
