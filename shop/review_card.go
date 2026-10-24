package shop

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReviewCardOptions word the helpful vote, "%d" taking the count.
type ReviewCardOptions struct{ Helpful string }

// ReviewCard shows a review and a helpful vote; cue Pop.
func ReviewCard(c *ui.Context, review Review, helpful *bool, o ReviewCardOptions) ui.Element {
	validReview("ReviewCard", review)
	if helpful == nil {
		panic(fmt.Sprintf("shop: ReviewCard by %q needs bound state", review.Author))
	}
	def(&o.Helpful, "Helpful (%d)")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Role(ui.RoleGroup).Label(review.Author)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	count := review.Helpful
	if *helpful && count < math.MaxInt {
		count++
	}
	var b ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			display.Avatar(c, review.Author, nil).Size(t.Space(9), t.Space(9)).FontSize(t.Rem(0.75)).Shrink(0)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, review.Author).Bold().FillWidth().SingleLine()
				stars(c, float64(review.Rating), t.Space(3.5))
			})
		})
		if review.Title != "" {
			ui.Text(c, review.Title).Bold().Wrap()
		}
		ui.Text(c, review.Text).FontSize(t.Rem(0.9)).LineHeight(1.5).Wrap()
		label := fmt.Sprintf(o.Helpful, count)
		// Keyed: an optional title above must not move its focus.
		c.Key("helpful")
		b = kit.Press(ui.ButtonBase(c).Label(label).AlignSelf(ui.Start).Padding(t.Space(1), t.Space(2.5)).Radius(999).Checked(*helpful)).
			Disabled(*helpful || review.Helpful == math.MaxInt)
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(*helpful)), ui.Snappy))
		b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(*helpful)), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.Accent, unit(*helpful)), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				thumb := ui.Icon(c, icon.Must("thumbs-up")).Size(14, 14)
				lift := thumb.AnimateTo("lift", unit(*helpful), ui.Bouncy)
				thumb.Rotate(-12*lift).Scale(1+0.15*lift*(1-lift)*4, 1+0.15*lift*(1-lift)*4)
				ui.Text(c, label).FontSize(t.Rem(0.8)).Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			})
		})
	})
	if b.Clicked() && !*helpful && review.Helpful < math.MaxInt {
		*helpful = true
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("ReviewCard", "helpful")
	}
	return root
}
