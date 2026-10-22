package navigation

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FloatingNavOptions tune a FloatingNav.
type FloatingNavOptions struct {
	// After is how far down the pill appears.
	After float32
}

// FloatingNav is a pill nav sliding down past After.
func FloatingNav(c *ui.Context, label string, links []string, active *int, st *ui.ScrollState, o FloatingNavOptions) ui.Element {
	tabs("FloatingNav", label, links, active)
	if st == nil || o.After < 0 || math.IsNaN(float64(o.After)) || math.IsInf(float64(o.After), 0) {
		panic(fmt.Sprintf("navigation: FloatingNav %q needs a scroll and a finite After", label))
	}
	t := c.Theme()
	if o.After == 0 {
		o.After = t.Space(40)
	}
	root := ui.Row(c).FillWidth().Justify(ui.Center)
	shown, cur := st.Y > o.After, *active
	root.Children(func() {
		if !shown {
			return
		}
		c.Key("pill")
		pill := ui.Row(c).Gap(t.Space(0.5)).Padding(t.Space(1)).Radius(999).Background(t.Background).Border(1, t.Border).
			Shadow(0, 4, 12, 0, kit.Shade(0.12)).FocusGroup(ui.Horizontal).Label(label).
			Transition(ui.ElementTransition{Enter: &ui.Motion{Y: -t.Space(4)}, Exit: &ui.Motion{Y: -t.Space(4)}})
		id := pill.ID()
		pill.Children(func() {
			for i, l := range links {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(l).Radius(999).Checked(i == cur))
				b.Children(func() {
					stack(c, i == cur, func() ui.Element {
						return ui.Box(c).Radius(999).Background(t.Accent.Alpha(0.12)).Transition(glide(id))
					}, func() ui.Element {
						return ui.Box(c).Padding(t.Space(1.5), t.Space(3)).Children(func() { ui.Text(c, l).SingleLine() })
					})
				})
				if kit.Clicked(b, sound.Tap) && i != *active {
					*active = i
					root.SetChanged()
					kit.Log("FloatingNav", "picked", "index", i)
				}
			}
		})
	})
	kit.Settle(c, root, []float32{unit(shown)}, func() []float32 { return []float32{unit(st.Y > o.After)} })
	return root
}
