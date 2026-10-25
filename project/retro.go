package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RetroColumn is one retrospective column.
type RetroColumn struct {
	Title string
	Cards []string
}

// RetroOptions word the add button and a new card.
type RetroOptions struct{ Add, NewCard string }

// Retro shows went-well, improve, and actions columns; cue Pop.
func Retro(c *ui.Context, columns *[]RetroColumn, o RetroOptions) ui.Element {
	if columns == nil || len(*columns) == 0 {
		panic("project: Retro needs columns")
	}
	titles := make([]string, len(*columns))
	for i, column := range *columns {
		titles[i] = column.Title
		for _, text := range column.Cards {
			if text == "" {
				panic(fmt.Sprintf("project: Retro column %q has an empty card", column.Title))
			}
		}
	}
	unique("Retro", titles)
	def(&o.Add, "Add card")
	def(&o.NewCard, "New card")
	t := c.Theme()
	looks := []struct {
		glyph string
		ink   ui.Color
	}{{"thumbs-up", t.Success}, {"wrench", t.Warning}, {"rocket", t.Accent}}
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Role(ui.RoleGroup)
	addColumn := -1
	root.Children(func() {
		for i, column := range *columns {
			c.Key(column.Title)
			look := looks[i%len(looks)]
			zone := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(2)).Padding(t.Space(2)).Radius(t.Radius).
				Background(look.ink.Alpha(0.07))
			kit.Arrive(zone, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			zone.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(0, t.Space(0.5)).Children(func() {
					glyph(c, look.glyph, 14).TextColor(look.ink)
					ui.Text(c, column.Title).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
					kit.Chip(c, fmt.Sprint(len(column.Cards)), look.ink)
				})
				cards := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				cards.Children(func() {
					for j, text := range column.Cards {
						c.Key(j)
						note := ui.Text(c, text).FontSize(t.Rem(0.78)).Wrap().FillWidth().Padding(t.Space(2)).Radius(t.Radius).
							Background(t.Background).Border(1, t.Border).Shadow(0, 1, 2, 0, kit.Shade(0.05))
						kit.Stagger(c, cards, note, j)
					}
				})
				add := kit.Press(ui.ButtonBase(c).Label(o.Add).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5)).Radius(t.Radius))
				add.Background(add.AnimateColor("bg", look.ink.Alpha(0).Mix(look.ink.Alpha(0.12), unit(add.Hovered())), ui.Snappy))
				add.TextColor(add.AnimateColor("ink", t.TextMuted.Mix(look.ink, unit(add.Hovered())), ui.Snappy))
				add.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						glyph(c, "plus", 14)
						ui.Text(c, o.Add).FontSize(t.Rem(0.75)).Shrink(0)
					})
				})
				if add.Clicked() {
					addColumn = i
				}
			})
		}
	})
	if addColumn >= 0 {
		(*columns)[addColumn].Cards = append((*columns)[addColumn].Cards, o.NewCard)
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Retro", "added", "column", addColumn, "cards", len((*columns)[addColumn].Cards))
	}
	return root
}
