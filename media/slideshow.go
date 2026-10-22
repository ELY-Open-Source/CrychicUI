package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SlideshowOptions tune a Slideshow.
type SlideshowOptions struct {
	// Every is a slide's time; zero gives 4 s.
	Every time.Duration
	// Height of a slide; zero gives the default.
	Height float32
}

// Slideshow advances slides on a timer under a progress bar.
func Slideshow(c *ui.Context, label string, n int, at *int, slide func(i int), o SlideshowOptions) ui.Element {
	if label == "" || n < 1 || slide == nil || o.Every < 0 {
		panic(fmt.Sprintf("media: Slideshow %q needs a label, slides, a builder and a time not negative", label))
	}
	chosen("Slideshow", label, at, 0, n)
	if o.Every == 0 {
		o.Every = 4 * time.Second
	}
	t := c.Theme()
	size("Slideshow", label, &o.Height, t.Space(48))
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Label(label)
	now := c.Now()
	since := ui.Local(root, "since", func() time.Time { return now })
	last := ui.Local(root, "at", func() int { return *at })
	if *last != *at {
		*since, *last = now, *at
	}
	if n > 1 && now.Sub(*since) >= o.Every {
		*at = (*at + 1) % n
		*since, *last = now, *at
		root.SetChanged()
		kit.Log("Slideshow", "advanced", "index", *at)
	}
	if n > 1 {
		c.After(o.Every - now.Sub(*since))
	}
	shown := *at
	return root.Children(func() {
		ui.Box(c).FillWidth().Height(o.Height).Radius(t.Radius).Clip().Background(t.Surface).Children(func() {
			kit.Slide(c, shown, func() { slide(shown) })
		})
		bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Background(t.Border).Clip()
		bar.Children(func() {
			fill := ui.Box(c).Height(t.Space(1)).Radius(999).Background(t.Accent)
			p := float32(1)
			if n > 1 {
				p = fill.Keyframes(*since, []ui.Keyframe{{Value: 0}, {At: o.Every, Value: 1}})
			}
			fill.WidthPercent(100 * p)
		})
	})
}
