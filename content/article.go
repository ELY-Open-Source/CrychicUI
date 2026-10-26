package content

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ArticleOptions dress an article; zero Date or Minutes hides it.
type ArticleOptions struct {
	Author, Role, Tag string
	Photo, Cover      *ui.Bitmap
	Date              time.Time
	Minutes           int
	// DateFormat is a layout; Read formats minutes.
	DateFormat, Read string
}

// Article is a title, a byline and paragraphs.
func Article(c *ui.Context, title string, body []string, o ArticleOptions) ui.Element {
	def(&o.DateFormat, "Jan 2, 2006")
	def(&o.Read, "%d min read")
	if strings.TrimSpace(title) == "" || strings.TrimSpace(o.Author) == "" || len(body) == 0 || o.Minutes < 0 || o.Minutes > 600 {
		panic(fmt.Sprintf("content: Article %q needs an author, paragraphs and 0 to 600 minutes", title))
	}
	for _, p := range body {
		if strings.TrimSpace(p) == "" {
			panic(fmt.Sprintf("content: Article %q has a blank paragraph", title))
		}
	}
	t := c.Theme()
	root := card(c).Clip().Role(ui.RoleGroup)
	root.Children(func() {
		if o.Cover != nil {
			cover := ui.Box(c).FillWidth().AspectRatio(16.0 / 5).Clip()
			cover.Children(func() {
				kit.Arrive(ui.Image(c, o.Cover).Fill().Fit(ui.Cover), ui.Pose{Scale: 1.06, Fade: 1}, 0)
			})
		}
		ui.Column(c).FillWidth().Gap(t.Space(4)).Padding(t.Space(6)).Children(func() {
			ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
				if o.Tag != "" {
					ui.Text(c, o.Tag).FontSize(t.Rem(0.72)).Bold().TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine()
				}
				kit.Arrive(ui.Text(c, title).FontSize(t.Rem(1.4)).Bold().FillWidth().Wrap(), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			})
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
				avatar(c, o.Author, o.Photo, 9)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, o.Author).FontSize(t.Rem(0.86)).Bold().FillWidth().MinWidth(0).SingleLine()
					date, read := "", ""
					if !o.Date.IsZero() {
						date = o.Date.Format(o.DateFormat)
					}
					if o.Minutes > 0 {
						read = fmt.Sprintf(o.Read, o.Minutes)
					}
					meta(c, o.Role, date, read)
				})
			})
			ui.Box(c).FillWidth().Height(1).Background(t.Border)
			list := ui.Column(c).FillWidth().Gap(t.Space(3))
			list.Children(func() {
				for i, p := range body {
					c.Key(i)
					size, ink := t.Rem(0.9), t.Text.Alpha(0.86)
					if i == 0 {
						size, ink = t.Rem(1.02), t.Text
					}
					kit.Stagger(c, list, ui.Text(c, p).FontSize(size).TextColor(ink).LineHeight(1.65).FillWidth().Wrap(), i)
				}
			})
		})
	})
	return root
}
