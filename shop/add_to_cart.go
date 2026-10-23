package shop

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AddToCartOptions are the add labels.
type AddToCartOptions struct{ Label string }

// AddToCart adds with a flying dot; cue Pop.
func AddToCart(c *ui.Context, label string, added *bool, o AddToCartOptions) ui.Element {
	if label == "" || added == nil {
		panic("shop: AddToCart needs a label and state")
	}
	def(&o.Label, "Add to cart")
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	flight := ui.Local(root, "flight", func() int { return 0 })
	until := ui.Local(root, "until", func() time.Time { return time.Time{} })
	root.Children(func() {
		button := addButton(c, label).Row().Gap(t.Space(1))
		button.Children(func() {
			ui.Icon(c, icon.Must("shopping-bag")).Size(16, 16)
			ui.Text(c, o.Label).SingleLine()
		})
		if button.Clicked() {
			*added = true
			*flight++
			*until = c.Now().Add(250 * time.Millisecond)
			c.After(250 * time.Millisecond)
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("AddToCart", "added")
		}
		if *added && c.Now().Before(*until) {
			c.Key(*flight)
			dot := ui.Box(c).Absolute().Right(0).Top(0).Size(t.Space(2), t.Space(2)).Radius(999).Background(t.Accent)
			kit.Arrive(dot, ui.Pose{X: t.Space(8), Y: -t.Space(6), Scale: 0.2, Fade: 1}, 0)
		}
	})
	return root
}
