package chat

import (
	"math"
	"time"

	"github.com/egoist/mygo/ui"
)

// Typing is three dots bouncing in a bubble.
func Typing(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("chat: Typing needs a label")
	}
	t := c.Theme()
	root := shell(c, false).Row().Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(label)
	ph := float64(root.Loop("bounce", 900*time.Millisecond, ui.Linear))
	return root.Children(func() {
		for i := range 3 {
			c.Key(i)
			lift := math.Max(0, math.Sin(2*math.Pi*(ph-float64(i)*0.15)))
			ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Radius(999).Background(t.TextMuted).
				Translate(0, -float32(lift)*t.Space(1.5))
		}
	})
}
