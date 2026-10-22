package navigation

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToTopOptions tune a ToTop.
type ToTopOptions struct {
	// After is how far down the button appears.
	After float32
}

// ToTop appears past a scroll and scrolls st back up.
func ToTop(c *ui.Context, label string, st *ui.ScrollState, o ToTopOptions) (root, up ui.Element) {
	if label == "" || st == nil || o.After < 0 || math.IsNaN(float64(o.After)) || math.IsInf(float64(o.After), 0) {
		panic(fmt.Sprintf("navigation: ToTop %q needs a label, a scroll and a finite After", label))
	}
	t := c.Theme()
	if o.After == 0 {
		o.After = t.Space(60)
	}
	root = ui.Box(c)
	shown := st.Y > o.After
	root.Children(func() {
		if !shown {
			return
		}
		c.Key("up")
		up = kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Size(t.Space(10), t.Space(10)).Center().Radius(999)).
			Transition(ui.ElementTransition{Enter: &ui.Motion{Y: t.Space(4)}, Exit: &ui.Motion{Y: t.Space(4)}})
		up.Background(t.Background).Border(1, t.Border).Shadow(0, 2, 8, 0, kit.Shade(0.12)).Children(func() {
			ui.Icon(c, icon.Must("arrow-up")).Size(16, 16)
		})
		if kit.Clicked(up, sound.Tap) {
			st.ScrollTo(st.X, 0, ui.Smooth)
			kit.Log("ToTop", "scrolled up")
		}
	})
	kit.Settle(c, root, []float32{unit(shown)}, func() []float32 { return []float32{unit(st.Y > o.After)} })
	return root, up
}
