package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// StaggerOptions tune a Stagger.
type StaggerOptions struct {
	// Step delays each item after the one before; 40 ms.
	Step time.Duration
}

// Stagger builds n items, each arriving a step later.
func Stagger(c *ui.Context, n int, item func(i int), o StaggerOptions) ui.Element {
	if n < 0 || item == nil || o.Step < 0 {
		panic("motion: Stagger needs n ≥ 0, an item and a positive Step")
	}
	if o.Step == 0 {
		o.Step = 40 * time.Millisecond
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2))
	start := ui.Local(root, "start", func() time.Time { return c.Now() })
	return root.Children(func() {
		for i := range n {
			c.Key(i)
			w := ui.Box(c)
			delay := ui.Local(w, "delay", func() time.Duration { return max(start.Add(time.Duration(i)*o.Step).Sub(c.Now()), 0) })
			k := w.Keyframes("crychic.in", arrive, ui.KeyframeOptions{Delay: *delay})
			w.Opacity(k).Translate(0, (1-k)*t.Space(2)).Children(func() { item(i) })
		}
	})
}
