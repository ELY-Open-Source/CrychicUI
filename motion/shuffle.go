package motion

import (
	"math/rand/v2"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShuffleOptions tune a Shuffle.
type ShuffleOptions struct {
	// Label names the shuffle button; "Shuffle" when empty.
	Label string
}

// Shuffle slides items to a new order on its button.
func Shuffle(c *ui.Context, items []string, o ShuffleOptions) ui.Element {
	unique("Shuffle", items)
	if o.Label == "" {
		o.Label = "Shuffle"
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Start)
	order := ui.Local(root, "order", func() []int { return identity(len(items)) })
	if len(*order) != len(items) {
		*order = identity(len(items))
	}
	slide := ui.ElementTransition{Animation: &ui.Snappy}
	return root.Children(func() {
		if kit.Clicked(ui.Button(c, o.Label), sound.Swipe) {
			next := rand.Perm(len(items))
			if slices.Equal(next, *order) {
				next = append(next[1:], next[0])
			}
			*order = next
			root.SetChanged()
			kit.Log("Shuffle", "shuffled", "count", len(items))
		}
		ui.Column(c).Gap(t.Space(1.5)).Children(func() {
			for _, i := range *order {
				c.Key(items[i])
				ui.Box(c).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
					Transition(slide).Children(func() { ui.Text(c, items[i]) })
			}
		})
	})
}
