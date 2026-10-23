package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StoriesOptions tune a Stories.
type StoriesOptions struct {
	// Every is a slide's time; zero gives 5 s.
	Every time.Duration
	// Height of a slide; zero gives the default.
	Height float32
}

// Stories advances slides under segment bars; a tap skips.
func Stories(c *ui.Context, label string, n int, at *int, slide func(i int), o StoriesOptions) ui.Element {
	if label == "" || n < 1 || slide == nil || o.Every < 0 {
		panic(fmt.Sprintf("media: Stories %q needs a label, slides, a builder and a time not negative", label))
	}
	chosen("Stories", label, at, 0, n)
	if o.Every == 0 {
		o.Every = 5 * time.Second
	}
	t := c.Theme()
	size("Stories", label, &o.Height, t.Space(80))
	root := ui.ButtonBase(c).Label(label).FillWidth().Height(o.Height).Radius(t.Radius).Clip().Background(t.Surface)
	now := c.Now()
	since := ui.Local(root, "since", func() time.Time { return now })
	last := ui.Local(root, "at", func() int { return *at })
	if *last != *at {
		*since, *last = now, *at
	}
	step := func(to int, event string) {
		to = max(0, min(n-1, to))
		if to == *at {
			return
		}
		*at = to
		*since, *last = now, to
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Stories", event, "index", to)
	}
	b := root.LayoutBounds()
	kit.Settle(c, root, []float32{b.W}, func() []float32 { return []float32{root.LayoutBounds().W} })
	_, _, prev, next := ends(c)
	switch {
	case root.Shortcut(0, ui.KeyEnter), root.Shortcut(0, ui.KeySpace), root.Shortcut(0, next):
		step(*at+1, "skipped")
	case root.Shortcut(0, prev):
		step(*at-1, "back")
	case root.Clicked():
		x, _, _ := root.PointerPosition()
		back := x < b.W/3
		if c.Direction() == ui.RightToLeft {
			back = x > b.W*2/3
		}
		if back {
			step(*at-1, "back")
		} else {
			step(*at+1, "skipped")
		}
	}
	if *at < n-1 && now.Sub(*since) >= o.Every {
		*at++
		*since, *last = now, *at
		root.SetChanged()
		kit.Log("Stories", "advanced", "index", *at)
	}
	if *at < n-1 {
		c.After(o.Every - now.Sub(*since))
	}
	shown := *at
	ink := t.Background
	if t.Dark {
		ink = t.Text
	}
	return root.Children(func() {
		ui.Box(c).FillWidth().Height(o.Height).Children(func() {
			kit.Slide(c, shown, func() {
				ui.Box(c).FillWidth().Height(o.Height).Children(func() { slide(shown) })
			})
		})
		ui.Row(c).Absolute().Left(0).Right(0).Top(0).Gap(t.Space(1)).Padding(t.Space(2)).Children(func() {
			for i := range n {
				c.Key(i)
				seg := ui.Box(c).Grow(1).Basis(0).Height(t.Space(0.75)).Radius(999).Background(ink.Alpha(0.35)).Clip()
				seg.Children(func() {
					fill := ui.Box(c).Height(t.Space(0.75)).Radius(999).Background(ink)
					p := unit(i < shown)
					if i == shown {
						p = fill.Keyframes(*since, []ui.Keyframe{{Value: 0}, {At: o.Every, Value: 1}})
					}
					fill.WidthPercent(100 * p)
				})
			}
		})
	})
}
