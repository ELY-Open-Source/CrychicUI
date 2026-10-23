package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Review is one customer review.
type Review struct {
	Rating  int
	Author  string
	Text    string
	Helpful int
}

// ReviewsOptions are the review labels.
type ReviewsOptions struct{ Count string }

// Reviews shows an average, star bars, and reviews.
func Reviews(c *ui.Context, reviews []Review, o ReviewsOptions) ui.Element {
	if len(reviews) == 0 {
		panic("shop: Reviews needs reviews")
	}
	def(&o.Count, "%d reviews")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	total, counts := 0, [5]int{}
	for _, review := range reviews {
		if review.Rating < 1 || review.Rating > 5 || review.Author == "" || review.Text == "" || review.Helpful < 0 {
			panic("shop: Reviews has an invalid review")
		}
		total += review.Rating
		counts[review.Rating-1]++
	}
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			kit.Numeric(c, fmt.Sprintf("%.1f", float64(total)/float64(len(reviews)))).Bold().FontSize(t.Rem(1.5))
			kit.Muted(c, fmt.Sprintf(o.Count, len(reviews))).FontSize(t.Rem(0.8))
		})
		for stars := 5; stars >= 1; stars-- {
			c.Key(fmt.Sprintf("stars.%d", stars))
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, fmt.Sprint(stars)).FontSize(t.Rem(0.75)).Width(t.Space(3))
				bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
				fraction := bar.AnimateTo("rating", float32(counts[stars-1])/float32(len(reviews)), ui.Smooth)
				bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fraction, H: r.H}, t.Warning, 999)
				}))
				kit.Numeric(c, fmt.Sprint(counts[stars-1])).FontSize(t.Rem(0.75))
			})
		}
		for i, review := range reviews {
			c.Key(fmt.Sprintf("review.%d", i))
			ui.Column(c).Gap(t.Space(1)).Children(func() {
				ui.Row(c).Gap(t.Space(1)).Children(func() {
					ui.Text(c, review.Author).Bold().Grow(1)
					for range review.Rating {
						ui.Icon(c, icon.Must("star")).Size(16, 16).TextColor(t.Warning)
					}
				})
				ui.Text(c, review.Text).FontSize(t.Rem(0.85)).Wrap()
			})
		}
	})
	return root
}
