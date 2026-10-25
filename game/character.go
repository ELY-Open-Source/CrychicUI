package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Hero is a character card's data; MaxMP 0 hides mana.
type Hero struct {
	Name, Class string
	Level       int
	HP, MaxHP   int
	MP, MaxMP   int
	Stats       []Stat
}

// CharacterOptions word the card; Points takes current and max.
type CharacterOptions struct{ Level, HP, MP, Points string }

// Character is a character card with level, points and stats.
func Character(c *ui.Context, h Hero, o CharacterOptions) ui.Element {
	if h.Name == "" || h.Class == "" || h.Level < 1 || h.MaxMP < 0 || h.MP < 0 || h.MP > h.MaxMP {
		panic(fmt.Sprintf("game: Character %q needs a name, class, level and 0 <= mp <= max", h.Name))
	}
	names := make([]string, len(h.Stats))
	for i, s := range h.Stats {
		names[i] = s.Name
	}
	if len(names) > 0 {
		unique("Character", names)
	}
	def(&o.Level, "Lv %d")
	def(&o.HP, "HP")
	def(&o.MP, "MP")
	def(&o.Points, "%d / %d")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Shrink(0).Children(func() {
				display.Avatar(c, h.Name, nil).Size(t.Space(14), t.Space(14))
				ui.Box(c).Absolute().Right(-t.Space(1)).Bottom(-t.Space(1)).Padding(t.Space(0.25), t.Space(1.5)).Radius(999).
					Background(t.Warning).TextColor(t.Background).Children(func() {
					ui.Text(c, fmt.Sprintf(o.Level, h.Level)).FontSize(t.Rem(0.62)).Bold().Shrink(0)
				})
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				ui.Text(c, h.Name).FontSize(t.Rem(1)).Bold().FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).Children(func() { kit.Chip(c, h.Class, t.Accent) })
			})
		})
		HealthBar(c, h.HP, h.MaxHP, HealthBarOptions{Label: o.HP, Count: o.Points})
		if h.MaxMP > 0 {
			ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "sparkles", 14).TextColor(t.Accent)
					ui.Text(c, o.MP).FontSize(t.Rem(0.75)).Bold().Grow(1).MinWidth(0).SingleLine()
					kit.Numeric(c, fmt.Sprintf(o.Points, h.MP, h.MaxMP)).Font("monospace").FontSize(t.Rem(0.75)).Bold().Shrink(0)
				})
				meter(c, float32(float64(h.MP)/float64(h.MaxMP)), t.Accent).Height(t.Space(3.5))
			})
		}
		if len(h.Stats) == 0 {
			return
		}
		grid := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2))
		grid.Children(func() {
			for i, s := range h.Stats {
				c.Key(s.Name)
				cell := ui.Column(c).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(2)).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
				kit.Stagger(c, grid, cell, i).Children(func() {
					kit.Muted(c, s.Name).FontSize(t.Rem(0.65)).Bold().FillWidth().MinWidth(0).SingleLine()
					kit.Numeric(c, fmt.Sprint(s.Value)).FontSize(t.Rem(1)).Bold().Shrink(0)
				})
			}
		})
	})
	return root
}
