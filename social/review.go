package social

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ReviewData is one star review with photos.
type ReviewData struct {
	ID     string
	Author User
	Rating int
	Text   string
	Photos []ui.ImageSource
}

// Review shows a star review with photos.
func Review(c *ui.Context, review ReviewData) ui.Element {
	validUser("Review", review.Author)
	if review.ID == "" || review.Rating < 1 || review.Rating > 5 || review.Text == "" {
		panic("social: invalid review")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			avatar(c, review.Author, t.Space(8))
			ui.Text(c, review.Author.Name).Bold()
		})
		ui.Row(c).Gap(t.Space(0.5)).Children(func() {
			for i := 1; i <= 5; i++ {
				color := t.TextMuted
				if i <= review.Rating {
					color = t.Warning
				}
				ui.Icon(c, icon.Must("star")).Size(16, 16).TextColor(color)
			}
		})
		ui.Text(c, review.Text).Wrap()
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			for i, photo := range review.Photos {
				if photo == nil {
					panic("social: review photo must not be nil")
				}
				c.Key(i)
				ui.Image(c, photo).Size(t.Space(20), t.Space(20)).Fit(ui.Cover).Radius(t.Radius).Clip()
			}
		})
	})
	return root
}
