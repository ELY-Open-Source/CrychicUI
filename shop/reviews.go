package shop

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Review is one customer review; Title is optional.
type Review struct {
	Rating        int
	Author, Title string
	Text          string
	Helpful       int
}

// ReviewsOptions word the summary; Count takes the number of reviews.
type ReviewsOptions struct{ Count string }

func validReview(component string, r Review) {
	if r.Rating < 1 || r.Rating > 5 || r.Author == "" || r.Text == "" || r.Helpful < 0 {
		panic(fmt.Sprintf("shop: %s review by %q needs a rating 1..5, text and helpful votes", component, r.Author))
	}
}

// Reviews shows the average, bars per star and the reviews.
func Reviews(c *ui.Context, reviews []Review, o ReviewsOptions) ui.Element {
	if len(reviews) == 0 {
		panic("shop: Reviews needs reviews")
	}
	total, counts := 0, [5]int{}
	for _, r := range reviews {
		validReview("Reviews", r)
		total += r.Rating
		counts[r.Rating-1]++
	}
	def(&o.Count, "%d reviews")
	mean := float64(total) / float64(len(reviews))
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(fmt.Sprintf("%.1f", mean))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(5)).Padding(t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Column(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
				kit.Numeric(c, fmt.Sprintf("%.1f", mean)).FontSize(t.Rem(2.2)).Bold()
				stars(c, mean, t.Space(3.5))
				kit.Muted(c, fmt.Sprintf(o.Count, len(reviews))).FontSize(t.Rem(0.75))
			})
			g := root.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Justify(ui.Center).Children(func() {
				for n := 5; n >= 1; n-- {
					c.Key(n)
					share := float32(counts[n-1]) / float32(len(reviews))
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						kit.Numeric(c, strconv.Itoa(n)).FontSize(t.Rem(0.75)).Width(t.Space(3)).Shrink(0)
						bar := ui.Box(c).Grow(1).Height(t.Space(2)).Radius(999).Background(t.Surface)
						w := bar.AnimateTo("share", share, ui.Smooth) * g
						bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
							if w > 0 {
								p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*w), H: r.H}, t.Warning, 999)
							}
						}))
						kit.Numeric(c, strconv.Itoa(counts[n-1])).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Width(t.Space(6)).Shrink(0).TextAlign(ui.End)
					})
				}
			})
		})
		ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(3), t.Space(2), t.Space(3)).Children(func() {
			for i, r := range reviews {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(2), 0).BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(reviews)-1)))
				kit.Stagger(c, root, row, i).Children(func() {
					display.Avatar(c, r.Author, nil).Size(t.Space(8), t.Space(8)).FontSize(t.Rem(0.7)).Shrink(0)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.75)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, r.Author).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							stars(c, float64(r.Rating), t.Space(3))
						})
						if r.Title != "" {
							ui.Text(c, r.Title).Bold().FontSize(t.Rem(0.85)).Wrap()
						}
						ui.Text(c, r.Text).FontSize(t.Rem(0.85)).TextColor(t.TextMuted).Wrap()
					})
				})
			}
		})
	})
	return root
}
