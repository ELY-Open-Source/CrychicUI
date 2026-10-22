package overlay

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// hoverGrace keeps the card while the pointer crosses the gap.
const hoverGrace = 300 * time.Millisecond

// HoverCard shows a card while anchor or card is hovered.
func HoverCard(c *ui.Context, label string, anchor ui.Element, content func()) ui.Element {
	if label == "" || content == nil {
		panic(fmt.Sprintf("overlay: HoverCard %q needs content", label))
	}
	t := c.Theme()
	root := ui.Box(c)
	shown := ui.Local(root, "shown", func() bool { return false })
	hideAt := ui.Local(root, "hideAt", func() time.Time { return time.Time{} })
	now := c.Now()
	if anchor.Hovered() || anchor.Focused() {
		*shown, *hideAt = true, time.Time{}
	} else if !hideAt.IsZero() && !now.Before(*hideAt) {
		*shown, *hideAt = false, time.Time{}
	}
	if *shown {
		ui.Overlay(c, func() {
			card := ui.Column(c).AttachTo(anchor, ui.AnchorBottom, ui.AnchorTop).Margin(t.Space(2), 0, 0, 0).Width(t.Space(72)).
				Padding(t.Space(4)).Gap(t.Space(2)).Radius(t.Radius).Background(t.Background).Border(1, t.Border).
				Shadow(0, 8, 24, 0, kit.Shade(0.16)).Label(label)
			kit.Arrive(card, ui.Pose{Y: -t.Space(1), Fade: 1}, 0).Children(content)
			switch {
			case card.Hovered() || anchor.Hovered() || anchor.Focused():
				*hideAt = time.Time{}
			case hideAt.IsZero():
				*hideAt = now.Add(hoverGrace)
			}
		})
		if !hideAt.IsZero() {
			c.After(hideAt.Sub(now))
		}
	}
	moved(root, *shown, "HoverCard")
	return root
}
