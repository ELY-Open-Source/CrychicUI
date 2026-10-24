package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AddToCartOptions word the button and its brief confirmation.
type AddToCartOptions struct{ Label, Added string }

// flightTime is how long the dot and check last.
const flightTime = 900 * time.Millisecond

// AddToCart adds with a flying dot and check; cue Pop.
func AddToCart(c *ui.Context, label string, added *bool, o AddToCartOptions) ui.Element {
	if label == "" || added == nil {
		panic(fmt.Sprintf("shop: AddToCart %q needs a label and state", label))
	}
	def(&o.Label, "Add to cart")
	def(&o.Added, "Added")
	t := c.Theme()
	root := ui.Grid(c).Columns(1).AlignSelf(ui.Start)
	flight := ui.Local(root, "flight", func() int { return 0 })
	since := ui.Local(root, "since", func() time.Time { return time.Time{} })
	flying := c.Now().Sub(*since) < flightTime
	if flying {
		c.After(flightTime - c.Now().Sub(*since))
	}
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(label).ColumnStart(1).RowStart(1).Padding(t.Space(2), t.Space(4)).Radius(t.Radius))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())).Mix(t.Success, unit(flying)), ui.Snappy)).TextColor(t.AccentText)
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph, text := "shopping-bag", o.Label
				if flying {
					glyph, text = "check", o.Added
				}
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
				ui.Text(c, text).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
			})
		})
		if b.Clicked() {
			*added = true
			*flight++
			*since = c.Now()
			c.After(flightTime)
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("AddToCart", "added", "times", *flight)
		}
		if flying {
			// Keyed by flight: each add launches a new dot.
			c.Key(*flight)
			dot := ui.Box(c).ColumnStart(1).RowStart(1).AlignSelf(ui.Center).JustifySelf(ui.Center).Size(t.Space(2.5), t.Space(2.5)).
				Radius(999).Background(t.Accent).Border(t.Space(0.5), t.Background)
			k := dot.Keyframes("crychic.fly", []ui.Keyframe{{Value: 0}, {At: 450 * time.Millisecond, Value: 1, Ease: ui.EaseInOut}})
			dot.Translate(t.Space(20)*k, -t.Space(10)*k*(2-k)).Scale(1-0.5*k, 1-0.5*k).Opacity(1 - k*k)
		}
	})
	return root
}
