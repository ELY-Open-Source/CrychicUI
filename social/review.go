package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ReviewData is one review; Title and Time are optional.
type ReviewData struct {
	ID, Title, Time string
	Author          User
	Rating          int
	Text            string
	Photos          []ui.ImageSource
}

// Review shows a star review with photos.
func Review(c *ui.Context, review ReviewData) ui.Element {
	validUser("Review", review.Author)
	if review.ID == "" || review.Rating < 1 || review.Rating > 5 || review.Text == "" {
		panic(fmt.Sprintf("social: Review %q needs id, text and a 1-5 rating", review.ID))
	}
	for i, photo := range review.Photos {
		if photo == nil {
			panic(fmt.Sprintf("social: Review %q photo %d is nil", review.ID, i))
		}
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(2.5)).Padding(t.Space(3.5))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Children(func() { byline(c, review.Author, review.Time, t.Space(10)) })
			stars(c, review.Rating, t.Space(4))
		})
		if review.Title != "" {
			ui.Text(c, review.Title).FontSize(t.Rem(0.92)).Bold().Wrap().FillWidth()
		}
		ui.Text(c, review.Text).FontSize(t.Rem(0.85)).Wrap().FillWidth()
		if len(review.Photos) == 0 {
			return
		}
		c.Key("photos")
		strip := ui.Row(c).FillWidth().Gap(t.Space(1.5))
		strip.Children(func() {
			shown := min(len(review.Photos), 4)
			for i, photo := range review.Photos[:shown] {
				c.Key(i)
				thumb := ui.Box(c).Size(t.Space(18), t.Space(18)).Shrink(0).Radius(t.Radius).Clip().Background(t.SurfaceHover)
				kit.Stagger(c, strip, thumb, i).Children(func() {
					ui.Image(c, photo).Fill().Fit(ui.Cover)
					if i == shown-1 && len(review.Photos) > shown {
						ui.Box(c).Absolute().Top(0).Left(0).Fill().Center().Background(t.Text.Alpha(0.5)).TextColor(t.Background).Children(func() {
							ui.Text(c, fmt.Sprintf("+%d", len(review.Photos)-shown)).FontSize(t.Rem(0.9)).Bold().Shrink(0)
						})
					}
				})
			}
		})
	})
	return root
}
