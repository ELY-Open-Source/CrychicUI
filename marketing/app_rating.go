package marketing

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AppRatingOptions word the summary; Ratings and One take the total.
type AppRatingOptions struct{ Store, Icon, Ratings, One string }

// AppRating is a store rating summary; counts[0] holds five stars.
func AppRating(c *ui.Context, score float64, counts [5]int, o AppRatingOptions) ui.Element {
	def(&o.Store, "App Store")
	def(&o.Icon, "smartphone")
	def(&o.Ratings, "%s ratings")
	def(&o.One, "%s rating")
	if !(score >= 0 && score <= 5) {
		panic(fmt.Sprintf("marketing: AppRating %q needs a score of 0 to 5, got %v", o.Store, score))
	}
	total, most := 0, 0
	for _, n := range counts {
		if n < 0 || n > math.MaxInt-total {
			panic(fmt.Sprintf("marketing: AppRating %q needs counts of 0 or more that sum without overflow", o.Store))
		}
		total += n
		most = max(most, n)
	}
	t := c.Theme()
	label := o.Ratings
	if total == 1 {
		label = o.One
	}
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			tile(c, o.Icon, t.Accent, t.Space(9))
			ui.Text(c, o.Store).FontSize(t.Rem(0.95)).Bold().Grow(1).MinWidth(0).SingleLine()
		})
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(5)).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				kit.Numeric(c, decimal(c, score, 1)).FontSize(t.Rem(2.6)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
				rating(c, score, 16)
				kit.Muted(c, fmt.Sprintf(label, number(c, total))).FontSize(t.Rem(0.74)).Shrink(0).MaxWidth(t.Space(40)).SingleLine()
			})
			bars := ui.Column(c).Grow(1).Basis(t.Space(40)).MinWidth(0).Gap(t.Space(1.5))
			bars.Children(func() {
				for i, n := range counts {
					share := float32(0)
					if most > 0 {
						share = float32(n) / float32(most)
					}
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						kit.Muted(c, strconv.Itoa(5-i)).FontSize(t.Rem(0.72)).Width(t.Space(2.5)).Shrink(0).TextAlign(ui.End)
						bar := ui.Box(c).Grow(1).MinWidth(0).Height(t.Space(2)).Radius(999).Background(t.SurfaceHover)
						grow := bar.Keyframes("grow", []ui.Keyframe{{Value: 0},
							{At: time.Duration(i)*80*time.Millisecond + 500*time.Millisecond, Value: 1, Ease: ui.EaseOut}})
						fill := bar.AnimateTo("share", share, ui.Smooth) * grow
						bar.Draw(func(p *ui.Painter, r ui.Rect) {
							if fill > 0 {
								p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*fill), H: r.H}, t.Warning, r.H/2)
							}
						})
						kit.Muted(c, number(c, n)).FontSize(t.Rem(0.7)).Width(t.Space(12)).Shrink(0).TextAlign(ui.End).MinWidth(0).SingleLine()
					})
				}
			})
		})
	})
	return root
}
