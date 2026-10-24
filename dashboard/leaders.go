package dashboard

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Leader is one performer; Detail is a team or role.
type Leader struct {
	Name, Detail string
	Score        int
}

// LeadersOptions word the card.
type LeadersOptions struct{ Title string }

// Leaders ranks performers by score; the top three wear medals.
func Leaders(c *ui.Context, people []Leader, o LeadersOptions) ui.Element {
	if len(people) == 0 {
		panic("dashboard: Leaders needs people")
	}
	names := make([]string, len(people))
	for i, person := range people {
		if person.Score < 0 || i > 0 && person.Score > people[i-1].Score {
			panic(fmt.Sprintf("dashboard: Leaders %q needs a score of zero or more, sorted high to low", person.Name))
		}
		names[i] = person.Name
	}
	unique("Leaders", names)
	def(&o.Title, "Leaderboard")
	peak := max(1, people[0].Score)
	t := c.Theme()
	medals := [...]ui.Color{t.Warning, t.TextMuted, t.Warning.Mix(t.Danger, 0.45)}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "trophy", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(2))
		list.Children(func() {
			for i, person := range people {
				c.Key(person.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					rank := ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999)
					ink := t.TextMuted
					if i < len(medals) {
						ink = medals[i]
						rank.Background(ink.Alpha(0.18))
					}
					rank.Children(func() {
						kit.Numeric(c, strconv.Itoa(i+1)).FontSize(t.Rem(0.75)).Bold().TextColor(ink).Shrink(0)
					})
					display.Avatar(c, person.Name, nil)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, person.Name).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							kit.Numeric(c, number(c, person.Score)).Bold().FontSize(t.Rem(0.85)).Shrink(0).
								ContentTransition(ui.ContentTransition{Numeric: true})
						})
						if person.Detail != "" {
							kit.Muted(c, person.Detail).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
						}
						progress(c, float64(person.Score)/float64(peak), t.Accent)
					})
				})
			}
		})
	})
	return root
}
