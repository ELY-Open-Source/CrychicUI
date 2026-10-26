package geo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Place is a spot; zero Rating hides the stars.
type Place struct {
	Name, Address, Icon, Hours string
	At                         Point
	Rating                     float64
	Reviews                    int
	Open                       bool
	Photo                      *ui.Bitmap
}

// checkPlace panics unless p is a named, real place.
func checkPlace(component string, p Place) {
	if strings.TrimSpace(p.Name) == "" || !(p.Rating >= 0 && p.Rating <= 5) || p.Reviews < 0 {
		panic(fmt.Sprintf("geo: %s place %q needs a name, rating 0 to 5 and reviews >= 0", component, p.Name))
	}
	valid(component, p.Name, p.At)
}

// PlaceCardOptions word the card; Reviews formats the count.
type PlaceCardOptions struct{ Open, Closed, Reviews string }

// PlaceCard shows a place with its photo, rating and hours.
func PlaceCard(c *ui.Context, p Place, o PlaceCardOptions) ui.Element {
	def(&o.Open, "Open")
	def(&o.Closed, "Closed")
	def(&o.Reviews, "(%s)")
	def(&p.Icon, "map-pin")
	checkPlace("PlaceCard", p)
	t := c.Theme()
	root := card(c).Clip().Role(ui.RoleGroup)
	return root.Children(func() {
		if p.Photo != nil {
			c.Key("photo")
			ui.Box(c).FillWidth().AspectRatio(2.4).Clip().Children(func() {
				kit.Arrive(ui.Image(c, p.Photo).Fill().Fit(ui.Cover), ui.Pose{Scale: 1.08, Fade: 1}, 0)
			})
		}
		c.Key("body")
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(4)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start).Children(func() {
				tile(c, p.Icon, t.Accent, t.Space(10))
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, p.Name).FontSize(t.Rem(0.95)).Bold().FillWidth().Wrap()
					if p.Address != "" {
						kit.Muted(c, p.Address).FontSize(t.Rem(0.76)).FillWidth().Wrap()
					}
				})
			})
			if p.Rating > 0 {
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, strconv.FormatFloat(p.Rating, 'f', 1, 64)).FontSize(t.Rem(0.8)).Bold().Shrink(0).SingleLine()
					rating(c, p.Rating, t.Space(3.5))
					kit.Muted(c, fmt.Sprintf(o.Reviews, number(c, p.Reviews))).FontSize(t.Rem(0.74)).Grow(1).MinWidth(0).SingleLine()
				})
			}
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				state, ink := o.Closed, t.Danger
				if p.Open {
					state, ink = o.Open, t.Success
				}
				c.Key(p.Open)
				kit.Arrive(ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(ink.Alpha(0.12)),
					ui.Pose{Scale: 0.9, Fade: 1}, 0).Children(func() {
					ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Radius(999).Background(ink)
					ui.Text(c, state).FontSize(t.Rem(0.72)).Bold().TextColor(ink).Shrink(0).MaxWidth(t.Space(30)).SingleLine()
				})
				if p.Hours != "" {
					glyph(c, "clock", 12).TextColor(t.TextMuted)
					kit.Muted(c, p.Hours).FontSize(t.Rem(0.74)).Grow(1).MinWidth(0).SingleLine()
				}
			})
		})
	})
}
