package shop

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReviewCardOptions are the review card labels.
type ReviewCardOptions struct{ Helpful string }

// ReviewCard shows a review with helpful; cue Pop.
func ReviewCard(c *ui.Context, review Review, helpful *bool, o ReviewCardOptions) ui.Element {
	if review.Rating < 1 || review.Rating > 5 || review.Author == "" || review.Text == "" || review.Helpful < 0 || helpful == nil {
		panic("shop: ReviewCard has an invalid review")
	}
	def(&o.Helpful, "Helpful (%d)")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			for range review.Rating {
				ui.Icon(c, icon.Must("star")).Size(16, 16).TextColor(t.Warning)
			}
			ui.Text(c, review.Author).Bold().Grow(1)
		})
		ui.Text(c, review.Text).Wrap()
		count := review.Helpful
		if *helpful && count < math.MaxInt {
			count++
		}
		b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Helpful, count)).AlignSelf(ui.Start).Padding(t.Space(1), t.Space(2)).
			Radius(999).Border(1, t.Border).Checked(*helpful).Disabled(*helpful || review.Helpful == math.MaxInt))
		b.Children(func() { ui.Text(c, fmt.Sprintf(o.Helpful, count)).FontSize(t.Rem(0.8)) })
		if b.Clicked() && !*helpful && review.Helpful < math.MaxInt {
			*helpful = true
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("ReviewCard", "helpful")
		}
	})
	return root
}
