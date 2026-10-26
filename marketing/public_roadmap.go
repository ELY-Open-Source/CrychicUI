package marketing

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RoadStage places a roadmap item.
type RoadStage int

const (
	RoadPlanned RoadStage = iota
	RoadBuilding
	RoadShipped
)

// RoadItem is one roadmap card; zero Votes hides the count.
type RoadItem struct {
	Title, Text string
	Stage       RoadStage
	Votes       int
}

// PublicRoadmapOptions word the board; Vote is the singular.
type PublicRoadmapOptions struct {
	Title, Text, Empty, Vote, Votes string
	Stages                          [3]string
}

// PublicRoadmap shows planned, in progress and done.
func PublicRoadmap(c *ui.Context, items []RoadItem, o PublicRoadmapOptions) ui.Element {
	def(&o.Empty, "Nothing here yet")
	def(&o.Vote, "%d vote")
	def(&o.Votes, "%d votes")
	for i, d := range [3]string{"Planned", "In progress", "Done"} {
		def(&o.Stages[i], d)
	}
	titles := make([]string, len(items))
	for i, it := range items {
		if it.Stage < RoadPlanned || it.Stage > RoadShipped || it.Votes < 0 {
			panic(fmt.Sprintf("marketing: PublicRoadmap %q item %q needs a known stage and votes of 0 or more", o.Title, it.Title))
		}
		titles[i] = it.Title
	}
	unique("PublicRoadmap", titles)
	t := c.Theme()
	marks := [3]string{"circle-dashed", "loader-circle", "circle-check"}
	inks := [3]ui.Color{t.TextMuted, t.Accent, t.Success}
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Start).Children(func() {
			for s := range 3 {
				var col []RoadItem
				for _, it := range items {
					if int(it.Stage) == s {
						col = append(col, it)
					}
				}
				lane := ui.Column(c).Grow(1).Basis(t.Space(40)).MinWidth(0).Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius * 1.5).
					Background(t.SurfaceHover.Alpha(0.6))
				lane.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, 0, t.Space(1), 0).TextColor(inks[s]).Children(func() {
						glyph(c, marks[s], 16)
						ui.Text(c, o.Stages[s]).FontSize(t.Rem(0.84)).Bold().TextColor(t.Text).Grow(1).MinWidth(0).SingleLine()
						kit.Chip(c, strconv.Itoa(len(col)), inks[s])
					})
					if len(col) == 0 {
						kit.Muted(c, o.Empty).FontSize(t.Rem(0.78)).FillWidth().TextAlign(ui.Center).Padding(t.Space(4), 0).Wrap()
						return
					}
					cards := ui.Column(c).FillWidth().Gap(t.Space(2))
					cards.Children(func() {
						for i, it := range col {
							c.Key(strconv.Quote(it.Title))
							box := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).
								Background(t.Background)
							kit.Stagger(c, cards, box, i).Children(func() {
								ui.Text(c, it.Title).FontSize(t.Rem(0.84)).Bold().FillWidth().Wrap()
								if it.Text != "" {
									kit.Muted(c, it.Text).FontSize(t.Rem(0.76)).FillWidth().Wrap()
								}
								if it.Votes > 0 {
									ui.Row(c).FillWidth().MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
										glyph(c, "thumbs-up", 12)
										kit.Muted(c, counted(it.Votes, o.Vote, o.Votes)).FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).SingleLine()
									})
								}
							})
						}
					})
				})
			}
		})
	})
	return root
}
