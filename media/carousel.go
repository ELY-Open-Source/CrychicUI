package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CarouselOptions tune a Carousel.
type CarouselOptions struct {
	// Height of a slide; zero gives the default.
	Height float32
	// Slide names a dot by number; "Slide %d" by default.
	Slide, Back, Forth string
}

// Carousel springs between slides; arrows, dots, keys, drags.
func Carousel(c *ui.Context, label string, n int, at *int, slide func(i int), o CarouselOptions) ui.Element {
	if label == "" || n < 1 || slide == nil {
		panic(fmt.Sprintf("media: Carousel %q needs a label, slides and a builder", label))
	}
	chosen("Carousel", label, at, 0, n)
	t := c.Theme()
	size("Carousel", label, &o.Height, t.Space(48))
	def(&o.Slide, "Slide %d")
	def(&o.Back, "Previous")
	def(&o.Forth, "Next")
	root := ui.Column(c).FillWidth().Gap(t.Space(3))
	shown := *at
	turn := func(to int) {
		if to < 0 || to >= n || to == shown {
			return
		}
		*at, shown = to, to
		root.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("Carousel", "shown", "index", to)
	}
	back, forth, prev, next := ends(c)
	step := 1
	if c.Direction() == ui.RightToLeft {
		step = -1
	}
	return root.Children(func() {
		view := ui.Grid(c).Columns(1).FillWidth().Height(o.Height).Radius(t.Radius).Clip().Background(t.Surface).
			Focusable().Label(label)
		w := view.LayoutBounds().W
		kit.Settle(c, view, []float32{w}, func() []float32 { return []float32{view.LayoutBounds().W} })
		drag := ui.Local(view, "drag", func() float32 { return 0 })
		was := ui.Local(view, "held", func() bool { return false })
		dx, _, held := view.Dragged()
		if held {
			*drag += dx
		}
		if *was && !held {
			switch {
			case *drag < -w/5:
				turn(shown + step)
			case *drag > w/5:
				turn(shown - step)
			}
			*drag = 0
		}
		*was = held
		switch {
		case view.Shortcut(0, prev):
			turn(shown - 1)
		case view.Shortcut(0, next):
			turn(shown + 1)
		}
		var track ui.Element
		var arrows [2]ui.Element
		view.Children(func() {
			track = ui.Row(c).ColumnStart(1).RowStart(1).Height(o.Height).Children(func() {
				for i := range n {
					c.Key(i)
					ui.Box(c).Width(w).Height(o.Height).Shrink(0).Children(func() { slide(i) })
				}
			})
			for k, a := range [2][2]string{{o.Back, back}, {o.Forth, forth}} {
				c.Key(k)
				b := round(c, a[0], a[1], t.Space(9)).Absolute().TopPercent(50).Translate(0, -t.Space(4.5))
				b.Background(b.AnimateColor("plate", t.Background.Alpha(0.75+0.25*unit(b.Hovered())), ui.Snappy))
				arrows[k] = b
			}
			arrows[0].Start(t.Space(2))
			arrows[1].End(t.Space(2))
		})
		dots := make([]ui.Element, n)
		ui.Row(c).FillWidth().Justify(ui.Center).Gap(t.Space(1.5)).Children(func() {
			for i := range n {
				c.Key(i)
				name := fmt.Sprintf(o.Slide, i+1)
				dots[i] = ui.ButtonBase(c).Label(name).Tooltip(name).Height(t.Space(2)).Radius(999)
			}
		})
		switch {
		case arrows[0].Clicked():
			turn(shown - 1)
		case arrows[1].Clicked():
			turn(shown + 1)
		}
		for i, d := range dots {
			if d.Clicked() {
				turn(i)
			}
		}
		arrows[0].Disabled(shown == 0)
		arrows[1].Disabled(shown == n-1)
		for i, d := range dots {
			ink := t.Border
			if i == shown {
				ink = t.Accent
			}
			d.Checked(i == shown).Width(d.AnimateTo("w", t.Space(2+4*unit(i == shown)), ui.Snappy)).
				Background(d.AnimateColor("ink", ink, ui.Snappy))
		}
		track.Translate(view.AnimateTo("x", -float32(step*shown)*w+*drag, ui.Bouncy), 0)
	})
}
