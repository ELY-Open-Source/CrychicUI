package motion

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WordRotateOptions tune a WordRotate.
type WordRotateOptions struct {
	// Every is how long each word shows; 2 s.
	Every time.Duration
}

// WordRotate cycles words, each sliding up; Changed marks rotations.
func WordRotate(c *ui.Context, words []string, o WordRotateOptions) ui.Element {
	if len(words) == 0 || o.Every < 0 {
		panic("motion: WordRotate needs words and a positive Every")
	}
	for _, w := range words {
		if w == "" {
			panic("motion: WordRotate words need text")
		}
	}
	if o.Every == 0 {
		o.Every = 2 * time.Second
	}
	h := c.Theme().Rem(1.25)
	root := ui.Box(c).ClipY()
	start := ui.Local(root, "start", func() time.Time { return c.Now() })
	n := int(c.Now().Sub(*start) / o.Every)
	if len(words) > 1 {
		c.After(start.Add(time.Duration(n+1) * o.Every).Sub(c.Now()))
	}
	i := n % len(words)
	last := ui.Local(root, "last", func() int { return n })
	if n != *last {
		*last = n
		kit.Log("WordRotate", "rotated", "index", i)
		root.SetChanged()
	}
	slide := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{Y: h}, Exit: &ui.Motion{Y: -h}}
	return root.Children(func() {
		c.Key(n)
		ui.Text(c, words[i]).Transition(slide)
	})
}
