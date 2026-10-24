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

// RetroOptions are the retro labels.
type RetroOptions struct{ Add string }

// Retro shows went-well, improve, and actions columns; cue Pop.
func Retro(c *ui.Context, columns *[]RetroColumn, o RetroOptions) ui.Element {
	if columns == nil || len(*columns) == 0 {
		panic("project: Retro needs columns")
	}
	def(&o.Add, "Add card")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	addColumn := -1
	root.Children(func() {
		for i, column := range *columns {
			if column.Title == "" {
				panic("project: Retro column title must not be empty")
			}
			c.Key(i)
			zone := ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Padding(t.Space(1)).Radius(t.Radius).Background(t.Surface)
			zone.Children(func() {
				ui.Text(c, column.Title).Bold()
				for _, card := range column.Cards {
					ui.Text(c, card).Padding(t.Space(1)).Radius(t.Space(0.5)).Background(t.Background).Wrap()
				}
				c.Key("add")
				add := kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(1)).Radius(t.Radius))
				add.Children(func() { ui.Text(c, o.Add).FontSize(t.Rem(0.75)).TextColor(t.TextMuted) })
				if add.Clicked() {
					addColumn = i
				}
			})
		}
	})
	if addColumn >= 0 {
		(*columns)[addColumn].Cards = append((*columns)[addColumn].Cards, fmt.Sprintf("%s %d", o.Add, len((*columns)[addColumn].Cards)+1))
		addColumn = -1
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Retro", "added")
	}
	return root
}
