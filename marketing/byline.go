package marketing

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BylineOptions word the byline; zero Minutes hides read time.
type BylineOptions struct {
	Role    string
	Minutes int
	PostWords
}

// Byline is an author with the date and read time.
func Byline(c *ui.Context, author string, photo *ui.Bitmap, date time.Time, o BylineOptions) ui.Element {
	o.fill()
	if strings.TrimSpace(author) == "" || date.IsZero() || o.Minutes < 0 || o.Minutes > 600 {
		panic(fmt.Sprintf("marketing: Byline %q needs a date and 0 to 600 minutes", author))
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
	kit.Arrive(root, ui.Pose{Y: t.Space(1.5), Fade: 1}, 0).Children(func() {
		kit.Arrive(avatar(c, author, photo, 10), ui.Pose{Scale: 0.7, Fade: 1}, 80*time.Millisecond)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Row(c).FillWidth().MinWidth(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, author).FontSize(t.Rem(0.88)).Bold().Shrink(0).MaxWidth(t.Space(44)).SingleLine()
				if o.Role != "" {
					dot(c)
					kit.Muted(c, o.Role).FontSize(t.Rem(0.78)).MinWidth(0).SingleLine()
				}
			})
			meta(c, Post{Date: date, Minutes: o.Minutes}, o.PostWords)
		})
	})
	return root
}
